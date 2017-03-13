package operations

import (
    "github.com/andir/lightsd/core"
    "reflect"
    "fmt"
)

type Factory struct {
    ConfigType reflect.Type

    Create func(name string, count int, config interface{}) (core.Operation, error)
}

var factories = make(map[string]*Factory)

func Register(t string, factory *Factory) {
    if _, found := factories[t]; found {
        panic(fmt.Errorf("Duplicated operation type: %s", t))
    }

    factories[t] = factory
}

func Get(t string) *Factory {
    f, found := factories[t]
    if !found {
        panic(fmt.Errorf("Unknown operation type: %s", t))
    }

    return f
}
