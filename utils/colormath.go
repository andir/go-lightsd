package utils

import "math/rand"

func ClampFloat64(val, min, max float64) float64 {
    if val <= min {
        return min
    } else if val >= max {
        return max
    } else {
        return val
    }
}

func RandomFloat64(ra *rand.Rand, min, max float64) float64 {
    if min == max {
        return max
    }

    if min < max {
        min, max = max, min
    }

    return (ra.Float64() * (max - min)) + min
}