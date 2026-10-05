package main

import "fmt"

type Team struct {
	Name    string
	Scores  [3]int
	Members []string
}

func main() {
	t1 := Team{Name: "Gophers", Scores: [3]int{1, 2, 3}, Members: []string{"Ali", "Bota"}}
	t2 := t1 // copy
	t2.Name = "Rustaceans"
	t2.Scores[0] = 100
	t2.Members[0] = "Nurlan"
	fmt.Println(t1.Name, t1.Scores, t1.Members)
	fmt.Println(t2.Name, t2.Scores, t2.Members)
}

// What will be the output? Why did one field change in both structs, while the others didn't?
//
// Output:
// Gophers [1 2 3] [Nurlan Bota]
// Rustaceans [100 2 3] [Nurlan Bota]
//
// t2 := t1 copies every field.
// - Name is a string, so t2 gets its own copy -> only t2 changed.
// - Scores is an ARRAY ([3]int), arrays are values too, so the whole array is copied -> only t2 changed.
// - Members is a SLICE. Copying it copies only the slice header (pointer, len, cap),
//   so t1.Members and t2.Members point to the same underlying array.
//   t2.Members[0] = "Nurlan" writes into that shared array, so t1 sees it too.
// That's the one field that changed in both.
