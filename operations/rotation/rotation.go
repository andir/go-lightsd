package operations

import (
    "github.com/andir/lightsd/core"
    "github.com/andir/lightsd/operations"
    "sync"
    "reflect"
)

type RotationConfig struct {
    Source         string `mapstructure:"source"`
    PixelPerSecond float64 `mapstructure:"speed"`
}

type Rotation struct {
    sync.RWMutex

    name string

    Source         string
    PixelPerSecond float64  `mqtt:"speed"`

    offset float64
}

type rotatedLEDStripeReader struct {
    source core.LEDStripeReader
    offset float64
}

func (this *rotatedLEDStripeReader) Count() int {
    return this.source.Count()
}

func(this *rotatedLEDStripeReader) Get(i int) (r, g, b uint8) {
    // TODO: Blending between colors?
    return this.source.Get((i+int(this.offset))%this.source.Count())
}

func (this *Rotation) Name() string {
    return this.name
}

func (this *Rotation) Render(context *core.RenderContext) core.LEDStripeReader {
    this.offset += context.Duration.Seconds() * this.PixelPerSecond

    // TODO: Ouch, this hurts
    source := context.Results[this.Source]

    return &rotatedLEDStripeReader{
        source: source,
        offset: this.offset,
    }
}

func init() {
    operations.Register("rotation", &operations.Factory{
        ConfigType: reflect.TypeOf(RotationConfig{}),
        Create: func(name string, count int, rconfig interface{}) (core.Operation, error) {
            config := rconfig.(*RotationConfig)

            //source := pipeline.ByName(config.Source)
            //if source == nil {
            //    return nil, fmt.Errorf("Unknown source: %s", config.Source)
            //}

            return &Rotation{
                name: name,

                Source:         config.Source,
                PixelPerSecond: config.PixelPerSecond,

                offset: 0.0,
            }, nil
        },
    })
}
