package operations

import (
    "math/rand"
    "sync"
    "github.com/andir/lightsd/core"
    "github.com/andir/lightsd/operations"
    "time"
    "reflect"
    "github.com/lucasb-eyer/go-colorful"
    "fmt"
)

type RaindropConfig struct {
    HueMin float64 `mapstructure:"hue_min"`
    HueMax float64 `mapstructure:"hue_max"`

    SaturationMin float64 `mapstructure:"sat_min"`
    SaturationMax float64 `mapstructure:"sat_max"`

    ValueMin float64 `mapstructure:"val_min"`
    ValueMax float64 `mapstructure:"val_max"`

    DecayMin float64 `mapstructure:"decay_min"`
    DecayMax float64 `mapstructure:"decay_max"`

    Chance float64 `mapstructure:"chance"`
}

type Raindrop struct {
    sync.RWMutex

    name string

    HueMin float64 `mqtt:"hue_min"`
    HueMax float64 `mqtt:"hue_max"`

    SaturationMin float64 `mqtt:"sat_min"`
    SaturationMax float64 `mqtt:"sat_max"`

    ValueMin float64 `mqtt:"val_min"`
    ValueMax float64 `mqtt:"val_max"`

    DecayMin float64 `mqtt:"decay_min"`
    DecayMax float64 `mqtt:"decay_max"`

    Chance float64 `mqtt:"chance"`

    rand *rand.Rand

    leds []raindropLED

    stripe core.LEDStripe
}

type raindropLED struct {
    color colorful.Color
    decay float64
}

func maxFloat64(a, b float64) float64 {
    if a > b {
        return a
    } else {
        return b
    }
}

func minFloat64(a, b float64) float64 {
    if a < b {
        return a
    } else {
        return b
    }
}

func randomFloat64(ra *rand.Rand, min, max float64) float64 {
    if min == max {
        return max
    }

    if min < max {
        min, max = max, min
    }

    return (ra.Float64() * (max - min)) + min
}

func (this *Raindrop) Name() string {
    return this.name
}

func (this *Raindrop) Render(context *core.RenderContext) core.LEDStripeReader {
    for i, l := range this.leds {
        roll := randomFloat64(this.rand, 0.0, 1.0)

        if roll < this.Chance {
            hue := randomFloat64(this.rand, this.HueMin, this.HueMax)
            saturation := randomFloat64(this.rand, this.SaturationMin, this.SaturationMax)
            value := randomFloat64(this.rand, this.ValueMin, this.ValueMax)
            decay := randomFloat64(this.rand, this.DecayMin, this.DecayMax)

            this.leds[i] = raindropLED{
                color: colorful.Hsv(hue, saturation, value),
                decay: decay,
            }

        } else {
            h, s, v := l.color.Hsv()
            v *= 1.0 - (1.0 / l.decay) * context.Duration.Seconds()
            this.leds[i].color = colorful.Hsv(h, s, v)
        }

        r, g, b := l.color.RGB255()
        this.stripe.Set(i, r, g, b)
    }

    return this.stripe
}

func init() {
    operations.Register("raindrops", &operations.Factory{
        ConfigType: reflect.TypeOf(RaindropConfig{}),
        Create: func(name string, count int, rconfig interface{}) (core.Operation, error) {
            config := rconfig.(*RaindropConfig)

            return &Raindrop{
                name: name,

                HueMin: config.HueMin,
                HueMax: config.HueMax,

                SaturationMin: config.SaturationMin,
                SaturationMax: config.SaturationMax,

                ValueMin: config.ValueMin,
                ValueMax: config.ValueMax,

                DecayMin: config.DecayMin,
                DecayMax: config.DecayMax,

                Chance: config.Chance,

                rand: rand.New(rand.NewSource(time.Now().Unix())),
                leds: make([]raindropLED, count),

                stripe: core.NewLEDStripe(count),
            }, nil
        },
    })
}
