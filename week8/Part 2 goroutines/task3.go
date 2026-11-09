package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func runTask3() {
	fmt.Println("\n=== Task 3: goroutine cost ===")
	const n = 100_000
	var before, after runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&before)
	stop := make(chan struct{})
	var wg sync.WaitGroup
	start := time.Now()
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() { defer wg.Done(); <-stop }()
	}
	elapsed := time.Since(start)
	runtime.ReadMemStats(&after)
	fmt.Printf("goroutines alive: %d (including main)\n", runtime.NumGoroutine())
	fmt.Printf("start time for %d goroutines: %s\n", n, elapsed)
	fmt.Printf("StackSys before: %d KB; after: %d KB; delta: %d KB\n", before.StackSys/1024, after.StackSys/1024, (int64(after.StackSys)-int64(before.StackSys))/1024)
	close(stop)
	wg.Wait()
	fmt.Printf("after close: %d goroutine(s)\n", runtime.NumGoroutine())
}
