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

    leds raindropLEDStripe
}

type raindropLEDStripe []struct {
    color core.LED
    decay float64
}

func (this raindropLEDStripe) Count() int {
    return len(this)
}

func (this raindropLEDStripe) Get(i int) core.LED {
    return this[i].color
}

func (this *Raindrop) Name() string {
    return this.name
}

func (this *Raindrop) Render(context *core.RenderContext) core.LEDStripeReader {
    for i := range this.leds {
        roll := randomFloat64(this.rand, 0.0, 1.0)

        if roll < this.Chance {
            hue := randomFloat64(this.rand, this.HueMin, this.HueMax)
            saturation := randomFloat64(this.rand, this.SaturationMin, this.SaturationMax)
            value := randomFloat64(this.rand, this.ValueMin, this.ValueMax)

            decay := 1.0 / randomFloat64(this.rand, this.DecayMin, this.DecayMax)

            color := colorful.Hsv(hue, saturation, value)

            this.leds[i].color = core.LED{R: color.R, G: color.G, B: color.B}
            this.leds[i].decay = decay

        } else {
            v := 1.0 - this.leds[i].decay*context.Duration.Seconds()

            this.leds[i].color.R *= v
            this.leds[i].color.G *= v
            this.leds[i].color.B *= v
        }
    }

    return this.leds
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
                leds: make(raindropLEDStripe, count),
            }, nil
        },
    })
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
