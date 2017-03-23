package utils

import (
    "github.com/lucasb-eyer/go-colorful"
    "github.com/andir/lightsd/core"
)

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
    for i := 0; i < len(self)-1; i++ {
        c1 := self[i]
        c2 := self[i+1]
        if c1.Pos <= t && t <= c2.Pos {
            // We are in between c1 and c2. Go blend them!
            t := (t - c1.Pos) / (c2.Pos - c1.Pos)
            return c1.Col.BlendHcl(c2.Col, t).Clamped()
        }
    }
    // Nothing found? Means we're at (or past) the last gradient keypoint.
    return self[len(self)-1].Col
}

func (self GradientTable) Fill(stripe core.LEDStripe) {
    for i := 0; i < stripe.Count(); i++ {
        pos := float64(i) / float64(stripe.Count()-1)

        r, g, b := self.GetInterpolatedColorFor(pos).RGB255()
        stripe.Set(i, r, g, b, )
    }
}
