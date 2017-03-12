package utils

import "github.com/lucasb-eyer/go-colorful"

func ParseColorHex(s string) colorful.Color {
    c, err := colorful.Hex(s)
    if err != nil {
        panic("Failed to parse color: " + err.Error())
    }
    return c
}
