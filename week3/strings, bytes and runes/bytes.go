/*
in Go, we dont use strings like most prog lans, most I/O ops are done with slice of bytes : []byte ((SANS S))
instead of: s := string(bytes.TrimSpace([]byte(data)))
its recommended to work directly with bytes: b := bytes.TrimSpace(data)

bytes package contains equivalent of all popular functions in strings: Split, Contains, Index, TrimSpace...

*/

// Task
/*What happens under the hood
In Go, string is an immutable sequence of bytes.

[]byte is a mutable slice of bytes.

When you convert []byte → string, Go copies the data into a new string object.

When you convert string → []byte, Go allocates a new slice and copies the string’s bytes into it.

That means every conversion creates extra allocations and memory copies.

For I/O (files, sockets, HTTP), Go works with []byte directly, so staying in []byte avoids this overhead.*/

package main

import (
	"bytes"
)

func processPayload(data []byte) []byte {
	// Trim spaces directly on []byte
	trimmed := bytes.TrimSpace(data)

	// Replace carriage returns directly on []byte
	clean := bytes.ReplaceAll(trimmed, []byte("\r"), []byte(""))

	return clean
}
