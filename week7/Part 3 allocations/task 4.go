// Task 4: value vs pointer, with vs without capacity.
package main

type Point struct{ X, Y int }

func NewPointValue(x, y int) Point { return Point{x, y} }

func NewPointPointer(x, y int) *Point { return &Point{x, y} }

func CollectNoCap(n int) []int {
	var s []int
	for i := 0; i < n; i++ {
		s = append(s, i)
	}
	return s
}

func CollectWithCap(n int) []int {
	s := make([]int, 0, n)
	for i := 0; i < n; i++ {
		s = append(s, i)
	}
	return s
}
