// Predict then check
package main

import "fmt"

type Point struct{ X, Y int }

var global *Point

func f1() int         { x := 42; return x }                           // stays
func f2() *int        { x := 42; return &x }                          // escapes
func f3(p *Point) int { return p.X + p.Y }                            // stays: p is only read inside the fct.the pointer isnt stored or returned
func f4()             { p := Point{1, 2}; fmt.Println(p) }            // escapes: print passes p thru an interface->cause escape
func f5() []int       { s := make([]int, 10); return s }              // escapes
func f6() int         { s := make([]int, 10); return len(s) }         // stays: only the length is returned. the slice isnt
func f7(n int) int    { s := make([]int, n); return len(s) }          //escapes
func f8() func() int  { x := 0; return func() int { x++; return x } } // escapes: the returned closure needs x after f8 finishes
func f9(p *Point)     { global = p }                                  //escapes
func f10() byte       { var big [200_000]byte; return big[0] }        //escapes

func main() {
	p := Point{3, 4}
	f9(&p) // p here escapes to the heap, because f9*(&p) passes the address of p to f9, which stores that pointer
	// in the global variable "global"
}

/*go build -gcflags='-m -l' output:
# task2
./task 2.go:11:25: moved to heap: x
./task 2.go:12:9: p does not escape
./task 2.go:13:54: ... argument does not escape
./task 2.go:13:55: p escapes to heap
./task 2.go:14:34: make([]int, 10) escapes to heap
./task 2.go:15:34: make([]int, 10) does not escape
./task 2.go:16:34: make([]int, n) does not escape
./task 2.go:17:25: moved to heap: x
./task 2.go:17:40: func literal escapes to heap
./task 2.go:18:9: leaking param: p
./task 2.go:19:29: moved to heap: big
./task 2.go:22:2: moved to heap: p
*/
