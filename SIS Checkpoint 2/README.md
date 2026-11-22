# Contrast Checker (Checkpoint 2: Concurrent Processing)

A small in-memory Go service that checks whether text is readable on a
background, following the WCAG 2.x contrast rules. Accessibility is the
domain on purpose: a "task" here is one contrast check.

Checkpoint 1 was the synchronous version. Checkpoint 2 adds a **worker pool**
with a **bounded queue**, so many checks can run at the same time. Everything
from Checkpoint 1 (validation, lifecycle, errors, tests) still works.

## What a task is

| Part       | In this project                                                         |
|------------|-------------------------------------------------------------------------|
| Input      | foreground colour, background colour (`#RGB` / `#RRGGBB`), point size, bold |
| Processing | compute both luminances, the contrast ratio, whether the text is "large", and the AA / AAA verdicts |
| Result     | ratio (1-21), large-text flag, AA pass/fail, AAA pass/fail, level (`AAA` / `AA` / `Fail`) |

WCAG thresholds used:

| Text                          | AA  | AAA |
|-------------------------------|-----|-----|
| Normal                        | 4.5 | 7   |
| Large (>=18pt, or >=14pt bold)| 3   | 4.5 |

## Validation

`Create` rejects (with `ErrInvalidInput`, wrapped with the reason): empty or
malformed colours, non-hex digits, and point sizes that are zero, negative,
NaN or infinite. Rejected input is never stored.

## How tasks are created and submitted

```go
svc := task.NewService(task.WCAGProcessor{})
pool, _ := task.NewPool(svc, 4, 16)   // 4 workers, queue capacity 16

t, _ := svc.Create(input)             // Pending
err := pool.Submit(t.ID)              // returns right away, task is Queued
done, _ := svc.Done(t.ID)             // channel closed when the task finishes
<-done
t, _ = svc.Get(t.ID)                  // Completed, with Result

pool.Stop()                           // stop accepting, finish accepted, wait for workers
```

`Submit` never waits for processing. `Service.Process(id)` from Checkpoint 1
still exists and runs a `Pending` task synchronously in the caller's goroutine.

## Concurrency design

```
                     Pool
                     ┌───────────────────────────────────────────┐
 Create ─► Service   │  Submit ─► queue (chan string, cap N) ─┬─► worker 1 ─┐
 (tasks map,         │                                        ├─► worker 2 ─┤ move / finish
  one RWMutex) ◄─────┼────────────────────────────────────────┴─► worker W ─┘ (short locks)
                     └───────────────────────────────────────────┘
```

- **Workers:** a fixed number (minimum 2, the demo uses `DefaultWorkers = 4`),
  started by `NewPool`. Each worker loops `for id := range queue`, so it is
  reused for many tasks. There is no goroutine per task.
- **Queue:** a buffered `chan string` of task IDs with a fixed capacity
  (the demo uses `DefaultQueueCapacity = 16`). Task data stays in the `Service`;
  the queue only carries IDs.
- **Processing runs outside the lock.** A worker locks briefly to move the task
  to `Processing`, unlocks, runs the processor, then locks briefly again to
  store the result. Reads (`Get`, `List`) are never blocked by slow processing.
- **Order is not guaranteed.** Tasks may finish in any order; nothing depends on it.
- The `Processor` must be safe for concurrent use. `WCAGProcessor` has no state,
  so it is.

### What happens when the queue is full

`Submit` **rejects** instead of blocking: it returns `ErrQueueFull` and the task
goes back to `Pending` (it is never left `Queued` or `Processing`). The caller
can resubmit later. Because `Submit` never blocks, a caller can't get stuck
because the pool stopped.

### Task lifecycle

```
            Submit()              worker picks up          success
 Pending ───────────────► Queued ───────────────► Processing ───────────► Completed
    │  ▲                    │                          │
    │  └────────────────────┘                          │ processor error
    │   queue was full                                 ▼
    │                                                Failed
    │ Process() (synchronous, Checkpoint 1)
    └──────────────────────────────────────────► Processing
```

Allowed transitions (the single source of truth is `CanTransition`):

| From       | To                    | Why                                   |
|------------|-----------------------|---------------------------------------|
| Pending    | Queued                | `Submit` accepted the task            |
| Pending    | Processing            | synchronous `Process`                 |
| Queued     | Processing            | a worker took it from the queue       |
| Queued     | Pending               | rollback: the queue turned out to be full |
| Processing | Completed, Failed     | processing finished                   |
| Completed  | none (terminal)       |                                       |
| Failed     | none (terminal)       |                                       |

Anything else returns `ErrInvalidTransition` and leaves the task unchanged:
submitting a task that is `Queued`, `Processing`, `Completed` or `Failed`, or
calling `Process` on a task that isn't `Pending`.

### How shared data is synchronized

| Shared data | Protected by |
|-------------|--------------|
| `Service`: tasks map, creation order, ID counter, and **every field of every stored task** (state, result, error) | one `sync.RWMutex` (`Service.mu`) |
| `Pool.stopped` and the rule "never send on the queue after it is closed" | `sync.RWMutex` (`Pool.mu`) |
| Queue contents | the channel itself |
| Count of live workers (`Running`) | `atomic.Int32` |

Key points:

- Locking the map alone would not be enough, because workers change the tasks
  *inside* it. So the one lock covers the map and the task fields.
- **Atomic check-and-change:** `move(id, from, to)` checks the state and changes
  it under one write lock. Only one goroutine can win `Pending -> Queued` or
  `Queued -> Processing`, so a task is never queued twice and never processed
  by two workers.
- `Get`, `List` and `Process` return **copies** (including a copied `Result`),
  so nobody outside holds a pointer into shared data.
- **Lock order** is always `Pool.mu` then `Service.mu`, and `Service` never
  calls back into `Pool`, so there's no lock cycle. No lock is held during
  processing or during a blocking channel operation.

### How workers are stopped

`Pool.Stop()`:

1. takes the write lock, sets `stopped = true` and **closes the queue** (only
   `Stop` ever closes it, and only once);
2. new `Submit` calls now return `ErrPoolStopped`;
3. workers keep receiving until the closed queue is empty, so every accepted
   task is still processed (a closed channel still hands out what's buffered);
4. `Stop` waits on a `sync.WaitGroup` until all workers have exited.

Why there is no "send on closed channel" panic: `Submit` does its check and its
send while holding the **read** lock, and `Stop` closes the channel while
holding the **write** lock, so they can never overlap. `Submit`'s send is
non-blocking, so it can't hold the read lock forever and starve `Stop`.

Calling `Stop` again (or from several goroutines) is safe: only the first call
closes the queue, every call waits for the workers and returns.

## Errors

| Error                  | When                                                  |
|------------------------|-------------------------------------------------------|
| `ErrInvalidInput`      | `Create` with bad input                               |
| `ErrTaskNotFound`      | `Get` / `Process` / `Submit` / `Done` with an unknown ID |
| `ErrInvalidTransition` | wrong state for the operation (e.g. submit twice)     |
| `ErrQueueFull`         | `Submit` when the queue has no free slot              |
| `ErrPoolStopped`       | `Submit` after `Stop`                                 |
| `ErrInvalidConfig`     | `NewPool` with fewer than 2 workers, capacity < 1, or nil service |

Check them with `errors.Is`. No panics for expected errors.

## Layout

```
main.go                      demo
internal/contrast/           pure colour maths (parse, luminance, ratio)
internal/task/
    task.go                  Task, State, transitions, Input/Result, errors
    processor.go             Processor interface + WCAGProcessor
    service.go               Service: in-memory store, state machine, RWMutex
    pool.go                  Pool: bounded queue, workers, Submit, Stop
    *_test.go                unit + concurrency tests
```

## Tests

Concurrency tests don't depend on finishing order and don't rely on sleeps.
They use channels and a "gate" processor that holds workers in the middle of a
task so the test can inspect the system. The only timeouts are safety nets so a
broken implementation fails instead of hanging.

| Required behaviour | Test |
|--------------------|------|
| Pool startup | `TestPoolStartup`, `TestNewPoolConfig` |
| Async submission | `TestSubmitDoesNotWaitForProcessing` |
| Many tasks, correct result in the correct task | `TestSubmitAndProcessMany`, `TestRealProcessorThroughPool` |
| Two tasks really overlapping | `TestTasksProcessConcurrently` |
| Concurrency capped by worker count | `TestConcurrencyIsBoundedByWorkers` |
| States after processing, failures | `TestSubmitAndProcessMany`, `TestFailedTaskThroughPool` |
| Invalid / duplicate submissions | `TestSubmitInvalid`, `TestSubmitQueuedTaskTwice`, `TestConcurrentSubmitSameTask` |
| Full queue | `TestFullQueueRejectsAndRecovers` |
| Safe concurrent access | `TestConcurrentReadsWhileProcessing`, `TestConcurrentCreate`, `TestConcurrentProcessSameTaskOnlyOneWins` |
| Worker termination | `TestStopWaitsForWorkersToExit` |
| Accepted tasks not lost on shutdown | `TestStopFinishesAcceptedTasks`, `TestConcurrentSubmitAndStop` |
| Rejection after shutdown, double `Stop` | `TestSubmitAfterStopIsRejected`, `TestStopTwiceAndConcurrently` |

## Run and test

```
go run .
gofmt -l .
go vet ./...
go test ./...
go test -race ./...
```
