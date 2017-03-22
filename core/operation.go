package core

import (
    "time"
    "sync"
)

type RenderContext struct {
    Pipeline *Pipeline

    Duration time.Duration

    Results map[string]LEDStripeReader
}

type Operation interface {
    sync.Locker

    Name() string

    Render(context *RenderContext) LEDStripeReader
}

func (this *RenderContext) Count() int {
    return this.Pipeline.Count
}
