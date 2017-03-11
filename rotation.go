package main

import (
	"time"
	"image/color"
	"sync"
	"log"
)

type Rotation struct {
	sync.RWMutex
	StepsPerSecond float64 `mqtt:"speed"`
	LastFrameTime time.Time
	Offset float64
}

func rotateLEDs(leds []color.RGBA, k int)  {
    nlen := len(leds)
    if nlen <= 1 {
        return
    }

    k = k % nlen
    if k == 0 {
        return
    }

    for i := 0; i < k; i++ {
        for j := nlen - 1; j > 0; j-- {
            leds[j], leds[j-1] = leds[j-1], leds[j]
        }
    }
}

func (r *Rotation) Render(stripe *LEDStripe) {
	if r.LastFrameTime.Equal(time.Time{}) {
		r.LastFrameTime = time.Now()
	} else {
		now := time.Now()
		diff := now.Sub(r.LastFrameTime)
		ndiff := diff.Nanoseconds()

		timePerStep := time.Second / time.Duration(r.StepsPerSecond)

		r.Offset += float64(ndiff) / float64(timePerStep)
		r.LastFrameTime = now

		iOffset := int(r.Offset) % len(stripe.LEDS)
		rotateLEDs(stripe.LEDS, iOffset)
	}
}

func NewRotation(StepsPerSecond float64) Operation {
	s := &Rotation{
		StepsPerSecond: StepsPerSecond,
		LastFrameTime: time.Time{},
		Offset: 0.0,
	}

	return s
}
