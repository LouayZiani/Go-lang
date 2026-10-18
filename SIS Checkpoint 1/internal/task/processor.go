package task

import "sis/internal/contrast"

type Processor interface {
	Process(in Input) (Result, error)
}

// WCAGProcessor checks WCAG 2.x text contrast.
type WCAGProcessor struct{}

// Process implements Processor.
func (WCAGProcessor) Process(in Input) (Result, error) {
	fg, err := contrast.ParseHex(in.Foreground)
	if err != nil {
		return Result{}, err
	}
	bg, err := contrast.ParseHex(in.Background)
	if err != nil {
		return Result{}, err
	}

	ratio := contrast.Ratio(fg, bg)
	large := in.PointSize >= 18 || (in.Bold && in.PointSize >= 14)

	aa, aaa := 4.5, 7.0
	if large {
		aa, aaa = 3.0, 4.5
	}

	return Result{
		Ratio:     ratio,
		LargeText: large,
		AA:        ratio >= aa,
		AAA:       ratio >= aaa,
	}, nil
}
