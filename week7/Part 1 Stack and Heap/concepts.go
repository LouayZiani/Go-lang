// every go program has two kinds of memory: stack and heap
// --> stack is fast, private to one goroutine and cleaned automatically when a fct returns
// --> heap is shared , lives as long as something points to it, and is cleaned by the garbage collector (GC)

// stack: good for loca, short-lived values
// heap: good for values that outlive their function or are shared

package main

// import "fmt"

func main1() {
	// when a function is called, Go creates a stack frame (small area of mem) for that fct's local vars and params

	// in Go we do not have dangling pointer problem, because we do not choose stack or heap, the compiler decides where
	// each value lives

	// goroutine is a lightweight thread managed by the Go runtime

}
