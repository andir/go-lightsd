package utils

func Clamp(val, min, max float64) float64 {
    if val <= min {
        return min
    } else if val >= max {
        return max
    } else {
        return val
    }
}
