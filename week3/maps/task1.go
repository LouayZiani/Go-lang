package main

import (
	"fmt"
)

func main1() {
	scores := map[string]int{
		"Alice": 90,
		"Bob":   0, // Bob exists but scored 0
	}

	// Check Bob
	if val, ok := scores["Bob"]; ok {
		fmt.Println("Bob is in the map with score", val)
	} else {
		fmt.Println("Bob is missing")
	}

	// Check Charlie
	if val, ok := scores["Charlie"]; ok {
		fmt.Println("Charlie is in the map with score", val)
	} else {
		fmt.Println("Charlie is missing")
	}

}
