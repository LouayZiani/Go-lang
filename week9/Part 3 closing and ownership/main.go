package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func main() {
	exe, _ := os.Executable()
	if len(os.Args) > 1 {
		wrong(3)
		return
	}
	fmt.Println("== 3.1 WRONG version x3")
	for i := 0; i < 3; i++ {
		out, _ := exec.Command(exe, "wrong").CombinedOutput()
		fmt.Println(strings.ReplaceAll(strings.SplitN(string(out), "\n\n", 2)[0], "\n", " | "))
	}
	fmt.Println("== 3.1 collect:", sum(collect(3)))
	fmt.Println("== 3.2 pipeline:", sum(square(generator(1, 2, 3, 4, 5, 6, 7, 8, 9, 10))))
	fmt.Println("== 3.3 with filterEven:", sum(filterEven(square(generator(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)))))
	fmt.Println("== Task 4")
	fmt.Println("same:", Same(New(1), New(1), true), "| different:", Same(New(1), New(2), true))
	time.Sleep(50 * time.Millisecond)
	fmt.Println("goroutines with done:", runtime.NumGoroutine())
	fmt.Println("different, no done:", Same(New(1), New(2), false))
	time.Sleep(50 * time.Millisecond)
	fmt.Println("goroutines without done:", runtime.NumGoroutine())
}
