package main

import (
	"fmt"
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

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

func filterEven(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			if n%2 == 0 {
				out <- n
			}
		}
	}()
	return out
}

func sum(in <-chan int) int {
	s := 0
	for n := range in {
		s += n
	}
	return s
}

func work(w int, out chan<- int) {
	for i := 0; i < 3; i++ {
		out <- w*10 + i
	}
}

func wrong(workers int) {
	results := make(chan int)
	var wg sync.WaitGroup
	for w := 1; w <= workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			work(w, results)
			close(results)
		}()
	}
	fmt.Println("sum:", sum(results))
	wg.Wait()
}

func collect(workers int) <-chan int {
	results := make(chan int)
	var wg sync.WaitGroup
	for w := 1; w <= workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			work(w, results)
		}()
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	return results
}
