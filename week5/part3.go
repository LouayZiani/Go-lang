package main

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	Field string
}

func (e *ValidationError) Error() string {
	return "invalid field: " + e.Field
}

func validate(age int) error {
	var err *ValidationError
	if age < 0 {
		err = &ValidationError{Field: "age"}
	}
	// there is a bug ehre: returning 'err' (a typed nil pointer) wraps it in an interface ( *ValidationError, nil ),
	// adn that is NOT a nil interface
	return err
}

func fixedValidate(age int) error {
	if age < 0 {
		return &ValidationError{Field: "age"}
	}
	return nil // this reutrns literal nil
}

// Task 4: Tour of Go -> Exercise: Errors

type MyError float64

func (e MyError) Error() string {
	// Infinite loop warning: fmt.Sprint(e) inside Error() calls e.String() or formatting on MyError, and that one calls Error() again recursively
	return fmt.Sprintf("myerror: %v", float64(e))
}

func Sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, MyError(x)
	}
	return 0, nil // simplified
}

func main121212() {
	// tsk 3 Demo
	err := validate(20)
	fmt.Printf("err == nil: %v\n", err == nil)
	fmt.Printf("err type: %T\n\n", err)

	// fixd version demo
	fixedErr := fixedValidate(-1)
	var valErr *ValidationError
	if errors.As(fixedErr, &valErr) {
		fmt.Printf("Caught expected validation error on field: %s\n", valErr.Field)
	}
}
