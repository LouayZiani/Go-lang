package main

import (
	"strings"
)

func WordCount(s string) map[string]int {
	slc := strings.Fields(s)

	counts := make(map[string]int)
	for i := 0; i < len(slc); i++ {
		counts[slc[i]]++
	}

	return counts
}
