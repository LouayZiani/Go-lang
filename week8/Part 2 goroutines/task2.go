package main

import (
	"fmt"
	"sync"
	"time"
)

func say(msg string) {
	for i := 0; i < 3; i++ {
		fmt.Println(msg, i)
		time.Sleep(100 * time.Millisecond)
	}
}

func runTask2() {
	fmt.Println("\n=== Task 2: first goroutines ===")
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); say("goroutine") }()
	go func() { defer wg.Done(); say("main") }()
	wg.Wait()

	var atomicWG sync.WaitGroup
	var counter atomicInt64
	atomicWG.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() { defer atomicWG.Done(); counter.Add(1) }()
	}
	atomicWG.Wait()
	fmt.Println("atomic counter:", counter.Load())

	values := make([]int, 1000)
	var sliceWG sync.WaitGroup
	sliceWG.Add(len(values))
	for i := range values {
		i := i
		go func() { defer sliceWG.Done(); values[i] = i + 1 }()
	}
	sliceWG.Wait()
	sum := 0
	for _, v := range values {
		sum += v
	}
	fmt.Println("sum of independently written slice:", sum)

	var panicWG sync.WaitGroup
	panicWG.Add(1)
	go func() {
		defer panicWG.Done()
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered inside goroutine:", r)
			}
		}()
		panic("demonstration panic")
	}()
	panicWG.Wait()
}
