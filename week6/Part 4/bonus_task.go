package main

import (
	"fmt"
	"os"
)

func main() {
	defer fmt.Println("bye!")
	fmt.Println("hello")
	os.Exit(3)
}
