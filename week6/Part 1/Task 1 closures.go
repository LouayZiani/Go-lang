package main

import "fmt"

func makeAccumulator() func(int) int {
	running_total := 0
	return func(n int) int {
		running_total += n
		return running_total
	}
}

func memoize(f func(int) int) func(int) int {
	cache := make(map[int]int)

	return func(n int) int {
		if res, ok := cache[n]; ok {
			return res
		}

		result := f(n)
		cache[n] = result
		return result
	}
}

func main11() {
	acc := makeAccumulator()

	fmt.Println(acc(10))
	fmt.Println(acc(5))
	fmt.Println(acc(-3))

	// 2nd task
	calls := 0

	slow := func(n int) int {
		calls++
		fmt.Println("calculating:", n)
		return n * n
	}

	fast := memoize(slow)

	fmt.Println(fast(4))
	fmt.Println(fast(4))
	fmt.Println(fast(5))
	fmt.Println(fast(4))

	fmt.Println("slow function calls:", calls)

	// participation task
	// When counter() returns, its local var count remains alive because the returned closure still references it.
	// The compiler/runtime ensures that the captured variable remains valid for as long as it is needed.
	// c1 and c2 do not share a counter because each call to counter() creates a distinct captured var
}
