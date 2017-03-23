package outputs

import (
    "github.com/andir/lightsd/core"
    "reflect"
    "log"
    "github.com/mitchellh/mapstructure"
    "fmt"
)

type Factory struct {
    ConfigType reflect.Type

    Create func(count int, source string, config interface{}) (core.Output, error)
}

var factories = make(map[string]*Factory)

func Register(t string, factory *Factory) {
    if _, found := factories[t]; found {
        log.Panicf("Duplicated output type: %s", t)
    }

    factories[t] = factory
}

func Make(t string, count int, source string, configData map[string]interface{}) (core.Output, error) {
    factory, found := factories[t]
    if !found {
        return nil, fmt.Errorf("outputs: Unknown type: %s", t)
    }

    config := reflect.New(factory.ConfigType).Interface()
    err := mapstructure.Decode(configData, config)
    if err != nil {
        return nil, fmt.Errorf("outputs: Failed to build config: %v", err)
    }

    output, err := factory.Create(count, source, config)
    if err != nil {
        return nil, err
    }

    return output, nil
}
