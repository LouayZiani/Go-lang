package main

import (
	"fmt"
	"os"
	"time"
)

func producer(out chan<- int) {
	out <- 1
	close(out)
}

func consumer(in <-chan int) {
	for v := range in {
		fmt.Println("consumer got", v)
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "deadlock" {
		var n chan int
		n <- 1
	}

	ch := make(chan int)
	go func() { ch <- 42 }()
	v := <-ch
	fmt.Println(v)
	go func() { close(ch) }()
	v, ok := <-ch
	fmt.Println(v, ok)

	start := time.Now()
	since := func() time.Duration { return time.Since(start).Round(100 * time.Millisecond) }
	pc := make(chan string)
	go func() {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("sender: about to send", since())
		pc <- "ping"
		fmt.Println("sender: send finished", since())
	}()
	fmt.Println("receiver: waiting... ", since())
	msg := <-pc
	fmt.Println("receiver: got", msg, " ", since())

	c := make(chan int)
	go producer(c)
	consumer(c)
	time.Sleep(50 * time.Millisecond)
}
