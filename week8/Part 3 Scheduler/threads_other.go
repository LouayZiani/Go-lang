//go:build !linux && !darwin && !freebsd

package main

import "fmt"

func runThreads() {
	fmt.Println("Blocking syscall thread comparison is implemented only for Linux, macOS, and FreeBSD in this lab project.")
}
