package operations

import (
    "time"
    "sync"
    "../core"
    "reflect"
    "fmt"
)

type RotationConfig struct {
    Source         string
    PixelPerSecond float64
}

type Rotation struct {
    sync.RWMutex

    name   string
    stripe core.LEDStripe

    Source         core.Operation
    PixelPerSecond float64  `mqtt:"speed"`

    offset float64

    LastFrameTime time.Time
}

func (this *Rotation) Name() string {
    return this.name
}

func (this *Rotation) Stripe() core.LEDStripe {
    return this.stripe
}

func (this *Rotation) Render() {
    this.RLock()
    defer this.RUnlock()

    now := time.Now()
    delta := now.Sub(this.LastFrameTime)

    this.offset += delta.Seconds() * this.PixelPerSecond
    this.LastFrameTime = now

    iOffset := int(this.offset) % len(this.stripe)

    for i, s := range this.Source.Stripe() {
        this.stripe[(i+iOffset)%len(this.stripe)] = s
    }
}

func init() {
    core.RegisterOperation("rotation", core.OperationFactory{
        ConfigType: reflect.TypeOf(RotationConfig{}),
        Create: func(pipeline *core.Pipeline, name string, count int, rconfig interface{}) (core.Operation, error) {
            config := rconfig.(*RotationConfig)

            source := pipeline.ByName(config.Source)
            if source == nil {
                return nil, fmt.Errorf("Unknown source: %s", config.Source)
            }

            return &Rotation{
                name: name,

                stripe: core.NewLEDStripe(count),

                Source:         source,
                PixelPerSecond: config.PixelPerSecond,

                offset: 0.0,

                LastFrameTime: time.Time{},
            }, nil
        },
    })
}
