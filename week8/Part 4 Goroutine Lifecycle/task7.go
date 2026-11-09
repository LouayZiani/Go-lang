package main

import (
	"sync"
)

func CountSeq(n int) int { return countPrimes(0, n) }
func CountPar(n, parts int) int {
	if n <= 0 {
		return 0
	}
	if parts < 1 {
		parts = 1
	}
	if parts > n {
		parts = n
	}
	results := make([]int, parts)
	var wg sync.WaitGroup
	wg.Add(parts)
	for p := 0; p < parts; p++ {
		p := p
		go func() {
			defer wg.Done()
			from := n * p / parts
			to := n * (p + 1) / parts
			results[p] = countPrimes(from, to)
		}()
	}
	wg.Wait()
	total := 0
	for _, v := range results {
		total += v
	}
	return total
}
