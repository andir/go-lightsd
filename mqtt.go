package main

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"reflect"
	"log"
	"fmt"
	"strconv"
)


func NewMqttConnection(broker string, clientId string, pipeline Pipeline) {

	opts := MQTT.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientId)

	client := MQTT.NewClient(opts)

	if token :=client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for name, op := range pipeline {
		t := reflect.TypeOf(op).Elem()

		for i:= 0; i < t.NumField(); i++ {
			f := t.Field(i)

			tag, found := f.Tag.Lookup("mqtt")
			if !found {
				continue
			}

			topic := fmt.Sprintf("lightsd/%s/%s/set", name, tag)
			log.Printf("Found MQTT exported parameter: %s:%s(%s) as %s", t.Name(), f.Name, f.Type.Name(), topic)

			client.Subscribe(topic, 0, func(c MQTT.Client, m MQTT.Message) {
				s := string(m.Payload())

				log.Printf("Setting value: %s:%s(%s) = %s", t.Name(), f.Name, f.Type.Name(), s)

				log.Printf("XXX: %v", reflect.ValueOf(op).Elem().Field(i))

				switch f.Type.Kind() {
				case reflect.Float64:
					val, err := strconv.ParseFloat(s, 64)
					if err != nil {
						log.Printf("Failed to parse float: %s: %v", s, err)
					}

					reflect.ValueOf(op).Field(i).SetFloat(val)

				case reflect.Int:
					val, err := strconv.ParseInt(s, 10, 64)
					if err != nil {
						log.Printf("Failed to parse int: %s: %v", s, err)
					}

					reflect.ValueOf(op).Field(i).SetInt(val)

				case reflect.Bool:
					val, err := strconv.ParseBool(s)
					if err != nil {
						log.Printf("Failed to parse bool: %s: %v", s, err)
					}

					reflect.ValueOf(op).Field(i).SetBool(val)

				case reflect.String:
					reflect.ValueOf(op).Field(i).SetString(s)
				}
			})
		}
	}

}