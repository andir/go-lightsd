package core

import (
    "fmt"
    "reflect"
    "sync"
    "github.com/mitchellh/mapstructure"
    "time"
)

type Operation interface {
    sync.Locker

    Name() string
    Stripe() LEDStripe

    Update(duration time.Duration)
    Render()
}

type Pipeline struct {
    operations []Operation

    lastRendered time.Time
}

func NewPipeline() *Pipeline {
    return &Pipeline{
        operations: make([]Operation, 0),
        lastRendered: time.Now(),
    }
}

func (p *Pipeline) Operations() []Operation {
    return p.operations
}

func (p *Pipeline) ByName(name string) Operation {
    for _, op := range p.operations {
        if op.Name() == name {
            return op
        }
    }

    return nil
}

func (p *Pipeline) Result() LEDStripe {
    return p.operations[len(p.operations) - 1].Stripe()
}

func (p *Pipeline) Render() time.Duration {
    now := time.Now()

    duration := now.Sub(p.lastRendered)

    for _, op := range p.operations {
        op.Lock()
        op.Update(duration)
        op.Unlock()
        op.Render()
    }

    p.lastRendered = now

    return duration
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

    p.operations = append(p.operations, op)
}
