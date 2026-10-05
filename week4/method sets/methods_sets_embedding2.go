package main

import "fmt"

// fmt.Println uses a type's String() string method if the type has
// one (the fmt.Stringer interface). What will this program print, and
// why are the two lines different?

type Book struct {
	Title  string
	Author string
}

func (b *Book) String() string { return b.Title + " by " + b.Author }
func main() {
	b := Book{"The Go Programming Language", "Donovan & Kernighan"}
	fmt.Println(b)  // {The Go Programming Language Donovan & Kernighan}
	fmt.Println(&b) // The Go Programming Language by Donovan & Kernighan
}

// Output:
// {The Go Programming Language Donovan & Kernighan}
// The Go Programming Language by Donovan & Kernighan
//
// String() has a pointer receiver, so it's only in the method set of *Book.
// Println(b) passes a Book VALUE as interface{}. Book's method set has no String(), so it isn't a fmt.Stringer
// and fmt falls back to the default struct printing {Title Author}. Go does NOT take the address here
// (the copy inside the interface isn't addressable, same reason as Dog{} not being a Speaker).
// Println(&b) passes *Book, whose method set contains String(), so fmt calls it.
// Fix: make String a value receiver, then both lines print the nice form.
