package main

import (
	"time"
	"sync"
)

type Rotation struct {
	name string

	sync.RWMutex

	StepsPerSecond float64 `mqtt:"speed"`
	LastFrameTime time.Time

	Offset float64
}

func (r *Rotation) Render(stripe LEDStripe) LEDStripe {
	now := time.Now()

	delta := now.Sub(r.LastFrameTime)

	r.Offset += delta.Seconds() * r.StepsPerSecond
	r.LastFrameTime = now

	iOffset := int(r.Offset) % len(stripe)

	if iOffset == 0 {
		return stripe
	}

	output := NewLEDStripe(len(stripe))
	for i, s := range(stripe) {
		output[(i + iOffset) % len(stripe)] = s
	}

	return output
}

func NewRotation(name string, StepsPerSecond float64) Operation {
	s := &Rotation{
		name: name,

		StepsPerSecond: StepsPerSecond,
		LastFrameTime: time.Time{},
		Offset: 0.0,
	}

	return s
}

func (r *Rotation) Name() string {
	return r.name
}
