package task

import (
	"fmt"

	"contrastchecker/internal/contrast"
)


type Processor interface {
	Process(in Input) (Result, error)
}


type WCAGProcessor struct{}

func (WCAGProcessor) Process(in Input) (Result, error) {
	fg, err := contrast.ParseHex(in.Foreground)
	if err != nil {
		return Result{}, fmt.Errorf("foreground: %w", err)
	}
	bg, err := contrast.ParseHex(in.Background)
	if err != nil {
		return Result{}, fmt.Errorf("background: %w", err)
	}

	ratio := contrast.Ratio(fg, bg)
	large := contrast.IsLargeText(in.PointSize, in.Bold)
	aa, aaa := contrast.Thresholds(large)

	res := Result{
		Ratio:     ratio,
		LargeText: large,
		PassAA:    ratio >= aa,
		PassAAA:   ratio >= aaa,
	}
	switch {
	case res.PassAAA:
		res.Level = LevelAAA
	case res.PassAA:
		res.Level = LevelAA
	default:
		res.Level = LevelFail
	}
	return res, nil
}
