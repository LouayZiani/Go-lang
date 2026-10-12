package main

import (
	"fmt"
)

func main33() {
	var a, b any = 1, 1
	fmt.Println(a == b) // true

	var c, d any = 1, int64(1)
	fmt.Println(c == d) // false (different dynamic types)

	var e, f any = []int{1}, []int{1}
	fmt.Println(e == f) // PANIC!

}
