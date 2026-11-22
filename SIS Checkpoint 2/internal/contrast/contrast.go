// Package contrast holds the pure WCAG 2.x colour maths: parsing hex colours, relative luminance, contrast ratio and the AA / AAA thresholds.
// also nothing in here has state, so it is safe to call from any goroutine.
package contrast

import (
	"errors"
	"fmt"
	"math"
)

// ErrInvalidColor is returned (wrapped) when a colour string can't be parsed
var ErrInvalidColor = errors.New("invalid colour")

// RGB is an 8-bit-per-channel colour
type RGB struct {
	R, G, B uint8
}

// ParseHex parses "#RGB" or "#RRGGBB"
func ParseHex(s string) (RGB, error) {
	if s == "" {
		return RGB{}, fmt.Errorf("%w: empty", ErrInvalidColor)
	}
	if s[0] != '#' {
		return RGB{}, fmt.Errorf("%w: %q must start with '#'", ErrInvalidColor, s)
	}
	digits := s[1:]
	if len(digits) != 3 && len(digits) != 6 {
		return RGB{}, fmt.Errorf("%w: %q must be #RGB or #RRGGBB", ErrInvalidColor, s)
	}

	vals := make([]uint8, len(digits))
	for i := 0; i < len(digits); i++ {
		v, ok := hexValue(digits[i])
		if !ok {
			return RGB{}, fmt.Errorf("%w: %q has non-hex digit %q", ErrInvalidColor, s, digits[i])
		}
		vals[i] = v
	}

	if len(vals) == 3 {
		//#RGB is shorthand for #RRGGBB: each digit is doubled
		return RGB{R: vals[0]*16 + vals[0], G: vals[1]*16 + vals[1], B: vals[2]*16 + vals[2]}, nil
	}
	return RGB{R: vals[0]*16 + vals[1], G: vals[2]*16 + vals[3], B: vals[4]*16 + vals[5]}, nil
}

func hexValue(c byte) (uint8, bool) {
	switch {
	case c >= '0' && c <= '9':
		return c - '0', true
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10, true
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10, true
	}
	return 0, false
}

// Luminance returns the WCAG relative luminance of c (0 = black, 1 = white)
func Luminance(c RGB) float64 {
	return 0.2126*linear(c.R) + 0.7152*linear(c.G) + 0.0722*linear(c.B)
}

func linear(v uint8) float64 {
	s := float64(v) / 255
	if s <= 0.03928 {
		return s / 12.92
	}
	return math.Pow((s+0.055)/1.055, 2.4)
}

// Ratio returns the contrast ratio between two colours, in the range 1..21.
// and the order of the arguments doesn't matter.
func Ratio(a, b RGB) float64 {
	la, lb := Luminance(a), Luminance(b)
	if la < lb {
		la, lb = lb, la
	}
	return (la + 0.05) / (lb + 0.05)
}

// IsLargeText reports whether WCAG treats the text as "large": at least 18pt or at least 14pt when bold
func IsLargeText(pointSize float64, bold bool) bool {
	return pointSize >= 18 || (bold && pointSize >= 14)
}

// thresholds returns the minimum ratios for AA and AAA
func Thresholds(large bool) (aa, aaa float64) {
	if large {
		return 3, 4.5
	}
	return 4.5, 7
}
