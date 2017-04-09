package operations

import (
    "github.com/andir/lightsd/core"
    "github.com/andir/lightsd/operations"
    "sync"
    "reflect"
    "github.com/lucasb-eyer/go-colorful"
    "fmt"
)

type AdjustmentConfig struct {
    Source string `mapstructure:"source"`

    Correction  string `mapstructure:"correction"`
    Temperature string `mapstructure:"temperature"`
}

type Adjustment struct {
    sync.RWMutex

    name string

    source string

    adjustment struct{ R, G, B float64 }
}

type adjustingLEDStripeReader struct {
    source     core.LEDStripeReader
    adjustment struct{ R, G, B float64 }
}

func (this *adjustingLEDStripeReader) Count() int {
    return this.source.Count()
}

func (this *adjustingLEDStripeReader) Get(i int) core.LED {
    c := this.source.Get(i)

    return core.LED{
        R: c.R * this.adjustment.R,
        G: c.G * this.adjustment.G,
        B: c.B * this.adjustment.B,
    }
}

func (this *Adjustment) Name() string {
    return this.name
}

func (this *Adjustment) Render(context *core.RenderContext) core.LEDStripeReader {
    // TODO: Ouch, this hurts
    source := context.Results[this.source]

    return &adjustingLEDStripeReader{
        source:     source,
        adjustment: this.adjustment,
    }
}

func init() {
    operations.Register("adjustment", &operations.Factory{
        ConfigType: reflect.TypeOf(AdjustmentConfig{}),
        Create: func(name string, count int, rconfig interface{}) (core.Operation, error) {
            config := rconfig.(*AdjustmentConfig)

            correction, err := colorful.Hex(config.Correction)
            if err != nil {
                return nil, fmt.Errorf("adjustment: failed to parse correction color: %v", err)
            }

            temperature, err := colorful.Hex(config.Temperature)
            if err != nil {
                return nil, fmt.Errorf("adjustment: failed to parse temperature color: %v", err)
            }

            return &Adjustment{
                name: name,

                source: config.Source,

                adjustment: struct{ R, G, B float64 }{
                    R: correction.R * temperature.R,
                    G: correction.R * temperature.G,
                    B: correction.R * temperature.B,
                },
            }, nil
        },
    })
}
