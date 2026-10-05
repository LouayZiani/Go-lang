package main

import "fmt"

// Task 5: Benchmark value vs pointer receivers
// In Week 3 you measured allocations. This time you'll measure copying.
//
// Report is deliberately large: 1024 * 8 bytes = 8 KB.
type Report struct {
	Title  string
	Values [1024]float64
}

//go:noinline
func (r Report) FirstValue() float64 { return r.Values[0] }

//go:noinline
func (r *Report) FirstPointer() float64 { return r.Values[0] }

// small struct for the second pair of benchmarks
type Point struct{ X, Y float64 }

//go:noinline
func (p Point) SumValue() float64 { return p.X + p.Y }

//go:noinline
func (p *Point) SumPointer() float64 { return p.X + p.Y }

func main() {
	r := Report{Title: "demo"}
	r.Values[0] = 42
	fmt.Println(r.FirstValue(), r.FirstPointer())
	p := Point{1, 2}
	fmt.Println(p.SumValue(), p.SumPointer())
	fmt.Println("run the benchmarks with: go test -bench=. -benchmem task5.go task5_test.go")
}

/*
My results (go test -bench=. -benchmem task5.go task5_test.go):

BenchmarkFirstValue      ~95 ns/op   0 B/op  0 allocs/op
BenchmarkFirstPointer    ~3 ns/op    0 B/op  0 allocs/op
BenchmarkPointValue      ~3 ns/op    0 B/op  0 allocs/op
BenchmarkPointPointer    ~3 ns/op    0 B/op  0 allocs/op

Comparison:
- Report (8 KB): value receiver is ~30x slower than pointer receiver. Every call copies 8 KB.
- Point (16 bytes): no real difference, copying 16 bytes is as cheap as passing a pointer.
  So value receivers are fine for small types, pointers pay off for big structs.

Why is B/op 0 even though the value receiver copies 8 KB per call?
B/op and allocs/op count HEAP allocations only. The copy of the receiver is made on the stack, in the
caller's frame / argument area, right before the call. Stack memory is reused for every call and is never
allocated by the GC allocator, so it isn't counted.

Where does the copy happen?
At the call site: the caller copies the whole Report into the argument slot for the receiver, then the method
runs on that copy. (A method is just a function with the receiver as first parameter.)

Which column shows the cost?
ns/op. Copying costs time, not heap memory: 95 ns vs 3 ns.

What does //go:noinline do and why do we need it?
It tells the compiler not to inline the function into the caller. These methods are tiny, so normally the
compiler would inline them. After inlining there is no call, so no receiver copy is needed. The compiler can see
that only Values[0] is read and skip copying the other 8 KB, so both versions would look the same (~1 ns) and
we would be benchmarking nothing. With noinline there's a real call, so the copy really happens.
*/
