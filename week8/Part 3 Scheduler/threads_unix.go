//go:build linux || darwin || freebsd

package main

import (
	"fmt"
	"runtime/pprof"
	"sync"
	"syscall"
	"time"
)

func threadCount() int { return pprof.Lookup("threadcreate").Count() }
func runThreads() {
	fmt.Println("\n=== Thread creation: timer vs blocking syscall ===")
	var wg sync.WaitGroup
	before := threadCount()
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() { defer wg.Done(); time.Sleep(200 * time.Millisecond) }()
	}
	wg.Wait()
	fmt.Printf("time.Sleep: threadcreate count %d -> %d\n", before, threadCount())
	before = threadCount()
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			ts := syscall.NsecToTimespec(int64(200 * time.Millisecond))
			_ = syscall.Nanosleep(&ts, nil)
		}()
	}
	wg.Wait()
	fmt.Printf("Nanosleep syscall: threadcreate count %d -> %d\n", before, threadCount())
}
