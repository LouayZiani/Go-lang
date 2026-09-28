package main

import (
	"fmt"
	"strings"
)

// slow string concatenation
// incorrect solution

func concat_wrong(values []string) string {
	st := ""
	for _, v := range values {
		st += v
	}

	return st
	// problem with this is that every += creates a new string --> very costly cuz we got many strings
	// --> best approach : strings.Builder
}

func concat1(values []string) string {
	var sb strings.Builder
	for _, v := range values {
		sb.WriteString(v)
	}

	return sb.String()
}

func main1() {
	vals := []string{"luk", "lu", "isl"}
	fmt.Println(concat_wrong(vals))
	fmt.Println(concat1(vals))

}
