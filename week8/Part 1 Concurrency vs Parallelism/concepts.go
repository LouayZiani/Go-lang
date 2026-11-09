package main

//goroutine: thread of execution managed by the Go runtime rather than OS
// Concurrency (about structure): a program is made of independent tasks that can make progress w/o waiting for each other
// Parallelism (about execution): several tasks run at the same time on diff CPU cores
// --> Concurrency is abt dealing with a lot of things at once, and Parallelism is about doing a lot of things at once
// in concurrency, tasks take turns on 1 CPU
// Concurrency: multiple tasks taking turns making progress
// Parallelism: multiple tasks execute simultaneously

// CPU-bound work: more goroutines can help when there are additional CPU cores available to execute the calculations
// I/O-bound work: more goroutines can help even with one CPU core cuz while one goroutine waits, another makes prgress

// GOMAXPROCS: controls how many Go tasks can execute Go code simultaenously
// e.g. GOMAXPROCS=1 --> only one goroutine executes Go code at a time
