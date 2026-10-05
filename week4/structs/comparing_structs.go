package main

import "fmt"

type Point struct{ X, Y int }

// Team is copied here from structs_are_values.go so this file runs on its own
type Team struct {
	Name    string
	Scores  [3]int
	Members []string
}

func main() {
	a, b := Point{1, 2}, Point{1, 2}
	c, d := &Point{1, 2}, &Point{1, 2}
	fmt.Println(a == b)   // true
	fmt.Println(c == d)   // false
	fmt.Println(*c == *d) // true

	t1, t2 := Team{}, Team{}
	_, _ = t1, t2
	// fmt.Println(t1 == t2) // uncomment to see the error
}

// Predict the three results and explain them. Then try t1 == t2 with the Team struct above. What does the compiler say?
//
// a == b   -> true.  Same type, all fields comparable, and X and Y are equal. Structs are compared field by field.
// c == d   -> false. c and d are POINTERS. Pointer comparison compares addresses, and
//             each &Point{1, 2} creates a different object, so the addresses differ.
// *c == *d -> true.  After dereferencing we compare two Point values again, field by field.
//
// t1 == t2 does not compile:
//   invalid operation: t1 == t2 (struct containing []string cannot be compared)
// Team has a slice field and slices are not comparable (only to nil), so the whole struct isn't comparable.
// (Same for maps and funcs.)
