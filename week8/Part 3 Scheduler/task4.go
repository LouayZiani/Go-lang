package main

import (
	"fmt"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
)

func busy(d time.Duration) {
	deadline := time.Now().Add(d)
	var x uint64
	for time.Now().Before(deadline) {
		x++
	}
	_ = x
}

func runTask4() {
	fmt.Println("\n=== Task 4: scheduler ===")
	fmt.Printf("NumCPU=%d GOMAXPROCS=%d\n", runtime.NumCPU(), runtime.GOMAXPROCS(0))
	var wg sync.WaitGroup
	wg.Add(8)
	for i := 0; i < 8; i++ {
		go func() { defer wg.Done(); busy(2 * time.Second) }()
	}
	wg.Wait()
	fmt.Println("Eight CPU-busy goroutines finished. Re-run with GODEBUG=schedtrace=500 go run . task4 to inspect scheduler snapshots.")

	fmt.Println("Run `go run . preempt` for the preemption demonstration.")
	fmt.Println("Run `go run . threads` to compare timer waits with blocking syscalls on supported Unix systems.")
	_ = pprof.Lookup("threadcreate")
}

func runPreemption() {
	runtime.GOMAXPROCS(1)
	go func() {
		for { /* intentionally CPU-bound loop */
		}
	}()
	time.Sleep(10 * time.Millisecond)
	fmt.Println("main is still alive (asynchronous preemption allowed)")
}
