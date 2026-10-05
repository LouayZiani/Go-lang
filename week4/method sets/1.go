package main

import "fmt"

type Speaker interface {
	Speak() string
}
type Cat struct{ Name string }

func (c Cat) Speak() string { return c.Name + ": meow" } // value receiver
type Dog struct{ Name string }

func (d *Dog) Speak() string { return d.Name + ": woof" } // pointer receiver
func main() {
	var s Speaker
	s = Cat{"Murka"}   // ✅ method set of Cat has Speak
	s = &Cat{"Barsik"} // ✅ method set of *Cat has Speak too
	s = &Dog{"Rex"}    // ✅ method set of *Dog has Speak
	// s = Dog{"Sharik"} // ❌ compile error (commented out so the file runs, uncomment to see it)
	fmt.Println(s.Speak()) // Rex: woof
}

// We just saw that d.Speak() works on a Dog variable (Go adds &
// for us). So why can't a Dog value be put into a Speaker ? Hint: the
// value inside an interface is a copy that is not addressable
//
// Answer:
// d.Speak() works because d is a variable, so it's addressable and the compiler quietly rewrites it to (&d).Speak().
// When you put a Dog into a Speaker, Go stores a COPY of the Dog inside the interface value. That copy has no
// variable and no stable address we could take, so there's nothing to write &x on. The compiler can't
// do the automatic & rewrite, and the method set of Dog (value) only contains value-receiver methods, so Speak isn't in it.
// It's also a safety thing: if it were allowed, Speak would modify a hidden copy and the changes
// would silently vanish. With &Dog{...} the interface holds the pointer, and the method works on the real object.
//
// Error: cannot use Dog{…} (value of struct type Dog) as Speaker value in assignment:
//        Dog does not implement Speaker (method Speak has pointer receiver)
