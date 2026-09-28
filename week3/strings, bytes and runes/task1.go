package main

import (
	"fmt"
	"strings"
)

func concat(values []string) string {
	// task: Demonstrate the concatenation of strings using "string.Builder"'s "Grow" method
	// Grow(n) pre-allocates space for at least n bytes in the builder
	// --> this avoids repeated mem allocations when we alr know the final size in advance
	// output is same but the process is faster and more efficient
	var sb strings.Builder

	// calculate total length of the final string
	totalLength := 0
	for _, v := range values {
		totalLength += len(v)
	}

	// reserve memory in advance
	sb.Grow(totalLength)

	// write strings into the builder
	for _, v := range values {
		sb.WriteString(v)
	}

	return sb.String()
}

func main() {
	vals := []string{"luk", "lu", "isl"}
	fmt.Println(concat(vals))
}
