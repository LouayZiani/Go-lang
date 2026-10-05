// Task 1: Student records
// 1. Declare the Student struct above and create students in 4 different
// ways (zero value, named fields, positional, pointer). Print each one with
// %+v .
// 2. Write func addCourse(s Student, course string) that appends a course to
// s.Courses . Call it and print the student. Did the course appear? Why?
// 3. Change the function to take *Student and call it again. What changed?
package main

import "fmt"

type Student struct {
	Name    string
	ID      int
	GPA     float64
	Courses []string
}

// part 2: value parameter
func addCourse(s Student, course string) {
	s.Courses = append(s.Courses, course)
}

// part 3: pointer parameter
func addCoursePtr(s *Student, course string) {
	s.Courses = append(s.Courses, course)
}

func main() {
	// 1. four ways
	var s1 Student                                             // zero value
	s2 := Student{Name: "Aruzhan", ID: 1, GPA: 3.8}            // named fields
	s3 := Student{"Dias", 2, 3.2, []string{"Go", "Databases"}} // positional
	s4 := &Student{Name: "Timur", ID: 3, GPA: 3.5}             // pointer

	fmt.Printf("%+v\n", s1) // {Name: ID:0 GPA:0 Courses:[]}
	fmt.Printf("%+v\n", s2) // {Name:Aruzhan ID:1 GPA:3.8 Courses:[]}
	fmt.Printf("%+v\n", s3) // {Name:Dias ID:2 GPA:3.2 Courses:[Go Databases]}
	fmt.Printf("%+v\n", s4) // &{Name:Timur ID:3 GPA:3.5 Courses:[]}  (note the & for a pointer)

	// 2. by value
	addCourse(s2, "Algorithms")
	fmt.Printf("after addCourse (value):   %+v\n", s2) // Courses:[] -> NOT added

	// 3. by pointer
	addCoursePtr(&s2, "Algorithms")
	fmt.Printf("after addCourse (pointer): %+v\n", s2) // Courses:[Algorithms] -> added
}

// 2. The course did NOT appear. s is a copy of the struct. append works on the copy's
//    Courses field (its own slice header), so the original's header (len 0) never changes.
//    Here Courses was nil, so append allocated a brand new array anyway. Even when there is spare
//    capacity and append writes into the shared array, the original's len stays the same, so you still wouldn't see it.
// 3. With *Student the function gets the address of the original, so s.Courses = append(...) assigns
//    to the real field and the course appears. (s.Courses is shorthand for (*s).Courses.)
