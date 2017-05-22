package operations

import (
    "github.com/andir/lightsd/core"
    "github.com/andir/lightsd/operations"
    "sync"
    "reflect"
    "fmt"
    "github.com/andir/lightsd/utils"
)

type MixerConfig struct {
    Source1 string `mapstructure:"source1"`
    Source2 string `mapstructure:"source2"`

    Mode utils.BlendMode `mapstructure:"mode"`

    Speed float64 `mapstructure:"speed"`
}

type Mixer struct {
    sync.RWMutex

    name string

    source1 string
    source2 string

    function utils.BlendFunc

    blend float64
    speed float64

    Target float64 `mqtt:"target"`
}

type blendingLEDStripeReader struct {
    source1 core.LEDStripeReader
    source2 core.LEDStripeReader

    function utils.BlendFunc
    blend    float64
}

func (this *blendingLEDStripeReader) Get(i int) core.LED {
    a := this.source1.Get(i)
    b := this.source2.Get(i)

    return core.LED{
        R: this.function.Blend(a.R, b.R, this.blend),
        G: this.function.Blend(a.G, b.G, this.blend),
        B: this.function.Blend(a.B, b.B, this.blend),
    }
}

func (this *Mixer) Name() string {
    return this.name
}

func (this *Mixer) Render(context *core.RenderContext) core.LEDStripeReader {
    // TODO: Ouch, this hurts
    source1 := context.Results[this.source1]
    source2 := context.Results[this.source2]

    this.blend += (utils.ClampFloat64(this.Target, -1.0, 1.0) - this.blend) * context.Duration.Seconds() * this.speed

    return &blendingLEDStripeReader{
        source1: source1,
        source2: source2,

        function: this.function,
        blend:    this.blend,
    }
}

func init() {
    operations.Register("mixer", &operations.Factory{
        ConfigType: reflect.TypeOf(MixerConfig{}),
        Create: func(name string, count int, rconfig interface{}) (core.Operation, error) {
            config := rconfig.(*MixerConfig)

            function, err := utils.BlendModeFunc(config.Mode)
            if err != nil {
                return nil, fmt.Errorf("mixer: failed to parse blend mode: %v: %v", config.Mode, err)
            }

            return &Mixer{
                name: name,

                source1: config.Source1,
                source2: config.Source2,

                function: function,

                blend: -1.0,
                speed: config.Speed,

                Target: -1.0,
            }, nil
        },
    })
}
