// Task 3: Compile or not?
// Without running the code, mark every numbered line as ✅ compiles or ❌
// does not compile, and explain why. Then check yourself with go build .
package main

type Counter struct{ n int }

func (c Counter) Get() int { return c.n }
func (c *Counter) Inc()    { c.n++ }

func NewCounter() Counter { return Counter{} }

func main() {
	var c Counter
	c.Inc() // 1 ✅ c is a variable -> addressable, Go rewrites it as (&c).Inc()

	(&c).Inc() // 2 ✅ we take the address ourselves

	Counter{}.Get() // 3 ✅ Get has a value receiver, no address needed

	// Counter{}.Inc() // 4 ❌ cannot call pointer method Inc on Counter
	//   a struct literal is not addressable, so Go can't do (&Counter{}).Inc()

	// NewCounter().Inc() // 5 ❌ cannot call pointer method Inc on Counter
	//   a function return value is a temporary, not addressable

	m := map[string]Counter{"a": {}}
	m["a"].Get() // 6 ✅ value receiver, works on a copy of the element

	// m["a"].Inc() // 7 ❌ cannot call pointer method Inc on Counter
	//   map elements are not addressable (the map can move entries around when it grows)

	s := []Counter{{}}
	s[0].Inc() // 8 ✅ slice elements are addressable -> (&s[0]).Inc()

	pm := map[string]*Counter{"a": {}}
	pm["a"].Inc() // 9 ✅ the map stores a *Counter, so the pointer is already there. We don't need the address of the map element, we use the pointer it holds
}

// Result: 1 ✅  2 ✅  3 ✅  4 ❌  5 ❌  6 ✅  7 ❌  8 ✅  9 ✅
// Checked with go build, the three errors were on 4, 5 and 7, each "cannot call pointer method Inc on Counter".
// Rule: a pointer method needs an addressable value (variable, slice element, field of addressable struct, *p)
// or an actual pointer.
