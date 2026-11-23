package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

func filled(s string) chan string {
	c := make(chan string, 1000)
	for range 1000 {
		c <- s
	}
	return c
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

type Job struct{ created time.Time }

func experiment(name string, size int, drop bool) {
	jobs := make(chan Job, size)
	var wg sync.WaitGroup
	handled, wait := 0, time.Duration(0)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := range jobs {
			wait = max(wait, time.Since(j.created))
			handled++
			time.Sleep(time.Millisecond)
		}
	}()
	start, dropped, queue := time.Now(), 0, 0
	for i := 0; i < 2000; i++ {
		if i > 0 && i%50 == 0 {
			time.Sleep(10 * time.Millisecond)
		}
		j := Job{time.Now()}
		if drop {
			select {
			case jobs <- j:
			default:
				dropped++
			}
		} else {
			jobs <- j
		}
		queue = max(queue, len(jobs))
	}
	busy := time.Since(start).Round(10 * time.Millisecond)
	close(jobs)
	wg.Wait()
	fmt.Printf("%-9s busy %v queue %d handled %d dropped %d wait %v\n", name, busy, queue, handled, dropped, wait.Round(10*time.Millisecond))
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

func main() {
	a, b, count := filled("a"), filled("b"), map[string]int{}
	for i := 0; i < 1000; i++ {
		select {
		case v := <-a:
			count[v]++
		case v := <-b:
			count[v]++
		}
	}
	fmt.Println("both ready, 1000 selects:", count)

	slow := make(chan string, 1)
	go func() {
		time.Sleep(200 * time.Millisecond)
		slow <- "answer"
	}()
	select {
	case r := <-slow:
		fmt.Println("got", r)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("timeout after 100ms")
	}

	queue := make(chan int, 2)
	for i := 1; i <= 4; i++ {
		select {
		case queue <- i:
			fmt.Println("queued", i)
		default:
			fmt.Println("queue full, dropped", i)
		}
	}

	for v := range merge(feed(2*time.Millisecond, 1, 2, 3), feed(time.Millisecond, 100, 101, 102, 103, 104)) {
		fmt.Print(v, " ")
	}
	fmt.Println()

	done := make(chan struct{})
	jobs := producer([]string{"go", "channels", "select"}, done)
	fmt.Println(<-jobs)
	close(done)

	experiment("blocking", 10, false)
	experiment("huge", 2000, false)
	experiment("dropping", 10, true)

	q := make(chan string, 1)
	fmt.Println(Submit(q, "job1", 50*time.Millisecond), Submit(q, "job2", 50*time.Millisecond))
}
