package core

import "image/color"

type LEDStripe []color.RGBA

func NewLEDStripe(count int) LEDStripe {
    stripe := make([]color.RGBA, count)

    return stripe
}

type LEDStripeReader interface {
    Count() int
    Get(i int) (r, g, b uint8)
}

func (this LEDStripe) Count() int {
    return len(this)
}

func (this LEDStripe) Get(i int) (r, g, b uint8) {
    return this[i].R, this[i].G, this[i].B
}

func (this LEDStripe) Set(i int, r, g, b uint8) {
    this[i].R = r
    this[i].G = g
    this[i].B = b
    this[i].A = 0xFF
}
