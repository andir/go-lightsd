package outputs

import (
    "github.com/andir/lightsd/core"
    "reflect"
    "log"
)

type Factory struct {
    ConfigType reflect.Type

    Create func(count int, operation string, config interface{}) (core.Output, error)
}

var factories = make(map[string]*Factory)

func Register(t string, factory *Factory) {
    if _, found := factories[t]; found {
        log.Panicf("Duplicated output type: %s", t)
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
