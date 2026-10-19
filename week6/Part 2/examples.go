package main

import (
	"errors"
	"fmt"
	"strconv"
)

var ErrNotFound = errors.New("student not found")

type GradeError struct {
	Student string
	Grade   int
}

func (e *GradeError) Error() string {
	return fmt.Sprintf("invalid grade %d for %s", e.Grade, e.Student)
}

func GetGrade(name string) (int, error) {
	if name != "Ali" {
		return 0, ErrNotFound
	}
	return 85, nil
}

func SetGrade(name string, grade int) error {
	if grade < 0 || grade > 100 {
		return &GradeError{Student: name, Grade: grade}
	}
	fmt.Println("Grade saved:", name, grade)
	return nil
}

func main() {
	n, err := strconv.Atoi("42a")
	fmt.Println(n, err)

	grade, err := GetGrade("Ali")
	fmt.Println(grade, err)

	_, err = GetGrade("Timur")
	fmt.Println(err)

	err = SetGrade("Ali", 120)
	fmt.Println(err)

	err = SetGrade("Ali", 90)
	fmt.Println(err)
}
