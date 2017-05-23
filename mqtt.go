package main

import (
    "github.com/andir/lightsd/core"
    MQTT "github.com/eclipse/paho.mqtt.golang"
    "reflect"
    "log"
    "fmt"
    conv "github.com/cstockton/go-conv" // TODO: Wait for merge in upstream
)

type MqttConnection struct {
    client MQTT.Client
    realm  string
}

func NewMqttConnection(config *Config) (*MqttConnection, error) {
    opts := MQTT.NewClientOptions()
    opts.AddBroker(fmt.Sprintf("tcp://%s:%d/", config.MQTT.Host, config.MQTT.Port))
    opts.SetClientID(config.MQTT.ClientID)

    client := MQTT.NewClient(opts)

    if token := client.Connect(); token.Wait() && token.Error() != nil {
        return nil, token.Error()
    }

    return &MqttConnection{
        client: client,
        realm:  config.MQTT.Realm,
    }, nil
}

func (this *MqttConnection) Register(pipeline *core.Pipeline) error {
    for _, operation := range pipeline.Operations {
        v := reflect.ValueOf(operation).Elem()
        t := v.Type()

        for i := 0; i < t.NumField(); i++ {
            fieldType := t.Field(i)
            fieldValue := v.Field(i)

            tag, found := fieldType.Tag.Lookup("mqtt")
            if !found {
                continue
            }

            topic := fmt.Sprintf("%s/%s/%s/%s", this.realm, pipeline.Name, operation.Name(), tag)

            log.Printf("Found MQTT exported parameter: %s:%s(%s) as %s/set", t.Name(), fieldType.Name, fieldType.Type.Name(), topic)

            handler := func(c MQTT.Client, m MQTT.Message) {
                msg := string(m.Payload())

                // Lock the operation for changes, parse the message into the field
                {
                    operation.Lock()
                    defer operation.Unlock()

                    // Parse the messages to a value
                    if err := conv.Infer(fieldValue, msg); err != nil {
                        log.Printf("Failed to parse: %s: %v", msg, err)
                        return
                    }
                }

                // Publish the updated value
                if t := this.client.Publish(topic, 0, false, msg); t.Wait() && t.Error() != nil {
                    log.Printf("Failed to publish: %s=%s: %v", topic, msg, t.Error())
                    return
                }

                log.Printf("Changed exported parameter: %s:%s(%s) = %v", t.Name(), fieldType.Name, fieldType.Type.Name(), msg)
            }

            if t := this.client.Subscribe(fmt.Sprintf("%s/set", topic), 0, handler); t.Wait() && t.Error() != nil {
                return t.Error()
            }
        }
    }

    return nil
}
