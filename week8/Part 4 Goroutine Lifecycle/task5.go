package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

func slowQuery() int { time.Sleep(50 * time.Millisecond); return 42 }
func fetchLeaky(timeout time.Duration) (int, error) {
	ch := make(chan int)
	go func() { ch <- slowQuery() }()
	select {
	case v := <-ch:
		return v, nil
	case <-time.After(timeout):
		return 0, errors.New("timeout")
	}
}
func fetchFixed(timeout time.Duration) (int, error) {
	ch := make(chan int, 1)
	go func() { ch <- slowQuery() }()
	select {
	case v := <-ch:
		return v, nil
	case <-time.After(timeout):
		return 0, errors.New("timeout")
	}
}
func startTickerLeaky() {
	go func() {
		for {
			time.Sleep(10 * time.Millisecond)
		}
	}()
}
func startTicker(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(10 * time.Millisecond):
			}
		}
	}()
}

func runTask5() {
	fmt.Println("\n=== Task 5: find and fix leaks ===")
	base := runtime.NumGoroutine()
	for i := 0; i < 10; i++ {
		_, _ = fetchLeaky(1 * time.Millisecond)
	}
	time.Sleep(5 * time.Millisecond)
	fmt.Printf("baseline=%d after 10 leaky fetches=%d\n", base, runtime.NumGoroutine())
	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
	for i := 0; i < 10; i++ {
		_, _ = fetchFixed(1 * time.Millisecond)
	}
	time.Sleep(70 * time.Millisecond)
	fmt.Printf("after fixed fetches finish=%d\n", runtime.NumGoroutine())

	for i := 0; i < 3; i++ {
		startTickerLeaky()
	}
	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < 3; i++ {
		startTicker(ctx)
	}
	fmt.Printf("with leaky + stoppable tickers=%d\n", runtime.NumGoroutine())
	cancel()
	time.Sleep(20 * time.Millisecond)
	fmt.Printf("after cancel (leaky tickers remain)=%d\n", runtime.NumGoroutine())
}
