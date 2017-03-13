package core

import (
    "time"
)

type Processor struct {
    operation Operation
    stripe    LEDStripe
}

func (this *Processor) Name() string {
    return this.operation.Name()
}

func (this *Processor) Lock() {
    this.operation.Lock()
}

func (this *Processor) Unlock() {
    this.operation.Unlock()
}

func (this *Processor) Operation() Operation {
    return this.operation
}

func (this *Processor) Stripe() LEDStripe {
    return this.stripe
}

type Pipeline struct {
    name string

    count int

    output     Output
    operations []Processor

    lastRendered time.Time
}

func NewPipeline(name string, count int, output Output, processors []Operation) *Pipeline {
    operations := make([]Processor, len(processors))
    for i, processor := range processors {
        operations[i] = Processor{
            operation: processor,
            stripe:    NewLEDStripe(count),
        }
    }

    return &Pipeline{
        name: name,

        count: count,

        output:     output,
        operations: operations,

        lastRendered: time.Now(),
    }
}

func (p *Pipeline) Name() string {
    return p.name
}

func (p *Pipeline) Count() int {
    return p.count
}

func (p *Pipeline) Output() Output {
    return p.output
}

func (p *Pipeline) Processors() []Processor {
    return p.operations
}

func (p *Pipeline) ByName(name string) *Processor {
    for _, op := range p.operations {
        if op.operation.Name() == name {
            return &op
        }
    }

    return nil
}

func (p *Pipeline) Render(duration time.Duration) {
    for _, op := range p.operations {
        op.operation.Lock()
        op.operation.Render(&RenderContext{
            Duration: duration,
            Stripe:   op.stripe,
            Pipeline: p,
        })
        op.operation.Unlock()
    }
}
