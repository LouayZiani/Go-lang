package main

import (
	"errors"
	"fmt"
)

type Task struct {
	Name string
	Run  func() error
}

func runOne(task Task) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("task %s: panic: %v", task.Name, r)
		}
	}()

	if task.Run == nil {
		return fmt.Errorf("task %s: nil Run function", task.Name)
	}

	if err := task.Run(); err != nil {
		return fmt.Errorf("task %s: %w", task.Name, err)
	}

	return nil
}

func SafeRun(tasks []Task) []error {
	var errs []error

	for _, task := range tasks {
		if err := runOne(task); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

var ErrSentinel = errors.New("expected failure")

func main1111() {
	tasks := []Task{
		{
			Name: "success",
			Run: func() error {
				fmt.Println("success task ran")
				return nil
			},
		},
		{
			Name: "index",
			Run: func() error {
				values := []int{1}
				fmt.Println(values[5])
				return nil
			},
		},
		{
			Name: "nil-map",
			Run: func() error {
				var m map[string]int
				m["x"] = 1
				return nil
			},
		},
		{
			Name: "config",
			Run: func() error {
				panic("config is missing")
			},
		},
		{
			Name: "sentinel",
			Run: func() error {
				return ErrSentinel
			},
		},
	}

	errs := SafeRun(tasks)

	for _, err := range errs {
		fmt.Println(err)
	}

	joined := errors.Join(errs...)
	fmt.Println("sentinel preserved:", errors.Is(joined, ErrSentinel))
}
