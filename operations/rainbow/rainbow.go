package operations

import (
    "github.com/andir/lightsd/core"
    "github.com/andir/lightsd/utils"
    "github.com/andir/lightsd/operations"
    "reflect"
    "sync"
)

type RainbowConfig struct {
    Gradient []struct {
        C string
        P float64
    }
}

type Rainbow struct {
    sync.RWMutex

    name   string
    stripe core.LEDStripe
}

func (this *Rainbow) Name() string {
    return this.name
}

func (this *Rainbow) Render(context *core.RenderContext) core.LEDStripeReader {
    return this.stripe
}

func init() {
    operations.Register("rainbow", &operations.Factory{
        ConfigType: reflect.TypeOf(RainbowConfig{}),
        Create: func(name string, count int, rconfig interface{}) (core.Operation, error) {
            config := rconfig.(*RainbowConfig)

            gradient := make(utils.GradientTable, len(config.Gradient))
            for i, e := range config.Gradient {
                gradient[i].Col = utils.ParseColorHex(e.C)
                gradient[i].Pos = e.P
            }

            stripe := core.NewLEDStripe(count)
            for i := 0; i < count; i++ {
                pos := float64(i) / float64(count-1)
                c := gradient.GetInterpolatedColorFor(pos)
                stripe[i].R, stripe[i].G, stripe[i].B = c.RGB255()
            }

            return &Rainbow{
                name:   name,
                stripe: stripe,
            }, nil
        },
    })
}
