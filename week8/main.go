package main

import (
	"fmt"
	"os"
)

func main() {
	task := "help"
	if len(os.Args) > 1 {
		task = os.Args[1]
	}
	switch task {
	case "task1":
		runTask1()
	case "task2":
		runTask2()
	case "racebad":
		runRaceBad()
	case "task3":
		runTask3()
	case "task4":
		runTask4()
	case "task5":
		runTask5()
	case "task6":
		runTask6()
	case "task6-leaky":
		runTask6Variant("leaky", handleLeaky)
	case "task6-buffered":
		runTask6Variant("buffered", handleBuffered)
	case "task6-context":
		runTask6Variant("context", handleContext)
	case "preempt":
		runPreemption()
	case "threads":
		runThreads()
	case "task7":
		fmt.Println("Task 7 is a benchmark/test task. Run: go test -run TestCountPar -bench 'Count(Seq|Par)' -benchmem -cpu 1,2")
	case "all":
		runTask1()
		runTask2()
		runTask3()
		runTask4()
		runTask5()
		runTask6()
		fmt.Println("\nTask 7: run go test -run TestCountPar -bench 'Count(Seq|Par)' -benchmem -cpu 1,2")
	default:
		fmt.Println("Go Concurrency Lab")
		fmt.Println("Usage: go run . task1|task2|racebad|task3|task4|task5|task6|task6-leaky|task6-buffered|task6-context|task7|preempt|threads|all")
		fmt.Println("Run tests/race detector: go test ./... ; go test -race ./...")
	}
}
