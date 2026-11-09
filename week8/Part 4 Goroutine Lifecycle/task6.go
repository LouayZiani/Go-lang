package main

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func buildReport(ctx context.Context) (*Report, error) {
	select {
	case <-time.After(20 * time.Millisecond):
		return &Report{Payload: make([]byte, 64*1024)}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

type Report struct{ Payload []byte }

// handleLeaky demonstrates an unbuffered result channel and a timed-out caller.
func handleLeaky(timeout time.Duration) error {
	ch := make(chan *Report)
	go func() { r, _ := buildReport(context.Background()); ch <- r }()
	select {
	case <-ch:
		return nil
	case <-time.After(timeout):
		return errors.New("timeout")
	}
}

// handleBuffered lets the worker send after the caller times out, so it can exit.
func handleBuffered(timeout time.Duration) error {
	ch := make(chan *Report, 1)
	go func() { r, _ := buildReport(context.Background()); ch <- r }()
	select {
	case <-ch:
		return nil
	case <-time.After(timeout):
		return errors.New("timeout")
	}
}

// handleContext also cancels the expensive work when the caller no longer needs it.
func handleContext(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := buildReport(ctx)
	return err
}

func runBatch(handle func(time.Duration) error) {
	var wg sync.WaitGroup
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() { defer wg.Done(); _ = handle(5 * time.Millisecond) }()
	}
	wg.Wait()
}

func printMemory(label string, round int) {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%s round %d: goroutines=%d heap_alloc=%d MB\n", label, round, runtime.NumGoroutine(), m.HeapAlloc/1024/1024)
}

func runTask6Variant(label string, handle func(time.Duration) error) {
	fmt.Printf("\n=== Task 6: %s server variant ===\n", label)
	for round := 1; round <= 3; round++ {
		runBatch(handle)
		time.Sleep(30 * time.Millisecond) // let 20 ms report builders finish
		printMemory(label, round)
	}
}

func runTask6() {
	fmt.Println("Run each task6-* command separately for clean, independent measurements.")
	runTask6Variant("leaky", handleLeaky)
	runTask6Variant("buffered", handleBuffered)
	runTask6Variant("context", handleContext)
}
