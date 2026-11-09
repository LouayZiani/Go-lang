package main

import (
	"runtime"
	"sync"
	"time"
)

// Used in tests to ensure timed-out fetch workers eventually finish.
func waitForGoroutines(target int, max time.Duration) bool {
	deadline := time.Now().Add(max)
	for time.Now().Before(deadline) {
		if runtime.NumGoroutine() <= target {
			return true
		}
		time.Sleep(time.Millisecond)
	}
	return runtime.NumGoroutine() <= target
}

var _ sync.WaitGroup
