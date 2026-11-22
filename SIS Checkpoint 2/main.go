package main

import (
	"contrastchecker/internal/task"
	"errors"
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	svc := task.NewService(task.WCAGProcessor{})

	pool, err := task.NewPool(svc, task.DefaultWorkers, task.DefaultQueueCapacity)
	if err != nil {
		return err
	}
	//stop is idempotent so deferring it is a safe backstop for early returns
	defer pool.Stop()

	fmt.Printf("pool: %d workers, queue capacity %d\n\n", pool.Workers(), pool.QueueCapacity())

	inputs := []task.Input{
		{Foreground: "#000000", Background: "#ffffff", PointSize: 12},
		{Foreground: "#767676", Background: "#ffffff", PointSize: 12},
		{Foreground: "#777777", Background: "#ffffff", PointSize: 12},
		{Foreground: "#777777", Background: "#ffffff", PointSize: 18},
		{Foreground: "#777", Background: "#fff", PointSize: 14, Bold: true},
		{Foreground: "#0b5fff", Background: "#ffffff", PointSize: 16},
		{Foreground: "#ffcc00", Background: "#ffffff", PointSize: 16},
		{Foreground: "#ffffff", Background: "#1a1a2e", PointSize: 11},
	}

	//create then submit. Submit returns as soon as the task is queued
	ids := make([]string, 0, len(inputs))
	for _, in := range inputs {
		tk, err := svc.Create(in)
		if err != nil {
			return err
		}
		if err := pool.Submit(tk.ID); err != nil {
			if errors.Is(err, task.ErrQueueFull) {
				fmt.Printf("%s: queue full, try again later\n", tk.ID)
				continue
			}
			return err
		}
		ids = append(ids, tk.ID)
	}

	//wait for each accepted task(no sleeping, no polling)
	for _, id := range ids {
		done, err := svc.Done(id)
		if err != nil {
			return err
		}
		<-done
	}

	for _, tk := range svc.List() {
		printTask(tk)
	}

	//shut down: no new tasks accepted, accepted ones finished, and workers gone
	pool.Stop()
	fmt.Printf("\nworkers still running after Stop: %d\n", pool.Running())

	if err := pool.Submit(ids[0]); errors.Is(err, task.ErrPoolStopped) {
		fmt.Println("submit after Stop:", err)
	}
	return nil
}

func printTask(tk task.Task) {
	if tk.Result == nil {
		fmt.Printf("%-7s %-9s %s on %s  %v\n", tk.ID, tk.State, tk.Input.Foreground, tk.Input.Background, tk.Error)
		return
	}
	r := tk.Result
	fmt.Printf("%-7s %-9s %-8s on %-8s %4.1fpt  ratio %5.2f  large=%-5v AA=%-5v AAA=%-5v -> %s\n",
		tk.ID, tk.State, tk.Input.Foreground, tk.Input.Background, tk.Input.PointSize,
		r.Ratio, r.LargeText, r.PassAA, r.PassAAA, r.Level)
}
