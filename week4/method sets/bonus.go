// r := Rectangle{3, 4}
// f := r.Area // a "method value"
// r.Width = 10
// fmt.Println(f(), r.Area()) // ?

// What is printed? When exactly is the receiver copied?

package main

import "fmt"

type Rectangle struct{ Width, Height float64 }

func (r Rectangle) Area() float64 { return r.Width * r.Height }

func main() {
	r := Rectangle{3, 4}
	f := r.Area // a "method value"
	r.Width = 10
	fmt.Println(f(), r.Area()) // ?  -> 12 40
}

// What is printed? When exactly is the receiver copied?
//
// Prints: 12 40
// The receiver is copied at the moment the method value is created (f := r.Area), NOT when f() is called.
// Area has a value receiver, so f captured a copy of r with Width=3. Changing r.Width to 10 afterwards
// doesn't touch that copy, so f() = 3*4 = 12. r.Area() is a fresh call on the current r: 10*4 = 40.
// (If Area had a pointer receiver, f would capture &r and f() would see the change -> 40 40.)
