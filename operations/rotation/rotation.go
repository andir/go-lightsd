package operations

import (
    "sync"
    "github.com/andir/lightsd/core"
    "github.com/andir/lightsd/operations"
    "reflect"
)

type RotationConfig struct {
    Source         string `mapstructure:"source"`
    PixelPerSecond float64 `mapstructure:"speed"`
}

type Rotation struct {
    sync.RWMutex

    name string

    Source         string
    PixelPerSecond float64  `mqtt:"speed"`

    offset float64
}

func (this *Rotation) Name() string {
    return this.name
}

func (this *Rotation) Render(context *core.RenderContext) {
    this.offset += context.Duration.Seconds() * this.PixelPerSecond

    // TODO: Ouch, this hurts
    source := context.Pipeline.ByName(this.Source)
    for i, s := range source.Stripe() {
        context.Stripe[(i+int(this.offset))%context.Count()] = s
    }
}

func init() {
    operations.Register("rotation", &operations.Factory{
        ConfigType: reflect.TypeOf(RotationConfig{}),
        Create: func(name string, count int, rconfig interface{}) (core.Operation, error) {
            config := rconfig.(*RotationConfig)

            //source := pipeline.ByName(config.Source)
            //if source == nil {
            //    return nil, fmt.Errorf("Unknown source: %s", config.Source)
            //}

            return &Rotation{
                name: name,

                Source:         config.Source,
                PixelPerSecond: config.PixelPerSecond,

                offset: 0.0,
            }, nil
        },
    })
}
