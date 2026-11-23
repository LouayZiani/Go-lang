package main

import (
	"fmt"
	"time"
)

func run(name string, ch chan int) {
	start := time.Now()
	since := func() time.Duration { return time.Since(start).Round(100 * time.Millisecond) }
	done := make(chan struct{})
	go func() {
		for i := 1; i <= 3; i++ {
			ch <- i
			fmt.Printf("%s: sent %d at %v\n", name, i, since())
		}
		close(done)
	}()
	time.Sleep(300 * time.Millisecond)
	for i := 0; i < 3; i++ {
		<-ch
		time.Sleep(100 * time.Millisecond)
	}
	<-done
}

func main() {
	run("unbuffered", make(chan int))
	run("buffered(3)", make(chan int, 3))

	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	fmt.Println("len:", len(ch), "cap:", cap(ch))
	close(ch)
	v, ok := <-ch
	fmt.Println(v, ok)
	v, ok = <-ch
	fmt.Println(v, ok)
	v, ok = <-ch
	fmt.Println(v, ok)
}
