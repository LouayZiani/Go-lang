// so i created a WCAG colour contrast checker. It creates a few tasks, processes them
// and prints the results.
package main

import (
	"fmt"
	"os"

	"sis/internal/task"
)

func main() {
	svc := task.NewService(task.WCAGProcessor{})

	inputs := []task.Input{
		{Foreground: "#000000", Background: "#FFFFFF", PointSize: 12},
		{Foreground: "#777777", Background: "#FFFFFF", PointSize: 12},
		{Foreground: "#777777", Background: "#FFFFFF", PointSize: 18},
		{Foreground: "#FF0000", Background: "#FFFFFF", PointSize: 12, Bold: true},
		{Foreground: "not-a-colour", Background: "#FFFFFF", PointSize: 12}, // rejected
	}

	for _, in := range inputs {
		t, err := svc.Create(in)
		if err != nil {
			fmt.Fprintln(os.Stderr, "create:", err)
			continue
		}
		if _, err := svc.Process(t.ID); err != nil {
			fmt.Fprintln(os.Stderr, "process:", err)
		}
	}

	for _, t := range svc.List() {
		fmt.Printf("%s  %s on %s  %gpt bold=%-5v  %-9s  ratio %5.2f  %s\n",
			t.ID, t.Input.Foreground, t.Input.Background, t.Input.PointSize, t.Input.Bold,
			t.State, t.Result.Ratio, t.Result.Level())
	}

	if _, err := svc.Process("task-1"); err != nil {
		fmt.Println("expected error:", err)
	}
	if _, err := svc.Get("task-42"); err != nil {
		fmt.Println("expected error:", err)
	}
}
