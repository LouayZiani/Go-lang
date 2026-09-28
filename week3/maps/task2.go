package main

import (
	"fmt"
)

func main12() {
	// inventory management program
	inventory := map[string]int{
		"apples":  10,
		"bananas": 5,
	}

	// Read
	fmt.Println("Apples stock:", inventory["apples"])

	// Update
	inventory["bananas"] = 12

	// Insert
	inventory["oranges"] = 8

	// Delete
	delete(inventory, "apples")

	fmt.Println("Updated inventory:", inventory)

}
