package main

import (
	"errors"
	"sync"
	"testing"
)

var ErrNegative = errors.New("negative value")

//go:noinline
func checkErr(n int) (int, error) {
	if n < 0 {
		return 0, ErrNegative
	}
	return n, nil
}

//go:noinline
func checkPanic(n int) int {
	if n < 0 {
		panic(ErrNegative)
	}
	return n
}

func safeCheck(n int) (v int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	return checkPanic(n), nil
}

func BenchmarkErrorReturn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = checkErr(-1)
	}
}

func BenchmarkPanicRecover(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = safeCheck(-1)
	}
}

type Counter struct {
	mu sync.Mutex
	n  int
}

func (c *Counter) IncNoDefer() {
	c.mu.Lock()
	c.n++
	c.mu.Unlock()
}

func (c *Counter) IncDefer() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.n++
}

func BenchmarkIncNoDefer(b *testing.B) {
	var c Counter
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.IncNoDefer()
	}
}

func BenchmarkIncDefer(b *testing.B) {
	var c Counter
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.IncDefer()
	}
}
