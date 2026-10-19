package main

import (
	"fmt"
	"strings"
)

func apply(nums []int, f func(int) int) []int {
	out := make([]int, len(nums))
	for i, n := range nums {
		out[i] = f(n)
	}

	return out
}

func double(n int) int {
	return n * 2
}

func divmod(a, b int) (q, r int) { // named returns: useful for doc and for "defer"== we'll see later
	q = a / b
	r = a % b
	return // naked return === returns q and r ( avoid for long fcts == user gotta scroll up to see whats returned)
}

// variadic function is a fct that accepts a variable num of arguments
func sum(nums ...int) int { // variadic: nums is []int
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// closures
func counter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

// closure that builds functions
func multiplier(k int) func(int) int {
	return func(n int) int { return n * k }
}

func main1() {
	// in Go, fcts are jst values. Fcts can be stored in a var, put into a map or slice, or passed to another fct and
	// returned from it, just like an int or a string
	fmt.Println(apply([]int{1, 2, 3}, double))                           // [2 4 6]
	fmt.Println(apply([]int{1, 2, 3}, func(n int) int { return n * n })) // [1 4 9] == using anonym fct, function literal

	var f func(string) string = strings.ToUpper
	fmt.Println(f("go")) // GO

	// zero value of a function type is nil. Calling a nil function panics

	// function type with a name
	type Op func(a, b float64) float64
	ops := map[string]Op{
		"+": func(a, b float64) float64 { return a + b },
		"*": func(a, b float64) float64 { return a * b },
	}
	fmt.Println(ops["+"](2, 3), ops["*"](2, 3)) // 5 6

	q, r := divmod(17, 5)
	fmt.Println(q, r)
	fmt.Println(sum(), sum(1, 2), sum([]int{4, 5, 6}...)) // ... here is variadic argument expansion operator

	// CLOSURE: is fct that uses variables declared outside of it
	c1 := counter()
	fmt.Println(c1(), c1(), c1()) // 1 2 3
	c2 := counter()
	fmt.Println(c2(), c1()) // 1 4

	// closures capture variables not values, it refers to var x itself so it sees its current value at moment it runs
	x := 10
	show := func() { fmt.Println("x =", x) }
	x = 20
	show()

	// closure that builds functions
	triple := multiplier(3)
	fmt.Println(triple(5)) // 15

	// closures in loops: Go version matters here!!!!!
	var funcs []func()
	for i := 0; i < 3; i++ {
		funcs = append(funcs, func() { fmt.Println(i, "") })
	}
	for _, f := range funcs {
		f() // this would print 3 3 3 for older versions of go (one var i for the whole loop, all closures share it),
		// and would print 0 1 2 for newer version 1.22 or newer (every iteration gets its own var i)
	}

}
