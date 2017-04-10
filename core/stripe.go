package core

import "github.com/andir/lightsd/utils"

type LED struct {
    R, G, B float64
}

func (this LED) RGB255() (r, g, b byte) {
    return byte(utils.Clamp(this.R, 0.0, 1.0) * 255),
        byte(utils.Clamp(this.G, 0.0, 1.0) * 255),
        byte(utils.Clamp(this.B, 0.0, 1.0) * 255)
}

type LEDStripe []LED

func NewLEDStripe(count int) LEDStripe {
    stripe := make([]LED, count)

    return stripe
}

type LEDStripeReader interface {
    Get(i int) LED
}

func (this LEDStripe) Get(i int) LED {
    return this[i]
}

func (this LEDStripe) Set(i int, led LED) {
    this[i].R = led.R
    this[i].G = led.G
    this[i].B = led.B
}
