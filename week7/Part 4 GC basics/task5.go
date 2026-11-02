// find the leak
package task5

import (
	"bytes"
	"fmt"
	"strings"
)

// readFile pretends to read a 1mb file-- the first 11 bytes are a header.
func readFile(i int) []byte {
	data := make([]byte, 1<<20)
	copy(data, fmt.Sprintf("HEADER-%04d", i))
	return data
}

// KeepHeadersSlice LEAKS: data[:11] has len 11 but still points into the 1mb backng array
// so the whole array stays reachable
func KeepHeadersSlice(n int) [][]byte {
	var headers [][]byte
	for i := 0; i < n; i++ {
		data := readFile(i)
		headers = append(headers, data[:11])
	}
	return headers
}

// KeepHeadersCopy does not leak: bytes.Clone copies just 11 bytes into a new small arr, the 1mb array becomes garbage
func KeepHeadersCopy(n int) [][]byte {
	var headers [][]byte
	for i := 0; i < n; i++ {
		data := readFile(i)
		headers = append(headers, bytes.Clone(data[:11]))
	}
	return headers
}

// bigString builds a distinct 1mb string
func bigString(i int) string {
	return fmt.Sprintf("HDR-%04d", i) + strings.Repeat("x", 1<<20)
}

// KeepSubstrings LEAKSSS: big[:5] shares memory with the 1 MB string
func KeepSubstrings(n int) []string {
	var out []string
	for i := 0; i < n; i++ {
		big := bigString(i)
		out = append(out, big[:5])
	}
	return out
}

// KeepClones doesnt leak: strings.Clone copies the 5 bytes
func KeepClones(n int) []string {
	var out []string
	for i := 0; i < n; i++ {
		big := bigString(i)
		out = append(out, strings.Clone(big[:5]))
	}
	return out
}
