package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrBadFormat = errors.New("bad record format")

type Record struct {
	Name  string
	Grade int
}

func ParseAll(input string) ([]Record, error) {
	var records []Record
	var errs []error

	lines := strings.Split(input, "\n")

	for i, line := range lines {
		lineNumber := i + 1
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		name, gradeText, ok := strings.Cut(line, ",")
		if !ok {
			errs = append(errs, fmt.Errorf(
				"line %d: %w", lineNumber, ErrBadFormat,
			))
			continue
		}

		name = strings.TrimSpace(name)
		gradeText = strings.TrimSpace(gradeText)

		grade, err := strconv.Atoi(gradeText)
		if err != nil {
			errs = append(errs, fmt.Errorf(
				"line %d: parse grade: %w", lineNumber, err,
			))
			continue
		}

		records = append(records, Record{
			Name:  name,
			Grade: grade,
		})
	}

	return records, errors.Join(errs...)
}

func main() {
	input := `Aruzhan, 95
Dias, 8O
Timur
Bota, 99999999999999999999
Ali, 70`

	records, err := ParseAll(input)

	fmt.Println("Records:", records)
	fmt.Println("Error:", err)

	fmt.Println("Bad format:", errors.Is(err, ErrBadFormat))
	fmt.Println("Syntax:", errors.Is(err, strconv.ErrSyntax))
	fmt.Println("Range:", errors.Is(err, strconv.ErrRange))

	var numErr *strconv.NumError
	if errors.As(err, &numErr) {
		fmt.Printf(
			"First numeric error: function=%s, input=%q, error=%v\n",
			numErr.Func,
			numErr.Num,
			numErr.Err,
		)
	}
}
