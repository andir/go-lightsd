package operations

import (
    "github.com/andir/lightsd/core"
    "github.com/andir/lightsd/utils"
    "github.com/andir/lightsd/operations"
    "reflect"
    "sync"
    "github.com/lucasb-eyer/go-colorful"
    "fmt"
)

type GradientConfig struct {
    Gradient []struct {
        C string `mapstructure:"c"`
        P float64 `mapstructure:"p"`
    } `mapstructure:"gradient"`
}

type Gradient struct {
    sync.RWMutex

    name   string
    stripe core.LEDStripe
}

func (this *Gradient) Name() string {
    return this.name
}

func (this *Gradient) Render(context *core.RenderContext) core.LEDStripeReader {
    return this.stripe
}

func init() {
    operations.Register("gradient", &operations.Factory{
        ConfigType: reflect.TypeOf(GradientConfig{}),
        Create: func(name string, count int, rconfig interface{}) (core.Operation, error) {
            config := rconfig.(*GradientConfig)

            gradient := make(utils.GradientTable, len(config.Gradient))
            for i, e := range config.Gradient {
                color, err := colorful.Hex(e.C)
                if err != nil {
                    return nil, fmt.Errorf("gradient: failed to parse color: %v: %v", e.C, err)
                }

                gradient[i].Col = color
                gradient[i].Pos = e.P
            }

            stripe := core.NewLEDStripe(count)
            for i := 0; i < count; i++ {
                c := gradient.GetInterpolatedColorFor(float64(i) / float64(count-1))
                stripe.Set(i, core.LED{R: c.R, G: c.G, B: c.B})
            }

            return &Gradient{
                name:   name,
                stripe: stripe,
            }, nil
        },
    })
}
