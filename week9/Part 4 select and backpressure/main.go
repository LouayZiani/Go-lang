package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func bonus(mode string) {
	ch := make(chan int)
	if mode == "closed" {
		close(ch)
	} else {
		ch = nil
	}
	if mode == "nilnodefault" {
		select {
		case v, ok := <-ch:
			fmt.Println("closed:", v, ok)
		}
		return
	}
	select {
	case v, ok := <-ch:
		fmt.Println("closed:", v, ok)
	default:
		fmt.Println("default")
	}
}

func main() {
	if len(os.Args) > 1 {
		bonus(os.Args[1])
		return
	}
	task5()
	task6()
	task7()
	exe, _ := os.Executable()
	for _, m := range []string{"closed", "nil", "nilnodefault"} {
		out, _ := exec.Command(exe, m).CombinedOutput()
		fmt.Printf("== bonus %s: %s\n", m, strings.SplitN(strings.TrimSpace(string(out)), "\n\n", 2)[0])
	}
}
