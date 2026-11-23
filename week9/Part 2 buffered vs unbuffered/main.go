package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func run(size int) [3]time.Duration {
	start := time.Now()
	ch := make(chan int, size)
	var t [3]time.Duration
	done := make(chan struct{})
	go func() {
		for i := range t {
			ch <- i + 1
			t[i] = time.Since(start).Round(100 * time.Millisecond)
		}
		close(done)
	}()
	time.Sleep(300 * time.Millisecond)
	for range t {
		<-ch
		time.Sleep(100 * time.Millisecond)
	}
	<-done
	return t
}

func snippet(n int) {
	ch := make(chan string, n)
	ch <- "a"
	ch <- "b"
	fmt.Println(<-ch, <-ch)
}

func main() {
	if len(os.Args) > 1 {
		snippet(1)
		return
	}
	fmt.Println("== 2.1 times of the three sends")
	for _, size := range []int{0, 1, 2, 3, 10} {
		fmt.Printf("cap %-3d %v\n", size, run(size))
	}
	fmt.Println("== 2.2 cap 2")
	snippet(2)
	fmt.Println("== 2.2 cap 1")
	exe, _ := os.Executable()
	out, _ := exec.Command(exe, "x").CombinedOutput()
	fmt.Println(strings.SplitN(string(out), "\n\n", 2)[0])
}
