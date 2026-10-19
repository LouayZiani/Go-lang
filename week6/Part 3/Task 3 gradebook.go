package main

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("student not found")

type GradeError struct {
	Student string
	Grade   int
}

func (e *GradeError) Error() string {
	return fmt.Sprintf(
		"invalid grade %d for %s",
		e.Grade,
		e.Student,
	)
}

type GradeBook struct {
	grades map[string]int
}

func NewGradeBook() *GradeBook {
	return &GradeBook{
		grades: make(map[string]int),
	}
}

func (gb *GradeBook) Get(name string) (int, error) {
	grade, ok := gb.grades[name]
	if !ok {
		return 0, fmt.Errorf("%s: %w", name, ErrNotFound)
	}
	return grade, nil
}

func (gb *GradeBook) Set(name string, grade int) error {
	if grade < 0 || grade > 100 {
		return &GradeError{
			Student: name,
			Grade:   grade,
		}
	}

	gb.grades[name] = grade
	return nil
}

func AddBonus(gb *GradeBook, name string, bonus int) error {
	grade, err := gb.Get(name)
	if err != nil {
		return fmt.Errorf("get grade: %w", err)
	}

	if err := gb.Set(name, grade+bonus); err != nil {
		return fmt.Errorf("add bonus for %s: %w", name, err)
	}

	return nil
}

func report(err error) {
	if err == nil {
		fmt.Println("OK")
		return
	}

	if errors.Is(err, ErrNotFound) {
		fmt.Println("404:", err)
		return
	}

	var ge *GradeError
	if errors.As(err, &ge) {
		fmt.Printf(
			"400: %s (student=%s, grade=%d)\n",
			err,
			ge.Student,
			ge.Grade,
		)
		return
	}

	fmt.Println("500:", err)
}

func main12() {
	gb := NewGradeBook()

	_ = gb.Set("Leia", 90)
	_ = gb.Set("Jakey", 98)

	report(AddBonus(gb, "Leia", 5))
	report(AddBonus(gb, "Jakey", 5))
	report(AddBonus(gb, "Louay", 5))

	// if we replace %w with %v: the not-found err no longer remains discoverable through the error chain. so,
	// errors.Is(err, ErrNotFound) returns false, and the reporting function falls through to the 500 branch.

	// The error's printed text may look the same but %v copies the msg rather than wrapping the og error
	// %w preserves the err relationship nd %v does not.
}
