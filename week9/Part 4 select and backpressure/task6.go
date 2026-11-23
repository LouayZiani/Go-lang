package main

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Result struct{ Word, Upper string }

var words = strings.Fields("alpha beta gamma delta epsilon zeta eta theta iota kappa lambda mu nu xi omicron pi rho sigma tau upsilon")

func producer(words []string, done <-chan struct{}) <-chan string {
	jobs := make(chan string)
	go func() {
		defer close(jobs)
		for _, w := range words {
			select {
			case jobs <- w:
			case <-done:
				return
			}
		}
	}()
	return jobs
}

func startWorkers(n int, jobs <-chan string, done <-chan struct{}) <-chan Result {
	results := make(chan Result)
	var wg sync.WaitGroup
	for range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for w := range jobs {
				select {
				case results <- Result{w, strings.ToUpper(w)}:
				case <-done:
					return
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	return results
}

func runPool(take int, useDone bool) int {
	var done chan struct{}
	if useDone {
		done = make(chan struct{})
	}
	n := 0
	for range startWorkers(3, producer(words, done), done) {
		n++
		if n == take {
			break
		}
	}
	if useDone {
		close(done)
	}
	time.Sleep(50 * time.Millisecond)
	return n
}

var ErrQueueFull = errors.New("queue is full")

func Submit(queue chan<- string, job string, timeout time.Duration) error {
	t := time.NewTimer(timeout)
	defer t.Stop()
	select {
	case queue <- job:
		return nil
	case <-t.C:
		return ErrQueueFull
	}
}

func task6() {
	fmt.Println("== 6 full run, goroutines:", runPool(0, true), runtime.NumGoroutine())
	fmt.Println("== 6 take 3 + close(done), goroutines:", runPool(3, true), runtime.NumGoroutine())
	q := make(chan string, 1)
	fmt.Println("== 6 Submit:", Submit(q, "job1", 50*time.Millisecond), Submit(q, "job2", 50*time.Millisecond))
	fmt.Println("== 6 take 3 without done, goroutines:", runPool(3, false), runtime.NumGoroutine())
}
