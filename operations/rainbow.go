package operations

import (
    "sync"
    "../core"
    "../utils"
    "reflect"
)

type RainbowConfig struct {
    Gradient map[string]float64
}

type Rainbow struct {
    sync.RWMutex

    name string
    stripe core.LEDStripe

    gradients *utils.GradientTable
}

func (this *Rainbow) Name() string {
    return this.name
}

func (this *Rainbow) Stripe() core.LEDStripe {
    return this.stripe
}

func (this *Rainbow) Render() {
    this.RLock()
    defer this.RUnlock()

    l := len(this.stripe)
    for i := range this.stripe {
        pos := float64(i) / float64(l)
        c := this.gradients.GetInterpolatedColorFor(pos)
        r, g, b := c.RGB255()

        this.stripe[i].R = uint8(r)
        this.stripe[i].G = uint8(g)
        this.stripe[i].B = uint8(b)
        this.stripe[i].A = 0.0
    }
}

func init() {
    core.RegisterOperation("rainbow", core.OperationFactory{
        ConfigType: reflect.TypeOf(struct{}{}),
        Create: func(pipeline *core.Pipeline, name string, count int, rconfig interface{}) (core.Operation, error) {
            //config := rconfig.(*RaindropConfig)

            keypoints := &utils.GradientTable{
                {utils.ParseColorHex("#ff0000"), 0.0},
                {utils.ParseColorHex("#d52a00"), 0.066},
                {utils.ParseColorHex("#ab5500"), 0.132},
                {utils.ParseColorHex("#ab7f00"), 0.198},
                {utils.ParseColorHex("#abab00"), 0.264},
                {utils.ParseColorHex("#56d500"), 0.330},
                {utils.ParseColorHex("#00ff00"), 0.396},
                {utils.ParseColorHex("#00d52a"), 0.462},
                {utils.ParseColorHex("#00AB55"), 0.528},
                {utils.ParseColorHex("#0056AA"), 0.594},
                {utils.ParseColorHex("#0000ff"), 0.660},
                {utils.ParseColorHex("#2a00d5"), 0.726},
                {utils.ParseColorHex("#5500ab"), 0.792},
                {utils.ParseColorHex("#7f0081"), 0.858},
                {utils.ParseColorHex("#ab0055"), 0.924},
                {utils.ParseColorHex("#ff0000"), 1.000},
            }

            return &Rainbow{
                name:      name,
                stripe: core.NewLEDStripe(count),
                gradients: keypoints,
            }, nil
        },
    })
}
