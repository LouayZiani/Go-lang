# Go Concurrency Lab

This project implements Tasks 1–7. It uses only the Go standard library and targets Go 1.23+ (the handout's `WaitGroup.Go` shortcut is not used, so it also works on Go 1.23).

## Run

```sh
go run . task1
go run . task2
go run -race . racebad  # intentionally incorrect example; run separately
go run . task3
go run . task4
go run . task5
go run . task6
go run . task6-leaky    # isolated Task 6 leaky baseline
go run . task6-buffered # isolated buffered-channel fix
go run . task6-context  # isolated context-cancellation fix
go test ./...
go test -race ./...
go test -bench . -benchmem -cpu 1,2
```

On macOS/Linux, compare Task 1 with one P:

```sh
GOMAXPROCS=1 go run . task1
GOMAXPROCS=1 go run . task4
GODEBUG=schedtrace=500 go run . task4
GODEBUG=asyncpreemptoff=1 go run . preempt
```

PowerShell equivalents:

```powershell
$env:GOMAXPROCS=1; go run . task1
$env:GODEBUG='schedtrace=500'; go run . task4
$env:GODEBUG='asyncpreemptoff=1'; go run . preempt
```

For thread creation comparison on Linux/macOS/FreeBSD:

```sh
go run . threads
```

The `preempt` and `threads` commands are available directly. The preemption demonstration intentionally leaves a busy goroutine running, so run it as a separate process. With `GODEBUG=asyncpreemptoff=1`, it may not return; terminate that process if necessary.

## Task map

- **Task 1:** `countPrimes`, `cpuWork`, `ioWork`; `task1` prints CPU/I/O timings and core counts. Run again with `GOMAXPROCS=1` and fill your handout's table with your machine's measurements.
- **Task 2:** `say` waits for both goroutines using a `WaitGroup`; demonstrates `atomic.Int64`, independent slice slots, and recovering a panic inside the goroutine that panics. `racebad` intentionally reproduces the incorrect `WaitGroup.Add` placement and shared-counter race; run it separately with `-race`.
- **Task 3:** `task3` starts 100,000 channel-blocked goroutines and reports `NumGoroutine`, `StackSys`, and start time. Benchmarks compare goroutine start/stop with a direct call.
- **Task 4:** `task4` launches eight CPU-busy goroutines. Use `GODEBUG=schedtrace=500`. `threads` compares `time.Sleep` and blocking `Nanosleep` using the `threadcreate` profile on supported Unix platforms. `preempt` is available via the small main-switch addition described above.
- **Task 5:** `fetchLeaky` / `startTickerLeaky` demonstrate leaks; `fetchFixed` uses a one-slot buffered channel; `startTicker` observes context cancellation. `task5` prints counts and a goroutine profile. `TestFetchFixedDoesNotLeak` checks cleanup.
- **Task 6:** `handleLeaky`, `handleBuffered`, and `handleContext` show the unbuffered leak, buffered-channel fix, and context-cancellation fix. `task6` runs all three approaches; for clean independent counters, run `task6-leaky`, `task6-buffered`, and `task6-context` in separate processes.
- **Task 7:** `CountSeq`, `CountPar`, `TestCountPar`, and the benchmarks are in `count.go` and `lab_test.go`.

## Important notes when interpreting results

1. Timings, `StackSys`, heap size, and thread counts depend on your OS, Go version, CPU, and background load. Record your own output; do not copy the sample numbers from the handout as if they were measured on your laptop.
2. `GOMAXPROCS` limits how many goroutines can execute Go code simultaneously; it does not limit how many goroutines exist.
3. The intentional leak demos leave goroutines alive in that process. Run each task command separately, or use a fresh process for clean measurements.
4. The handout's `worker 5` ordering is not guaranteed by the Go specification. Scheduler `runnext` behavior may make it appear frequently, but application correctness must never depend on it.
5. `runtime.MemStats.StackSys` is a runtime-wide stack-memory measurement, so its delta is approximate and may include runtime activity; it is not a precise per-goroutine allocation measurement.
6. The context-aware `buildReport` can stop its work when the request times out. A buffered channel alone only prevents the sender from becoming stuck; it does not cancel unnecessary work.


## Participation tasks and defence notes

- **Concurrency without parallelism:** a single-core web server has many request goroutines. One goroutine waits for a database response while another handles a ready request; they make progress over overlapping time, but only one executes Go code at a time when `GOMAXPROCS=1`.
- **Concurrency with parallelism:** on a multi-core server, two independent CPU-heavy requests (for example, image transformations) can execute on different cores at the same time.
- **Why I/O still overlaps with one P:** waiting goroutines are parked; the scheduler runs another goroutine. `time.Sleep` is handled by the runtime timer and does not occupy an OS thread for the whole wait.
- **Why CPU work does not scale with one P:** only one goroutine can execute Go code at once. Extra goroutines add scheduling overhead but no additional CPU capacity.
- **Why 4 workers may not beat 2 on 2 cores:** two workers already keep both cores busy; extra workers add scheduling/coordination overhead. Two workers are not exactly twice as fast because ranges differ in difficulty, and there is startup, scheduling, and aggregation overhead.
- **Why `main` output can omit goroutine lines:** returning from `main` exits the process; it does not wait for other goroutines. Sleep only changes the scheduling opportunity; it is not a synchronization guarantee. Use `WaitGroup`.
- **Why `wg.Add` comes before `go`:** otherwise `Wait` can observe a zero counter and return before the worker increments it. Call `Done` with `defer` in each worker.
- **Data race:** concurrent unsynchronized access to the same memory location, with at least one write. Use `go test -race` / `go run -race`; protect shared state with a mutex/atomic or give each goroutine a distinct slot.
- **Panic recovery:** `recover` only catches a panic in the same goroutine. Put `defer` with `recover` inside the goroutine that may panic.
- **G/M/P:** G is a goroutine, M is an OS thread, P is the scheduler resource/permission to execute Go code. There are `GOMAXPROCS` Ps. An idle P can steal runnable goroutines from another P's queue. A blocking syscall blocks its M; the runtime can hand the P to another M.
- **Preemption:** since Go 1.14, asynchronous preemption lets the runtime interrupt long-running goroutines that do not reach function calls, improving scheduler responsiveness.
- **Worker 5 puzzle:** the most recently created goroutine may be placed in a P's `runnext` slot and run promptly. That is an implementation detail, not a language guarantee; never rely on print order.
- **Goroutine leak:** a goroutine remains blocked or running after its result is no longer needed. Common causes are an unreceived unbuffered send and a worker without a cancellation path. Find them with `runtime.NumGoroutine()` and `pprof.Lookup("goroutine")`.
- **Why buffered channel alone is not enough:** it lets a completed worker deliver one result and exit, but it still performs work after the caller timed out. Passing a context lets the worker stop early.
- **Task 7:** `parts=2` may not double speed because overhead and uneven work remain. More chunks can improve load balancing, up to a point. `parts=200000` creates one goroutine per integer and huge allocation/scheduling overhead, so it can be slower than sequential counting.

## Local benchmark evidence (this build environment only)

A short one-iteration smoke benchmark was run with Go 1.23.2 in a container reporting `runtime.NumCPU()=3` on an AMD EPYC 9V74 host. It is only a build/benchmark check, **not your laptop's result**. In that run `CountSeq` was about 15.2 ms; `parts=2` about 14.7 ms; `parts=8` about 14.8 ms; `parts=1000` about 8.9 ms; and `parts=200000` about 64.8 ms. Use the full benchmark command on your own machine and report its output.

- **Bonus loop-variable question:** on Go 1.22+, each iteration has its own loop variable, so the five goroutines print values 1 through 5, each once. The order is unspecified even with `GOMAXPROCS=1`; one P means one goroutine executes Go code at a time, not that goroutines run in creation order.
