// gcdemo: around 22 MB of live data + arouund2.7 GB of short-lived garbage
package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

type Record struct {
	ID   int
	Data [256]byte
}

var (
	live []*Record //stays reachable: the live heap
	sum  int
	sink *Record //storing r here makes it escape --> it really goes to the heap
)

func consume(r *Record) int {
	sink = r // without this escape analysi keeps r on the stack so no garbage at all
	return r.ID & 1
}

func main() {
	start := time.Now()

	for i := 0; i < 80_000; i++ {
		live = append(live, &Record{ID: i})
	}
	for i := 0; i < 10_000_000; i++ {
		r := &Record{ID: i} //garbage after one use
		sum += consume(r)
	}
	elapsed := time.Since(start)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	gogc := os.Getenv("GOGC")
	if gogc == "" {
		gogc = "100 (default)"
	}
	fmt.Printf("GOGC=%-14s GOMEMLIMIT=%-8s\n", gogc, os.Getenv("GOMEMLIMIT"))
	fmt.Printf("  GC cycles:      %d\n", m.NumGC)
	fmt.Printf("  heap from OS:   %d MB\n", m.HeapSys>>20)
	fmt.Printf("  total GC pause: %.1f ms\n", float64(m.PauseTotalNs)/1e6)
	fmt.Printf("  allocated:      %d MB total, live now: %d MB\n", m.TotalAlloc>>20, m.HeapAlloc>>20)
	fmt.Printf("  time:           %.2f s\n", elapsed.Seconds())
	fmt.Println("  (sum =", sum, ", live records =", len(live), ")")
}
