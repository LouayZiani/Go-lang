package main

import (
	"fmt"
	"sync"
)

// runRaceBad intentionally reproduces the two mistakes from the handout.
// Run it as a separate process with: go run -race . racebad
func runRaceBad() {
	fmt.Println("First mistake: Add inside the goroutine (unsafe; Wait may return too early).")
	var wrongWG sync.WaitGroup
	wrongCounter := 0
	for i := 0; i < 1000; i++ {
		go func() {
			wrongWG.Add(1) // intentionally wrong
			defer wrongWG.Done()
			wrongCounter++ // intentionally raced
		}()
	}
	wrongWG.Wait()
	fmt.Println("unsafe counter (may be incomplete):", wrongCounter)

	fmt.Println("Second mistake: even with correct Add placement, unsynchronized counter++ races.")
	var wg sync.WaitGroup
	counter := 0
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() { defer wg.Done(); counter++ }()
	}
	wg.Wait()
	fmt.Println("racy counter (run with -race to detect):", counter)
}
