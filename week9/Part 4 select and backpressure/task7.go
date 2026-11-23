package main

import (
	"fmt"
	"runtime"
	"testing"
)

const N = 10000

func passN(size int) int {
	ch := make(chan int, size)
	go func() {
		for i := 0; i < N; i++ {
			ch <- i
		}
		close(ch)
	}()
	sum := 0
	for v := range ch {
		sum += v
	}
	return sum
}

func passBatches() int {
	ch := make(chan []int, 100)
	go func() {
		for i := 0; i < N; i += 100 {
			b := make([]int, 100)
			for j := range b {
				b[j] = i + j
			}
			ch <- b
		}
		close(ch)
	}()
	sum := 0
	for b := range ch {
		for _, v := range b {
			sum += v
		}
	}
	return sum
}

var cases = []struct {
	name string
	f    func() int
}{
	{"buf=0", func() int { return passN(0) }},
	{"buf=1", func() int { return passN(1) }},
	{"buf=100", func() int { return passN(100) }},
	{"buf=10000", func() int { return passN(10000) }},
	{"batches", passBatches},
}

func task7() {
	fmt.Println("== 7 cost per value")
	for _, cpu := range []int{1, 2} {
		runtime.GOMAXPROCS(cpu)
		for _, c := range cases {
			r := testing.Benchmark(func(b *testing.B) {
				b.ReportAllocs()
				for range b.N {
					c.f()
				}
			})
			fmt.Printf("cpu=%d %-10s %6.1f ns/value %7d B/op %4d allocs/op\n", cpu, c.name, float64(r.NsPerOp())/N, r.AllocedBytesPerOp(), r.AllocsPerOp())
		}
	}
}
