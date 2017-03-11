package main

import (
	"math/rand"
	"github.com/lucasb-eyer/go-colorful"
	"image/color"
	"log"
	"sync"
)

type Raindrop struct {
	sync.RWMutex

	HueMin float64 `mqtt:"hue_min"`
	HueMax float64 `mqtt:"hue_max"`

	SaturationMin float64 `mqtt:"saturation_min"`
	SaturationMax float64 `mqtt:"saturation_max"`

	ValueMin float64 `mqtt:"value_min"`
	ValueMax float64 `mqtt:"value_max"`

	Chance float64 `mqtt:"chance"`

	DecayLow float64 `mqtt:"decay_low"`
	DecayHigh float64 `mqtt:"decay_high"`

	leds []RaindropLED
}

type RaindropLED struct {
	Color colorful.Color
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

func randomFloat64(min, max float64) float64 {

	diff := max - min

	return (rand.Float64() * diff) + min
}

func (r *RaindropLED) Decay() {
	factor := maxFloat64(minFloat64(r.DecayRate, 1.0), 0.0)
	h, s, v := r.Color.Hsv()
	v *= float64(factor)
	r.Color = colorful.Hsv(h, s, v)
}


func (r *Raindrop) HitLED(led *RaindropLED) {
	log.Println("hit")
	saturation := randomFloat64(r.SaturationMin, r.SaturationMax)
	hue := randomFloat64(r.HueMin, r.HueMax)
	value := randomFloat64(r.ValueMin, r.ValueMax)

	decayRate := randomFloat64(r.DecayLow, r.DecayHigh)

	//log.Printf("H: %v S: %v V: %v, R: %v", hue, saturation, value, decayRate)
	led.Color = colorful.Hsv(hue, saturation, value)
	led.DecayRate = 1.0 - decayRate
}

func NewRaindrop() Operation {
	return &Raindrop{
		HueMin: 0.0,
		HueMax: 360.0,

		SaturationMin: 0.0,
		SaturationMax: 1.0,

		ValueMin: 0.0,
		ValueMax: 1.0,

		Chance: 0.95,

		DecayLow: 0.001,
		DecayHigh: 0.5,
	}
}

func (r *Raindrop) Render(stripe *LEDStripe) {
	r.RLock()
	defer r.RUnlock()

	if r.leds == nil || len(r.leds) != len(stripe.LEDS) {
		r.leds = make([]RaindropLED, len(stripe.LEDS))
	}

	for i := range stripe.LEDS {
		roll := randomFloat64(0.0, 1.0)

		l := &r.leds[i]
		if roll > r.Chance {
			r.HitLED(l)

		}
		l.Decay()

		r, g, b := l.Color.RGB255()
		stripe.LEDS[i] = color.RGBA{R:r, G:g, B:b, A:0}
	}
}