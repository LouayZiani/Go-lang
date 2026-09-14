package main

import "fmt"

// SLICES
// 1. According to the tour of Go and other materials, declare a slice in 4 different ways.
// premier c slice literal
s1 := []int{11, 221, 54}

// from an arr
arr1 := [5]int{1, 2, 3, 4, 3}
s2 := arr1[1:4]

// make
s3 = make([]int, 4, 5)

// empty slice + appendd
var s4 []int
s5 = append(s5, 44, 55, 23)

// 2. modif to slice and impact on the underlaying arr

arr := [5]int{13, 12, 45, 90, 11}
s := arr[1:4]

fmt.Println("og array:", arr)
fmt.Println("slice:", s)

s[1] = 99
fmt.Println("after modifying slice:")
fmt.Println("arr:", arr)
fmt.Println("slice:", s)

arr[2] = 77
fmt.Println("after modifying arr:")
fmt.Println("arr:", arr)
fmt.Println("slice:", s)

// 3. theoretical questions

// array is a fixed size collection of elements of same type, size is part of the type, stores contiguously in memory
// slice is a flexible dynamic view into an array, internally stores pointer to backing arr, capacity and length, can grow
// and shrink with append unlike array

// participation task

// fct that modifies slice elements
func modify(s []int) {
	for i := range s {
		s[i] = s[i] * 2
	}
}

func main() {
	arr_nums := []int{3, 4, 1, 2}
	fmt.Println("before:", arr_nums)

	modify(arr_nums) 

	fmt.Println("after:", arr_nums)
}

// the og slice in main() is modified , cuz slices r refs to an underlying backing arr
