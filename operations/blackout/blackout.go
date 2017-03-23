package operations

import (
    "github.com/andir/lightsd/core"
    "github.com/andir/lightsd/operations"
    "sync"
    "reflect"
)

type BlackoutConfig struct {
    Source string `mapstructure:"source"`

    Enabled bool
    From    int
    To      int
}

type Blackout struct {
    sync.RWMutex

    name string

    source string

    Enabled bool  `mqtt:"enabled"`

    from   int
    to     int
    invert bool
}

type blackoutLEDStripeReader struct {
    source  core.LEDStripeReader
    enabled bool
    from    int
    to      int
    invert  bool
}

func (this *blackoutLEDStripeReader) Count() int {
    return this.source.Count()
}

func (this *blackoutLEDStripeReader) Get(i int) (r, g, b uint8) {
    if this.enabled && ((this.from <= i && i <= this.to) == this.invert) {
        return 0.0, 0.0, 0.0
    }

    return this.source.Get(i)
}

func (this *Blackout) Name() string {
    return this.name
}

func (this *Blackout) Render(context *core.RenderContext) core.LEDStripeReader {
    // TODO: Ouch, this hurts
    source := context.Results[this.source]

    return &blackoutLEDStripeReader{
        source:  source,
        enabled: this.Enabled,
        from:    this.from,
        to:      this.to,
    }
}

func init() {
    operations.Register("blackout", &operations.Factory{
        ConfigType: reflect.TypeOf(BlackoutConfig{}),
        Create: func(name string, count int, rconfig interface{}) (core.Operation, error) {
            config := rconfig.(*BlackoutConfig)

            from, to, invert := config.From, config.To, false
            if from > to {
                from, to, invert = to, from, true
            }

            return &Blackout{
                name: name,

                source: config.Source,

                Enabled: config.Enabled,

                from:   from,
                to:     to,
                invert: invert,
            }, nil
        },
    })
}
