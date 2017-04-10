package utils

import (
    "fmt"
)

type BlendMode string

const (
    NORMAL      BlendMode = "normal"
    AVERAGE     BlendMode = "average"
    MULTIPLY    BlendMode = "multiply"
    SCREEN      BlendMode = "screen"
    DARKEN      BlendMode = "darken"
    LIGHTEN     BlendMode = "lighten"
    OVERLAY     BlendMode = "overlay"
    COLOR_DODGE BlendMode = "color_dodge"
    COLOR_BURN  BlendMode = "color_burn"
    ADDITIVE    BlendMode = "additive"
    SUBTRACTIVE BlendMode = "subtractive"
)

type BlendFunc func(float64, float64) float64

func BlendModeFunc(mode BlendMode) (BlendFunc, error) {
    switch mode {
    case NORMAL:
        return func(a float64, b float64) float64 {
            return b
        }, nil
    case AVERAGE:
        return func(a float64, b float64) float64 {
            return (a + b) / 2.0
        }, nil
    case MULTIPLY:
        return func(a float64, b float64) float64 {
            return a * b
        }, nil
    case SCREEN:
        return func(a float64, b float64) float64 {
            return 1.0 - (1.0-a)*(1.0-b)
        }, nil
    case DARKEN:
        return func(a float64, b float64) float64 {
            if a < b {
                return a
            } else {
                return b
            }
        }, nil
    case LIGHTEN:
        return func(a float64, b float64) float64 {
            if a > b {
                return a
            } else {
                return b
            }
        }, nil
    case OVERLAY:
        return func(a float64, b float64) float64 {
            if a < 0.5 {
                return 2 * a * b
            } else {
                return 1.0 - 2.0*(1.0-a)*(1.0-b)
            }
        }, nil
    case COLOR_DODGE:
        return func(a float64, b float64) float64 {
            return a / (1.0 - b)
        }, nil
    case COLOR_BURN:
        return func(a float64, b float64) float64 {
            return 1.0 - (1.0-a)/b
        }, nil
    case ADDITIVE:
        return func(a float64, b float64) float64 {
            return a + b
        }, nil
    case SUBTRACTIVE:
        return func(a float64, b float64) float64 {
            return a + b - 1.0
        }, nil
    default:
        return nil, fmt.Errorf("Invalid blend mode: %s", mode)
    }
}

func (f BlendFunc) Blend(a float64, b float64, blend float64) float64 {
    if blend < 0.0 {
        return (f(a, b) * (1.0 - (-blend))) + (a * (-blend))
    } else {
        return (f(a, b) * (1.0 - (+blend))) + (b * (+blend))
    }
}
