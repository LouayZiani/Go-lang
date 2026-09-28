package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func main11() {
	// strings and intro to unicode
	s := "世界"
	fmt.Println(len(s)) // 6

	fmt.Println(utf8.RuneCountInString(s)) // 2

	// runes

	// incorrect string iteration
	s1 := "hêllo"
	for i := range s1 {
		fmt.Printf("position: %d: %c\n", i, s1[i]) // s[i] returns a byte not a rune. and char "ê" takes up 2 bytes. This loop
		// iterates over BYTE indices, not rune indices.
	}
	// correct approach
	for i, r := range s1 {
		fmt.Printf("position: %d: %c\n", i, r)
	}

	// Trim functions
	s2 := "123oxo"
	s3 := strings.TrimRight(s2, "xo")
	for _, r := range s3 {
		fmt.Printf("%c", r)
	}

	// TrimRight will remove ALL chars present in cutset (xo). It keeps on going after removing cutsett at first cuz the end
	// of the string still contains a char of cutset

	// TrimSuffix will remove the suffix wanted exactly ONE time

	s4 := strings.TrimSuffix(s2, "xo")
	fmt.Printf("\n")
	for _, r := range s4 {
		fmt.Printf("%c", r)
	}
	fmt.Printf("\n")

	// TrimPrefix similar to TrimSuffix, instead executed from the beginning
	s5 := strings.TrimPrefix(s2, "123")
	for _, r := range s5 {
		fmt.Printf("%c", r)
	}
	fmt.Printf("\n")

}
