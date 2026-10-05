// Task 2: University
// 1. Create Person with fields Name , Email and a method Contact() string
// that returns "Name <Email>" .
// 2. Create Student that embeds Person and adds GPA . Create Teacher that
// embeds Person and adds Department .
// 3. Give Teacher its own Contact() that also includes the department. Call
// Contact() on a student, on a teacher, and call the "hidden"
// Person.Contact() of the teacher.
// Participation Task: create a struct that embeds two types which both
// have a field ID . Show the compile error, then fix the code so that both IDs
// are printed.
package main

import "fmt"

type Person struct {
	Name  string
	Email string
}

func (p Person) Contact() string { return p.Name + " <" + p.Email + ">" }

type Student struct {
	Person
	GPA float64
}

type Teacher struct {
	Person
	Department string
}

// Teacher's own Contact() shadows Person.Contact()
func (t Teacher) Contact() string {
	return t.Person.Contact() + ", " + t.Department + " department"
}

// Participation task
type Passport struct{ ID int }
type StudentCard struct{ ID int }

type Holder struct {
	Passport
	StudentCard
}

func main() {
	s := Student{Person{"Louay", "louay@exmaple.kz"}, 3.8}
	t := Teacher{Person{"jake", "jake@gmail.com"}, "Computer Science"}

	fmt.Println(s.Contact())        // Louay <louay@example.kz>                      (promoted from Person)
	fmt.Println(t.Contact())        // jake <jake@gmail.com>, Computer Science department (Teacher's own)
	fmt.Println(t.Person.Contact()) // jake <jake@gmail.com>                  (hidden one, via full path)

	h := Holder{Passport{111}, StudentCard{222}}

	// fmt.Println(h.ID) // compile error: ambiguous selector h.ID
	//
	// Both Passport and StudentCard are embedded at the same depth and both have ID, so Go doesn't know which
	// one h.ID means. (Declaring Holder is fine, the error only appears when you USE the ambiguous name.)
	// Fix: use the full path.
	fmt.Println(h.Passport.ID, h.StudentCard.ID) // 111 222
}
