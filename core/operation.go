package core

import (
    "time"
    "sync"
)

type RenderContext struct {
    Duration time.Duration

    Stripe LEDStripe

    Pipeline *Pipeline
}

func (this *RenderContext) Count() int {
    return len(this.Stripe)
}

func (this *RenderContext) Set(i int, r uint8, g uint8, b uint8) {
    this.Stripe[i].R = r
    this.Stripe[i].G = g
    this.Stripe[i].B = b
    this.Stripe[i].A = 0
}

type Operation interface {
    sync.Locker

    Name() string

    Render(context *RenderContext)
}
