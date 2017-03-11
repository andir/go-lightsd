package core

import (
	"fmt"
)

type Pipeline []Operation

type Operation interface {
	Name() string

	Lock()
	Unlock()
	RLock()
	RUnlock()

	Render(LEDStripe) LEDStripe
}

type OperationFactory func(settings interface{}) (Operation, error)

var operationFactories = make(map[string]OperationFactory)

func RegisterOperation(name string, factory OperationFactory) {
	if _, found := operationFactories[name]; !found {
		panic(fmt.Errorf("Duplicated operation name: %s", name))
	}

	operationFactories[name] = factory
}

