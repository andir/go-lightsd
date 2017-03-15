package core

import (
    "time"
    "sync"
)

type RenderContext struct {
    Count int
    Duration time.Duration

    Results map[string]LEDStripeReader
}

type Operation interface {
    sync.Locker

    Name() string

    Render(context *RenderContext) LEDStripeReader
}
