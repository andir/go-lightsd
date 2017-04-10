package operations

import (
    "github.com/andir/lightsd/core"
    "github.com/andir/lightsd/operations"
    "sync"
    "reflect"
)

type DimmerConfig struct {
    Source string `mapstructure:"source"`

    Value float64 `mapstructure:"value"`
}

type Dimmer struct {
    sync.RWMutex

    name string

    source string

    Value float64 `mqtt:"value"`
}

type dimmerLEDStripeReader struct {
    source core.LEDStripeReader
    value  float64
}

func (this *dimmerLEDStripeReader) Get(i int) core.LED {
    c := this.source.Get(i)

    return core.LED{
        R: c.R * this.value,
        G: c.G * this.value,
        B: c.B * this.value,
    }
}

func (this *Dimmer) Name() string {
    return this.name
}

func (this *Dimmer) Render(context *core.RenderContext) core.LEDStripeReader {
    // TODO: Ouch, this hurts
    source := context.Results[this.source]

    return &dimmerLEDStripeReader{
        source: source,
        value:  this.Value,
    }
}

func init() {
    operations.Register("dimmer", &operations.Factory{
        ConfigType: reflect.TypeOf(DimmerConfig{}),
        Create: func(name string, count int, rconfig interface{}) (core.Operation, error) {
            config := rconfig.(*DimmerConfig)

            return &Dimmer{
                name: name,

                source: config.Source,

                Value: config.Value,
            }, nil
        },
    })
}
