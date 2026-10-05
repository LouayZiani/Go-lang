package main

// Which of the lines A–D compile? Explain using the table.
type Speaker interface{ Speak() string }
type Base struct{}

func (Base) Hello() string  { return "hello" }
func (*Base) Speak() string { return "speak" }

type ByValue struct{ Base }    // embeds Base
type ByPointer struct{ *Base } // embeds *Base
var _ Speaker = &ByValue{}     // A ✅
// var _ Speaker = ByValue{}              // B ❌
var _ Speaker = ByPointer{Base: &Base{}}  // C ✅
var _ Speaker = &ByPointer{Base: &Base{}} // D ✅

func main() {}

// Speak has a pointer receiver on Base, so it belongs to the method set of *Base only.
//
// A ✅ &ByValue is *S where S embeds T (Base). Table: method set of *S = methods of T and *T -> Speak is in.
// B ❌ ByValue is S where S embeds T. Method set of S = methods of T only (Hello, no Speak).
//      Error: ByValue does not implement Speaker (method Speak has pointer receiver)
//      (the embedded Base lives inside a non-addressable copy in the interface, so no &Base is possible.)
// C ✅ ByPointer embeds *T. Table: method set of S = methods of T and *T. The struct already holds the
//      pointer, so no address needs to be taken -> Speak is in.
// D ✅ *ByPointer: for embedded *T both S and *S have T and *T methods.
