package operations

import (
    "math/rand"
    "sync"
    "github.com/andir/lightsd/core"
    "github.com/andir/lightsd/operations"
    "time"
    "reflect"
    "github.com/lucasb-eyer/go-colorful"
)

type RaindropConfig struct {
    HueMin float64 `mapstructure:"hue_min"`
    HueMax float64 `mapstructure:"hue_max"`

    SaturationMin float64 `mapstructure:"sat_min"`
    SaturationMax float64 `mapstructure:"sat_max"`

    ValueMin float64 `mapstructure:"val_min"`
    ValueMax float64 `mapstructure:"val_max"`

    DecayLow  float64 `mapstructure:"decay_low"`
    DecayHigh float64 `mapstructure:"decay_high"`

    Chance float64 `mapstructure:"chance"`
}

type Raindrop struct {
    sync.RWMutex

    name   string

    HueMin float64 `mqtt:"hue_min"`
    HueMax float64 `mqtt:"hue_max"`

    SaturationMin float64 `mqtt:"sat_min"`
    SaturationMax float64 `mqtt:"sat_max"`

    ValueMin float64 `mqtt:"val_min"`
    ValueMax float64 `mqtt:"val_max"`

    DecayLow  float64 `mqtt:"decay_low"`
    DecayHigh float64 `mqtt:"decay_high"`

    Chance float64 `mqtt:"chance"`

    rand *rand.Rand

    leds []RaindropLED
}

type RaindropLED struct {
    Color     colorful.Color
    DecayRate float64
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
    return (ra.Float64() * (max - min)) + min
}

func (this *Raindrop) Name() string {
    return this.name
}

func (this *Raindrop) Render(context *core.RenderContext) {
    for i, l := range this.leds {
        roll := randomFloat64(this.rand, 0.0, 1.0)

        if roll > this.Chance {
            saturation := randomFloat64(this.rand, this.SaturationMin, this.SaturationMax)
            hue := randomFloat64(this.rand, this.HueMin, this.HueMax)
            value := randomFloat64(this.rand, this.ValueMin, this.ValueMax)

            decayRate := randomFloat64(this.rand, this.DecayLow, this.DecayHigh)

            l.Color = colorful.Hsv(hue, saturation, value)
            l.DecayRate = 1.0 - decayRate
        }

        factor := maxFloat64(minFloat64(l.DecayRate, 1.0), 0.0)
        h, s, v := l.Color.Hsv()
        v *= float64(factor)
        l.Color = colorful.Hsv(h, s, v)

        r, g, b := l.Color.RGB255()
        context.Set(i, r, g, b)
    }
}

func init() {
    operations.Register("raindrops", &operations.Factory{
        ConfigType: reflect.TypeOf(RaindropConfig{}),
        Create: func(name string, count int, rconfig interface{}) (core.Operation, error) {
            config := rconfig.(*RaindropConfig)

            return &Raindrop{
                name:   name,

                HueMin: config.HueMin,
                HueMax: config.HueMax,

                SaturationMin: config.SaturationMin,
                SaturationMax: config.SaturationMax,

                ValueMin: config.ValueMin,
                ValueMax: config.ValueMax,

                DecayLow:  config.DecayLow,
                DecayHigh: config.DecayHigh,

                Chance: config.Chance,

                rand: rand.New(rand.NewSource(time.Now().Unix())),
                leds: make([]RaindropLED, count),
            }, nil
        },
    })
}
