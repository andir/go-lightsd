package main

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"reflect"
	"log"
	"fmt"
	"strconv"
	"github.com/andir/lightsd/core"
)


func NewMqttConnection(config *Config, pipeline *core.Pipeline) {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d/", config.MQTT.Host, config.MQTT.Port))
	opts.SetClientID(config.MQTT.ClientID)

	client := MQTT.NewClient(opts)

	if token :=client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for _, op := range pipeline.Operations() {
		v := reflect.ValueOf(op).Elem()
		t := v.Type()

		for i:= 0; i < t.NumField(); i++ {
			fieldType := t.Field(i)
			fieldValue := v.Field(i)

			tag, found := fieldType.Tag.Lookup("mqtt")
			if !found {
				continue
			}

			topic := fmt.Sprintf("%s/%s/%s/set", config.MQTT.Realm, op.Name(), tag)

			var parse func(s string) (reflect.Value, error)

			switch k := fieldType.Type.Kind(); k {
			case reflect.Float64:
				parse = func(s string) (reflect.Value, error) {
					val, err := strconv.ParseFloat(s, 64)
					if err != nil {
						return reflect.ValueOf(nil), err
					}

					return reflect.ValueOf(val), nil
				}

			case reflect.Int:
				parse = func(s string) (reflect.Value, error) {
					val, err := strconv.ParseInt(s, 10, 64)
					if err != nil {
						return reflect.ValueOf(nil), err
					}

					return reflect.ValueOf(val), nil
				}

			case reflect.Bool:
				parse = func(s string) (reflect.Value, error) {
					val, err := strconv.ParseBool(s)
					if err != nil {
						return reflect.ValueOf(nil), err
					}

					return reflect.ValueOf(val), nil
				}

			case reflect.String:
				parse = func(s string) (reflect.Value, error) {
					return reflect.ValueOf(s), nil
				}

			default:
				log.Fatalf("Unsupported type: %v", k)
			}

			log.Printf("Found MQTT exported parameter: %s:%s(%s) as %s", t.Name(), fieldType.Name, fieldType.Type.Name(), topic)

			if t := client.Subscribe(topic, 0, func(c MQTT.Client, m MQTT.Message) {
				val, err := parse(string(m.Payload()))
				if err != nil {
					log.Printf("Failed to parse: %s: %v", m.Payload(), err)
					return
				}

				log.Printf("Setting fieldValue: %s:%s(%s) = %v", t.Name(), fieldType.Name, fieldType.Type.Name(), val)
				op.Lock()
				defer op.Unlock()
				fieldValue.Set(val)
			}); t.Wait() && t.Error() != nil {
				log.Fatal(t.Error())
			}
		}
	}

}