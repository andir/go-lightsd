package core

import "image/color"

type LEDStripe []color.RGBA

func NewLEDStripe(count int) LEDStripe {
    stripe := make([]color.RGBA, count)

    return stripe
}
