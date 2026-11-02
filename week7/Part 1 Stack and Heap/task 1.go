package main

import "fmt"

func depth(n int) int {
	var pad [1024]byte //after increasing the size of pad, i got a stack overflow
	pad[n%len(pad)] = byte(n)
	if n == 0 {
		return 0
	}

	// return depth(n-1) + 1 + int(pad[(n+1)%len(pad)]) // reads a 0, but the compiler can't know that
	return depth(n-1) + 1 // we removed the place where the program reads pad
}

func main11() {
	// fmt.Println("depth:", depth(1_000_000)) // works!
	// fmt.Println(forever(0))// got a stack overflow
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered:", r)
		}
	}()

	forever(0) // doesnt work still terminates with a stack ovreflow
	// a regular panic: can be caught with recover()
	// a fatal runtime error: terminates the program and cannot be recovered from = stack overflow
}

// depth(1000000) doesnt crash cuz Go automatically grows goroutine stacks when they run out of space

// increased pad to 1024 --> meaning each function call requires more memory now
// Go has max goroutine stack size: 1GB on a 64bit system, so when each recursive call needs a much larger
// stack frame (1024) the toal stack req for 1 Million becomes enormous

func forever(n int) int {
	return forever(n+1) + 1
}

//participation task:
/* Why the goroutine stack needs to be copied when it grows nd why that is safe even if some pointers
point into the stack?

The goroutine stack is copied to a larger memory area when it runs out of space,
and Go safely updates pointers that refer to the old stack so they point to the correct locations in the new one.

*/
