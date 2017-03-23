package operations

import (
    "github.com/andir/lightsd/core"
    "github.com/andir/lightsd/utils"
    "github.com/andir/lightsd/operations"
    "reflect"
    "sync"
)

type GradientConfig struct {
    Gradient []struct {
        C string
        P float64
    }
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
                gradient[i].Col = utils.ParseColorHex(e.C)
                gradient[i].Pos = e.P
            }

            stripe := core.NewLEDStripe(count)
            gradient.Fill(stripe)

            return &Gradient{
                name:   name,
                stripe: stripe,
            }, nil
        },
    })
}
