package main

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"reflect"
	"log"
)


func NewMqttConnection(broker, clientId string, pipeline []Operation) {

	opts := MQTT.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientId)



	client := MQTT.NewClient(opts)

	if token :=client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for _, op := range pipeline {
		var foo interface{} = op

		reflect.Indirect()

		t := reflect.TypeOf(foo)
		log.Println(t)
		for i:= 0; i < t.NumField(); i++ {
			f := t.Field(i)
			name := f.Name
			switch f.Type.Name() {
			default:
				log.Printf("type: %v", f.Type.Name())
			case "float32":
				log.Printf("float32 %v", name)
			}
		}
	}

}