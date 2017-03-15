package core

import (
    "time"
)

type Pipeline struct {
    name string

    count int

    operations []Operation
    output     Output

    lastRendered time.Time
}

func NewPipeline(name string, count int, output Output, operations []Operation) *Pipeline {
    return &Pipeline{
        name: name,

        count: count,

        operations: operations,
        output:     output,

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

func (p *Pipeline) Operations() []Operation {
    return p.operations
}

func (p *Pipeline) ByName(name string) *Operation {
    for _, op := range p.operations {
        if op.Name() == name {
            return &op
        }
    }

    return nil
}

func (p *Pipeline) Render(duration time.Duration) *RenderContext {
    context := &RenderContext{
        Count: p.count,
        Duration: duration,
        Results: make(map[string]LEDStripeReader, len(p.operations)),
    }

    for _, op := range p.operations {
        op.Lock()

        result := op.Render(context)
        context.Results[op.Name()] = result

        op.Unlock()
    }

    p.output.Render(context.Results[p.output.Source()])

    return context
}
