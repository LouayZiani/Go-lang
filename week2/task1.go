package main

import "fmt"

func main() {
	// ARRAYS
	//1. declare arr of 5 float nums and print em
	var a = [5]float64{1.2, 1.44, 4.51, 22.4, 33.5}

	for _, v := range a {
		fmt.Println(v)
	}

	// 2. What will be the output if you initialize a 5-element array with only 3 values? Why?

	// if we only intialize the array of 5 elemnts with only 3 values, the rest of the unitialiazed to 0
	// Go automatically fills the unitialized elements with the default value and cuz we're using floats it will be 0 (def value)
	// ansd also cuz every variable needs to have defined value it cant just be uninitialized

	// 3. Declare a simple 2D array.
	var darr [3][5]int

	darr = [3][5]int{
		{1, 2, 3, 9, 7},
		{20, 54, 98, 55, 7},
		{77, 63, 21, 4, 2},
	}

	for i, v1 := range darr {
		for j, v2 := range v1 {
			fmt.Printf("darr[%d][%d] = %d\n", i, j, v2) // printf for format printting like in C, println prints values sep by space
		}
	}
}
