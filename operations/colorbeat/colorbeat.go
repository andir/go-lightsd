package operations

import (
    "sync"
    "time"
    "reflect"
    "math/rand"
    "github.com/lucasb-eyer/go-colorful"
    "github.com/andir/lightsd/core"
    "github.com/andir/lightsd/operations"
    "github.com/andir/lightsd/utils"
    "math"
)

type ColorbeatConfig struct {
    SpeedMin float64 `mapstructure:"speed_min"`
    SpeedMax float64 `mapstructure:"speed_max"`

    PeakChange    float64 `mapstructure:"Peak_chance"`
    PeakThreshold float64 `mapstructure:"Peak_threshold"`
    PeakDecay     float64 `mapstructure:"Peak_decay"`
}

type Colorbeat struct {
    sync.RWMutex

    name   string
    stripe core.LEDStripe

    Speed    float64 `mqtt:"speed"`
    SpeedMin float64
    SpeedMax float64

    PeakChange    float64 `mqtt:"peak_chance"`
    PeakThreshold float64 `mqtt:"peak_threshold"`
    PeakDecay     float64 `mqtt:"peak_decay"`

    rand *rand.Rand

    leds colorbeatLEDStripe
}

type colorbeatLEDStripe []struct {
    hueValue float64
    hueSpeed float64

    peak float64
}

func (this colorbeatLEDStripe) Count() int {
    return len(this)
}

func (this colorbeatLEDStripe) Get(i int) core.LED {
    c := colorful.Hcl(
        this[i].hueValue,
        0.5,
        0.5+math.Pow(this[i].peak, 1.25)*0.5).Clamped()
    return core.LED{R: c.R, G: c.G, B: c.B}
}

func (this *Colorbeat) Name() string {
    return this.name
}

func (this *Colorbeat) Render(context *core.RenderContext) core.LEDStripeReader {
    beat := this.rand.Float64() < this.PeakThreshold*context.Duration.Seconds() // TODO: Beat detection

    for i := range this.leds {
        if this.leds[i].hueSpeed == 0.0 || (beat && utils.RandomFloat64(this.rand, 0.0, 1.0) < this.PeakChange) {
            this.leds[i].peak = 1.0
            this.leds[i].hueValue = utils.RandomFloat64(this.rand, 0.0, 360.0)
            this.leds[i].hueSpeed = RandomInverse(this.rand, utils.RandomFloat64(this.rand, this.SpeedMin, this.SpeedMax))

        } else {
            this.leds[i].peak -= this.PeakDecay * context.Duration.Seconds()
            this.leds[i].peak = math.Max(0.0, this.leds[i].peak)

            this.leds[i].hueValue += this.leds[i].hueSpeed * this.Speed * context.Duration.Seconds()
        }
    }

    return this.leds
}

func init() {
    operations.Register("colorbeat", &operations.Factory{
        ConfigType: reflect.TypeOf(ColorbeatConfig{}),
        Create: func(name string, count int, rconfig interface{}) (core.Operation, error) {
            config := rconfig.(*ColorbeatConfig)

            println(config.PeakThreshold)

            return &Colorbeat{
                name:   name,
                stripe: core.NewLEDStripe(count),

                Speed:    1.0,
                SpeedMin: config.SpeedMin,
                SpeedMax: config.SpeedMax,

                PeakChange:    config.PeakChange,
                PeakThreshold: config.PeakThreshold,
                PeakDecay:     config.PeakDecay,

                rand: rand.New(rand.NewSource(time.Now().Unix())),
                leds: make(colorbeatLEDStripe, count),
            }, nil
        },
    })
}

func RandomInverse(ra *rand.Rand, val float64) float64 {
    if ra.Uint32()%2 == 0 {
        return val
    } else {
        return -val
    }
}
