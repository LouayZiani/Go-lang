package main

import "fmt"

type Animal struct{ Name string }

func (a Animal) Sound() string    { return "..." }
func (a Animal) Describe() string { return a.Name + " says " + a.Sound() }

type Dog struct{ Animal }

func (d Dog) Sound() string { return "Woof!" }
func main() {
	d := Dog{Animal{Name: "Rex"}}
	fmt.Println(d.Sound())
	fmt.Println(d.Animal.Sound())
	fmt.Println(d.Describe())
}

// Question: What will be printed? Why doesn't the last line print Rex says Woof! ?
//
// Output:
// Woof!
// ...
// Rex says ...
//
// d.Sound() -> Dog's own Sound shadows the promoted one, so "Woof!".
// d.Animal.Sound() -> full path reaches the hidden inner method, "...".
// d.Describe() is promoted from Animal, which means it's called with receiver d.Animal (an Animal, not a Dog).
// Inside Describe, a.Sound() is resolved statically on type Animal, so it calls Animal.Sound -> "...".
// Animal knows nothing about Dog. Go has no virtual methods / overriding: embedding is not inheritance,
// it's just shortcut forwarding to an inner field.
