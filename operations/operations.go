package operations

import (
    "github.com/andir/lightsd/core"
    "reflect"
    "log"
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

func Get(t string) *Factory {
    f, found := factories[t]
    if !found {
        return nil
    }

    return f
}
