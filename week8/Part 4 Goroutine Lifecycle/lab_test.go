package main

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestCountPar(t *testing.T) {
	want := CountSeq(200_000)
	for _, parts := range []int{1, 2, 8, 1000, 200_000, 250_000} {
		if got := CountPar(200_000, parts); got != want {
			t.Fatalf("parts=%d: CountPar=%d, CountSeq=%d", parts, got, want)
		}
	}
}

func TestFetchFixedDoesNotLeak(t *testing.T) {
	before := runtime.NumGoroutine()
	for i := 0; i < 20; i++ {
		_, _ = fetchFixed(time.Millisecond)
	}
	if !waitForGoroutines(before, time.Second) {
		t.Fatalf("goroutines did not return near baseline: before=%d after=%d", before, runtime.NumGoroutine())
	}
}

var sink int

const benchmarkN = 200_000

func BenchmarkCountSeq(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sink = CountSeq(benchmarkN)
	}
}
func BenchmarkCountPar(b *testing.B) {
	for _, parts := range []int{2, 8, 1000, 200_000} {
		b.Run(fmt.Sprintf("parts=%d", parts), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				sink = CountPar(benchmarkN, parts)
			}
		})
	}
}

func BenchmarkGoroutineStartStop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		done := make(chan struct{})
		go func() { close(done) }()
		<-done
	}
}
func BenchmarkFuncCall(b *testing.B) {
	f := func() {}
	for i := 0; i < b.N; i++ {
		f()
	}
}
