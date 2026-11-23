package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

var order = []string{"nil send", "nil recv", "nil close", "nil len/cap", "open send", "open recv", "open close", "open len/cap", "closed send", "closed recv", "closed close", "closed len/cap"}

var cells = map[string]func(){
	"nil send":     func() { var c chan int; c <- 1 },
	"nil recv":     func() { var c chan int; <-c },
	"nil close":    func() { var c chan int; close(c) },
	"nil len/cap":  func() { var c chan int; fmt.Println(len(c), cap(c)) },
	"open send":    func() { c := make(chan int); c <- 1 },
	"open recv":    func() { c := make(chan int); <-c },
	"open close":   func() { c := make(chan int); close(c); fmt.Println("ok") },
	"open len/cap": func() { c := make(chan int, 3); c <- 1; c <- 2; fmt.Println(len(c), cap(c)) },
	"closed send":  func() { c := make(chan int, 1); close(c); c <- 1 },
	"closed recv": func() {
		c := make(chan int, 1)
		c <- 7
		close(c)
		v, ok := <-c
		fmt.Println(v, ok)
		v, ok = <-c
		fmt.Println(v, ok)
	},
	"closed close":   func() { c := make(chan int); close(c); close(c) },
	"closed len/cap": func() { c := make(chan int, 3); c <- 1; close(c); fmt.Println(len(c), cap(c)) },
}

func ping(receiverSleeps bool) {
	start := time.Now()
	since := func() time.Duration { return time.Since(start).Round(100 * time.Millisecond) }
	ch := make(chan string)
	go func() {
		if !receiverSleeps {
			time.Sleep(500 * time.Millisecond)
		}
		fmt.Println("sender: about to send", since())
		ch <- "ping"
		fmt.Println("sender: send finished", since())
	}()
	fmt.Println("receiver: waiting... ", since())
	if receiverSleeps {
		time.Sleep(500 * time.Millisecond)
	}
	msg := <-ch
	fmt.Println("receiver: got", msg, " ", since())
	time.Sleep(50 * time.Millisecond)
}

func say(s string, done chan<- struct{}) {
	for i := 0; i < 3; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
	done <- struct{}{}
}

func main() {
	if len(os.Args) > 1 {
		cells[os.Args[1]]()
		return
	}
	fmt.Println("== 1.1 sleep in sender")
	ping(false)
	fmt.Println("== 1.1 sleep in receiver")
	ping(true)
	fmt.Println("== 1.2 table")
	exe, _ := os.Executable()
	for _, name := range order {
		out, _ := exec.Command(exe, name).CombinedOutput()
		fmt.Printf("%-15s %s\n", name, strings.ReplaceAll(strings.SplitN(strings.TrimSpace(string(out)), "\n\n", 2)[0], "\n", " | "))
	}
	fmt.Println("== 1.3 say with channel")
	done := make(chan struct{})
	go say("goroutine", done)
	<-done
}
