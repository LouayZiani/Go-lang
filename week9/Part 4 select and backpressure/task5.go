package main

import (
	"errors"
	"fmt"
	"runtime"
	"time"
)

func filled(s string) chan string {
	c := make(chan string, 1000)
	for range 1000 {
		c <- s
	}
	return c
}

func count(preferA bool) map[string]int {
	a, b, m := filled("a"), filled("b"), map[string]int{}
	for range 1000 {
		if preferA {
			select {
			case v := <-a:
				m[v]++
				continue
			default:
			}
		}
		select {
		case v := <-a:
			m[v]++
		case v := <-b:
			m[v]++
		}
	}
	return m
}

func fetchWithTimeout(timeout time.Duration) (string, error) {
	ch := make(chan string, 1)
	go func() {
		time.Sleep(200 * time.Millisecond)
		ch <- "answer"
	}()
	select {
	case r := <-ch:
		return r, nil
	case <-time.After(timeout):
		return "", errors.New("timeout")
	}
}

func merge(a, b <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for a != nil || b != nil {
			select {
			case v, ok := <-a:
				if !ok {
					a = nil
					continue
				}
				out <- v
			case v, ok := <-b:
				if !ok {
					b = nil
					continue
				}
				out <- v
			}
		}
	}()
	return out
}

func mergeWithoutNil(a, b <-chan int) int {
	n, doneA, doneB := 0, false, false
	for !doneA || !doneB {
		select {
		case _, ok := <-a:
			if !ok {
				if doneA {
					n++
				}
				doneA = true
			}
		case _, ok := <-b:
			if !ok {
				if doneB {
					n++
				}
				doneB = true
			}
		}
	}
	return n
}

func feed(delay time.Duration, vals ...int) <-chan int {
	c := make(chan int)
	go func() {
		defer close(c)
		for _, v := range vals {
			time.Sleep(delay)
			c <- v
		}
	}()
	return c
}

func pingPong() (int, int) {
	table, done := make(chan int), make(chan struct{})
	player := func() {
		for {
			select {
			case b := <-table:
				select {
				case table <- b + 1:
				case <-done:
					return
				}
			case <-done:
				return
			}
		}
	}
	go player()
	go player()
	table <- 0
	time.Sleep(100 * time.Millisecond)
	ball := <-table
	close(done)
	time.Sleep(50 * time.Millisecond)
	return ball, runtime.NumGoroutine()
}

func task5() {
	fmt.Println("== 5.1 random choice x3:", count(false), count(false), count(false))
	fmt.Println("== 5.1 prefer a:", count(true))
	for _, t := range []time.Duration{100, 300} {
		r, err := fetchWithTimeout(t * time.Millisecond)
		fmt.Printf("== 5.2 timeout %dms: %q %v\n", t, r, err)
	}
	time.Sleep(300 * time.Millisecond)
	fmt.Println("goroutines after fetch:", runtime.NumGoroutine())
	var got []int
	for v := range merge(feed(time.Millisecond, 1, 2, 3), feed(2*time.Millisecond, 100, 101, 102)) {
		got = append(got, v)
	}
	fmt.Println("== 5.3 merge:", got)
	fmt.Println("== 5.3 without a = nil, useless zero values:", mergeWithoutNil(feed(0, 1, 2, 3), feed(time.Millisecond, 100, 101, 102)))
	hits, g := pingPong()
	fmt.Println("== participation: hits", hits, "goroutines after stop", g)
}
