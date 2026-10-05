package main

import "fmt"

type Student struct {
	Name    string
	ID      int
	GPA     float64
	Courses []string
}

func main() {
	var s1 Student // all values r set to Zero

	s2 := Student{Name: "Louay", ID: 7555222, GPA: 3.5} // named fields (recommended)

	s3 := Student{"JJ", 554521, 2.1, []string{"Go", "Databases"}} // positional: ALLL fields, in order

	p := &Student{Name: "Timur"} // pointer to a new struct
	p.GPA = 2.8                  // same as (*p).GPA = 2.8 --> Go dereferences struct pointers

	// q := new(Student) // pointer to a zero Student

	// prefer NAMED FIELDS. Positional breaks as soon as someone adds/reorders a field

	// Fields that start with a CAPITAL LETTER are exported (visible from other packages)
	// Lowercase fields are private to the package

	// What will be the output here? What is the difference between %v ,%+v and %#v ?

	fmt.Printf("%v\n", s1)         // output: { 0 0 []}   (Name is "" so nothing is printed before the first space)
	fmt.Printf("%+v\n", s2)        // output: {Name:Louay ID:7555222 GPA:3.5 Courses:[]}
	fmt.Printf("%#v\n", s3)        // output: main.Student{Name:"JJ", ID:554521, GPA:2.1, Courses:[]string{"Go", "Databases"}}
	fmt.Println(s1.Courses == nil) // output: true

	// %v : values only no fields names
	// %+v : values WITH fields names
	// %#v : Go syntax representation (like source code)

	fmt.Println(p.GPA) // 2.8, just so p is used
}
