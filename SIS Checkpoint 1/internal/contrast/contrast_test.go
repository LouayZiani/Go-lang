package contrast

import (
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
		{"six digits", "#FF8800", RGB{255, 136, 0}, false},
		{"lowercase", "#ff8800", RGB{255, 136, 0}, false},
		{"three digits expand", "#f80", RGB{255, 136, 0}, false},
		{"no hash", "000000", RGB{0, 0, 0}, false},
		{"surrounding spaces", "  #FFFFFF ", RGB{255, 255, 255}, false},
		{"empty", "", RGB{}, true},
		{"only hash", "#", RGB{}, true},
		{"whitespace only", "   ", RGB{}, true},
		{"too short", "#12", RGB{}, true},
		{"too long", "#1234567", RGB{}, true},
		{"non hex", "#GGGGGG", RGB{}, true},
		{"sign characters", "#+12-34", RGB{}, true},
		{"colour name", "red", RGB{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseHex(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseHex(%q) error = %v, wantErr %v", tt.in, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseHex(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestRatio(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want float64 // published WCAG values, to 2 decimals
	}{
		{"black on white is the maximum", "#000000", "#FFFFFF", 21.00},
		{"same colour is the minimum", "#336699", "#336699", 1.00},
		{"#767676 is the lightest grey that passes AA on white", "#767676", "#FFFFFF", 4.54},
		{"#777777 just misses AA on white", "#777777", "#FFFFFF", 4.48},
		{"pure red on white", "#FF0000", "#FFFFFF", 4.00},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, _ := ParseHex(tt.a)
			b, _ := ParseHex(tt.b)
			got := Ratio(a, b)
			if math.Abs(got-tt.want) > 0.005 {
				t.Errorf("Ratio(%s, %s) = %.4f, want %.2f", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestRatioIsSymmetric(t *testing.T) {
	a, _ := ParseHex("#1A73E8")
	b, _ := ParseHex("#F1F3F4")
	if Ratio(a, b) != Ratio(b, a) {
		t.Errorf("Ratio(a,b)=%v but Ratio(b,a)=%v", Ratio(a, b), Ratio(b, a))
	}
}
