package main

import (
	"fmt"
	"os"
	"sync"
)

func generator(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

func wrong() {
	results := make(chan int)
	var wg sync.WaitGroup
	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 3; i++ {
				results <- w*10 + i
			}
			close(results)
		}()
	}
	sum := 0
	for v := range results {
		sum += v
	}
	fmt.Println("sum:", sum)
	wg.Wait()
}

func collect(workers int) <-chan int {
	results := make(chan int)
	var wg sync.WaitGroup
	for w := 1; w <= workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 3; i++ {
				results <- w*10 + i
			}
		}()
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	return results
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "wrong" {
		wrong()
		return
	}
	for v := range generator(1, 2, 3) {
		fmt.Print(v, " ")
	}
	fmt.Println()
	sum := 0
	for v := range collect(3) {
		sum += v
	}
	fmt.Println("sum:", sum)
}
