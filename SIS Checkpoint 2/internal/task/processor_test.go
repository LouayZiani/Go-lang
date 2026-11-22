package task

import (
	"math"
	"testing"
)

func TestWCAGProcessor(t *testing.T) {
	tests := []struct {
		name      string
		in        Input
		wantLevel Level
		wantAA    bool
		wantAAA   bool
		wantLarge bool
		wantRatio float64
	}{
		{"black on white", Input{"#000", "#fff", 12, false}, LevelAAA, true, true, false, 21},
		{"#767676 on white, normal", Input{"#767676", "#ffffff", 12, false}, LevelAA, true, false, false, 4.54},
		{"#777777 on white, normal fails", Input{"#777777", "#ffffff", 12, false}, LevelFail, false, false, false, 4.48},
		{"#777777 on white, large passes AA", Input{"#777777", "#ffffff", 18, false}, LevelAA, true, false, true, 4.48},
		{"#777777 on white, bold 14pt counts as large", Input{"#777777", "#ffffff", 14, true}, LevelAA, true, false, true, 4.48},
		{"same colour", Input{"#336699", "#336699", 24, true}, LevelFail, false, false, true, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WCAGProcessor{}.Process(tt.in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Level != tt.wantLevel || got.PassAA != tt.wantAA || got.PassAAA != tt.wantAAA || got.LargeText != tt.wantLarge {
				t.Errorf("got %+v", got)
			}
			if math.Abs(got.Ratio-tt.wantRatio) > 0.01 {
				t.Errorf("ratio = %.4f, want about %.2f", got.Ratio, tt.wantRatio)
			}
		})
	}
}

func TestWCAGProcessorBadColour(t *testing.T) {
	if _, err := (WCAGProcessor{}).Process(Input{"nope", "#fff", 12, false}); err == nil {
		t.Error("want error for bad foreground")
	}
	if _, err := (WCAGProcessor{}).Process(Input{"#fff", "nope", 12, false}); err == nil {
		t.Error("want error for bad background")
	}
}
