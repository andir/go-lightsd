package main

import (
	"github.com/lucasb-eyer/go-colorful"
	"image/color"
	"sync"
)


type Rainbow struct {
	sync.RWMutex
	gradients *GradientTable
}


func NewRainbow() Operation {
	keypoints := &GradientTable{
		//{MustParseHex("#9e0142"), 0.0},
		//{MustParseHex("#d53e4f"), 0.1},
		//{MustParseHex("#f46d43"), 0.2},
		//{MustParseHex("#fdae61"), 0.3},
		//{MustParseHex("#fee090"), 0.4},
		//{MustParseHex("#ffffbf"), 0.5},
		//{MustParseHex("#e6f598"), 0.6},
		//{MustParseHex("#abdda4"), 0.7},
		//{MustParseHex("#66c2a5"), 0.8},
		//{MustParseHex("#3288bd"), 0.9},
		//{MustParseHex("#5e4fa2"), 1.0},


		{MustParseHex("#ff0000"), 0.0},
		{MustParseHex("#d52a00"), 0.066},
		{MustParseHex("#ab5500"), 0.132},
		{MustParseHex("#ab7f00"), 0.198},
		{MustParseHex("#abab00"), 0.264},
		{MustParseHex("#56d500"), 0.330},
		{MustParseHex("#00ff00"), 0.396},
		{MustParseHex("#00d52a"), 0.462},
		{MustParseHex("#00AB55"), 0.528},
		{MustParseHex("#0056AA"), 0.594},
		{MustParseHex("#0000ff"), 0.660},
		{MustParseHex("#2a00d5"), 0.726},
		{MustParseHex("#5500ab"), 0.792},
		{MustParseHex("#7f0081"), 0.858},
		{MustParseHex("#ab0055"), 0.924},
		{MustParseHex("#ff0000"), 1.000},
	}


	s := &Rainbow{
		gradients: keypoints,
	}
	return s
}

// This table contains the "keypoints" of the colorgradient you want to generate.
// The position of each keypoint has to live in the range [0,1]
type GradientTable []struct {
    Col colorful.Color
    Pos float64
}

// This is the meat of the gradient computation. It returns a HCL-blend between
// the two colors around `t`.
// Note: It relies heavily on the fact that the gradient keypoints are sorted.
func (self GradientTable) GetInterpolatedColorFor(t float64) colorful.Color {
    for i := 0 ; i < len(self) - 1 ; i++ {
        c1 := self[i]
        c2 := self[i+1]
        if c1.Pos <= t && t <= c2.Pos {
            // We are in between c1 and c2. Go blend them!
            t := (t - c1.Pos)/(c2.Pos - c1.Pos)
            return c1.Col.BlendHcl(c2.Col, t).Clamped()
        }
    }
    // Nothing found? Means we're at (or past) the last gradient keypoint.
    return self[len(self)-1].Col
}

// This is a very nice thing Golang forces you to do!
// It is necessary so that we can write out the literal of the colortable below.
func MustParseHex(s string) colorful.Color {
    c, err := colorful.Hex(s)
    if err != nil {
        panic("MustParseHex: " + err.Error())
    }
    return c
}


func (r *Rainbow) Render(stripe LEDStripe) {
	r.RLock()
	defer r.RUnlock()
	l := len(stripe)
	for i := range stripe {
		pos := float64(i)/float64(l)
		c := r.gradients.GetInterpolatedColorFor(pos)
		r,g,b := c.RGB255()
		stripe[i] = color.RGBA{
			R: uint8(r),
			G: uint8(g),
			B: uint8(b),
			A: 0.0,
		}
	}
}