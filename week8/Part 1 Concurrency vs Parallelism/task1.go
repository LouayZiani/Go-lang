package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// countPrimes counts prime integers in [from, to).
func countPrimes(from, to int) int {
	count := 0
	for n := from; n < to; n++ {
		if n < 2 {
			continue
		}
		prime := true
		for d := 2; d*d <= n; d++ {
			if n%d == 0 {
				prime = false
				break
			}
		}
		if prime {
			count++
		}
	}
	return count
}

// cpuWork splits [0, 6_000_000) into approximately equal non-overlapping ranges.
func cpuWork(workers int) (int, time.Duration) {
	const n = 6_000_000
	if workers < 1 {
		workers = 1
	}
	start := time.Now()
	results := make([]int, workers)
	var wg sync.WaitGroup
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		w := w // explicit capture; also works with pre-Go-1.22 language semantics
		go func() {
			defer wg.Done()
			from := n * w / workers
			to := n * (w + 1) / workers
			results[w] = countPrimes(from, to)
		}()
	}
	wg.Wait()
	total := 0
	for _, v := range results {
		total += v
	}
	return total, time.Since(start)
}

// ioWork simulates independent requests that each wait for 100 ms.
func ioWork(workers int) time.Duration {
	if workers < 1 {
		workers = 1
	}
	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func() { defer wg.Done(); time.Sleep(100 * time.Millisecond) }()
	}
	wg.Wait()
	return time.Since(start)
}

func runTask1() {
	fmt.Println("\n=== Task 1: concurrency vs parallelism ===")
	fmt.Printf("runtime.NumCPU(): %d; GOMAXPROCS: %d\n", runtime.NumCPU(), runtime.GOMAXPROCS(0))
	cores := runtime.NumCPU()
	workersToTry := []int{1, cores, 2 * cores}
	seen := make(map[int]bool)
	for _, workers := range workersToTry {
		if seen[workers] {
			continue
		}
		seen[workers] = true
		count, elapsed := cpuWork(workers)
		fmt.Printf("CPU work, %d goroutine(s): %d primes in %s\n", workers, count, elapsed.Round(time.Millisecond))
	}
	for _, workers := range []int{1, 10, 1000} {
		fmt.Printf("I/O work, %d goroutine(s): %s\n", workers, ioWork(workers).Round(time.Millisecond))
	}
	fmt.Println("Repeat with one P using: GOMAXPROCS=1 go run . task1 (Unix) or set GOMAXPROCS=1 in your shell (Windows PowerShell: $env:GOMAXPROCS=1).")
}
