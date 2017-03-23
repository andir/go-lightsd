package operations

import (
    "github.com/andir/lightsd/core"
    "reflect"
    "log"
    "fmt"
    "github.com/mitchellh/mapstructure"
)

type Factory struct {
    ConfigType reflect.Type

    Create func(name string, count int, config interface{}) (core.Operation, error)
}

var factories = make(map[string]*Factory)

func Register(t string, factory *Factory) {
    if _, found := factories[t]; found {
        log.Panic("operations: Duplicated operation type:", t)
    }

    factories[t] = factory
}

func Make(t string, name string, count int, configData map[string]interface{}) (core.Operation, error) {
    factory, found := factories[t]
    if !found {
        return nil, fmt.Errorf("operations: Unknown type: %s", t)
    }

    config := reflect.New(factory.ConfigType).Interface()
    err := mapstructure.Decode(configData, config)
    if err != nil {
        return nil, fmt.Errorf("operations: Failed to build config: %v", err)
    }

    operation, err := factory.Create(name, count, config)
    if err != nil {
        return nil, err
    }

    return operation, nil
}
