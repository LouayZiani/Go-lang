package main

import (
	"fmt"
	"os"
	"time"
)

// named return values
func double() (result int) {
	defer func() {
		result *= 2
	}()

	return 21
}

func notDouble() int {
	result := 21

	defer func() {
		result *= 2
	}()

	return result
}

// deferred cleanup and err handling
func save(name string) (err error) {
	f, err := os.Create(name)
	if err != nil {
		return err
	}

	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	_, err = f.WriteString("Aruzhan,95\n")
	return err
}

// task5.2: trace and defer function calls
func trace(name string) func() {
	fmt.Println("enter", name)

	return func() {
		fmt.Println("exit", name)
	}
}

func work() {
	defer trace("work")()
	fmt.Println("working...")
}

func timer(name string) func() {
	start := time.Now()

	return func() {
		fmt.Println(name, "took", time.Since(start))
	}
}

// timer with defer timer("sort")()
func sortExample() {
	defer timer("sort")()

	values := []int{5, 3, 1, 4, 2}
	for i := 0; i < len(values); i++ {
		for j := 0; j < len(values)-1-i; j++ {
			if values[j] > values[j+1] {
				values[j], values[j+1] = values[j+1], values[j]
			}
		}
	}
}

func main111() {
	// the four defer snippets
	defer fmt.Println("1: deferred first")
	defer fmt.Println("2: deferred second")
	defer fmt.Println("3: deferred third")

	fmt.Println("main body")

	//deferrred args vs closures
	x := 1

	defer fmt.Println("deferred x =", x)
	defer func() {
		fmt.Println("closure x =", x)
	}()

	x = 2

	fmt.Println("current x =", x)

	//named return values
	fmt.Println(double(), notDouble())

	//deferred cleanup and err handling
	fmt.Println(save("grades.csv"))

	//task5.2: trace and defer function calls
	work()

	// timer with defer timer("sort")()
	sortExample()
}
