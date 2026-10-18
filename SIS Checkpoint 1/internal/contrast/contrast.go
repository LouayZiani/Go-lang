// Package contrast implements the WCAG 2.x colour contrast maths.
// It knows nothing about tasks; it only turns colours into numbers.
package contrast

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// RGB is an 8-bit-per-channel sRGB colour.
type RGB struct {
	R, G, B uint8
}

// ParseHex parses "#RGB" or "#RRGGBB" (the leading '#' is optional,
// surrounding spaces are ignored).
func ParseHex(s string) (RGB, error) {
	h := strings.TrimPrefix(strings.TrimSpace(s), "#")
	if h == "" {
		return RGB{}, fmt.Errorf("colour is empty")
	}

	if len(h) == 3 { // "#f80" -> "#ff8800"
		h = string([]byte{h[0], h[0], h[1], h[1], h[2], h[2]})
	}
	if len(h) != 6 {
		return RGB{}, fmt.Errorf("colour %q must be #RGB or #RRGGBB", s)
	}

	var ch [3]uint8
	for i := range ch {
		v, err := strconv.ParseUint(h[i*2:i*2+2], 16, 8)
		if err != nil {
			return RGB{}, fmt.Errorf("colour %q contains non-hex digits", s)
		}
		ch[i] = uint8(v)
	}
	return RGB{R: ch[0], G: ch[1], B: ch[2]}, nil
}

// Luminance returns the relative luminance of c, from 0 (black) to 1 (white).
func (c RGB) Luminance() float64 {
	return 0.2126*linear(c.R) + 0.7152*linear(c.G) + 0.0722*linear(c.B)
}

// linear converts one sRGB channel to linear light.
func linear(v uint8) float64 {
	s := float64(v) / 255
	if s <= 0.04045 {
		return s / 12.92
	}
	return math.Pow((s+0.055)/1.055, 2.4)
}

// Ratio returns the WCAG contrast ratio between two colours, from 1 to 21.
// The order of the arguments does not matter.
func Ratio(a, b RGB) float64 {
	l1, l2 := a.Luminance(), b.Luminance()
	if l1 < l2 {
		l1, l2 = l2, l1
	}
	return (l1 + 0.05) / (l2 + 0.05)
}
