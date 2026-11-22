package contrast

import (
	"errors"
	"math"
	"testing"
)

func TestParseHex(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    RGB
		wantErr bool
	}{
		{"long form", "#112233", RGB{0x11, 0x22, 0x33}, false},
		{"short form", "#abc", RGB{0xaa, 0xbb, 0xcc}, false},
		{"upper case", "#FFAA00", RGB{0xff, 0xaa, 0x00}, false},
		{"empty", "", RGB{}, true},
		{"no hash", "112233", RGB{}, true},
		{"only hash", "#", RGB{}, true},
		{"bad length 4", "#1234", RGB{}, true},
		{"bad length 7", "#1234567", RGB{}, true},
		{"non hex digit", "#12345g", RGB{}, true},
		{"sign char", "#+12345", RGB{}, true},
		{"spaces", "# 12 34", RGB{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseHex(tt.in)
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidColor) {
					t.Fatalf("want ErrInvalidColor, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestRatio(t *testing.T) {
	white := RGB{255, 255, 255}
	black := RGB{0, 0, 0}
	tests := []struct {
		name string
		a, b RGB
		want float64
	}{
		{"black on white", black, white, 21},
		{"order does not matter", white, black, 21},
		{"same colour", white, white, 1},
		{"#767676 on white", RGB{0x76, 0x76, 0x76}, white, 4.54},
		{"#777777 on white", RGB{0x77, 0x77, 0x77}, white, 4.48},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Ratio(tt.a, tt.b)
			if math.Abs(got-tt.want) > 0.01 {
				t.Fatalf("got %.4f, want about %.2f", got, tt.want)
			}
		})
	}
}

func TestIsLargeText(t *testing.T) {
	tests := []struct {
		size float64
		bold bool
		want bool
	}{
		{12, false, false},
		{14, false, false},
		{14, true, true},
		{13.9, true, false},
		{18, false, true},
		{17.9, false, false},
		{24, true, true},
	}
	for _, tt := range tests {
		if got := IsLargeText(tt.size, tt.bold); got != tt.want {
			t.Errorf("IsLargeText(%v, %v) = %v, want %v", tt.size, tt.bold, got, tt.want)
		}
	}
}

func TestThresholds(t *testing.T) {
	if aa, aaa := Thresholds(false); aa != 4.5 || aaa != 7 {
		t.Errorf("normal: got %v/%v", aa, aaa)
	}
	if aa, aaa := Thresholds(true); aa != 3 || aaa != 4.5 {
		t.Errorf("large: got %v/%v", aa, aaa)
	}
}
