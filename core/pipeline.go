package core

import (
    "time"
)

type Pipeline struct {
    Name string

    Count int

    Operations []Operation
    Output     Output

    lastRendered time.Time
    results map[string]LEDStripeReader
}

func NewPipeline(name string, count int, output Output, operations []Operation) *Pipeline {
    return &Pipeline{
        Name: name,

        Count: count,

        Operations: operations,
        Output:     output,

        lastRendered: time.Now(),
        results:make(map[string]LEDStripeReader, len(operations)),
    }
}

func (p *Pipeline) Render(duration time.Duration) *RenderContext {
    context := &RenderContext{
        Pipeline: p,
        Duration: duration,
        Results: p.results,
    }

    for _, op := range p.Operations {
        op.Lock()

        result := op.Render(context)
        context.Results[op.Name()] = result

        op.Unlock()
    }

    p.Output.Render(context.Results[p.Output.Source()])

    return context
}
