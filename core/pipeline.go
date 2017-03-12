package core

import (
    "fmt"
    "reflect"
    "sync"
    "github.com/mitchellh/mapstructure"
)

type Pipeline []Operation

func NewPipeline() Pipeline {
    return make([]Operation, 0)
}

func (p *Pipeline) ByName(name string) Operation {
    for _, op := range *p {
        if op.Name() == name {
            return op
        }
    }

    return nil
}

type Operation interface {
    sync.Locker

    Name() string
    Stripe() LEDStripe

    Render()
}

type OperationFactory struct {
    ConfigType reflect.Type

    Create func(pipeline *Pipeline, name string, count int, config interface{}) (Operation, error)
}

var operationFactories = make(map[string]OperationFactory)

func RegisterOperation(t string, factory OperationFactory) {
    if _, found := operationFactories[t]; found {
        panic(fmt.Errorf("Duplicated operation type: %s", t))
    }

    operationFactories[t] = factory
}

func (p *Pipeline) NewOperation(t string, name string, count int, config map[string]interface{}) {
    f, found := operationFactories[t]
    if !found {
        panic(fmt.Errorf("Unknown operation type: %s", t))
    }

    s := reflect.New(f.ConfigType).Interface()

    err := mapstructure.Decode(config, s)
    if err != nil {
        panic(err)
    }

    op, err := f.Create(p, name, count, s)
    if err != nil {
        panic(err)
    }

    *p = append(*p, op)
}
