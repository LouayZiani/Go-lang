# Contrast Checker (Checkpoint 1)

A small in-memory Go service that checks whether text is readable on a
background, following the WCAG 2.x contrast rules. Accessibility is the
domain on purpose: a "task" here is one contrast check.

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

## Task lifecycle

```
            Process()              success
 Pending ───────────────► Processing ───────────► Completed
                               │
                               │ processor error
                               ▼
                            Failed
```

Allowed transitions:

| From       | To                    |
|------------|-----------------------|
| Pending    | Processing            |
| Processing | Completed, Failed     |
| Completed  | none (terminal)       |
| Failed     | none (terminal)       |

Anything else returns `ErrInvalidTransition` and leaves the task unchanged,
e.g. processing a task twice.

## Errors

| Error                  | When                                             |
|------------------------|--------------------------------------------------|
| `ErrInvalidInput`      | `Create` with bad input                          |
| `ErrTaskNotFound`      | `Get` / `Process` with an unknown ID             |
| `ErrInvalidTransition` | `Process` on a task that isn't `Pending`         |

Check them with `errors.Is`. No panics for expected errors.

## Layout

```
main.go                     demo
internal/contrast/          pure colour maths (parse, luminance, ratio)
internal/task/              Task, State, Input/Result, Processor interface, Service
```

- `Processor` is the interface: `Service` depends on it, `WCAGProcessor` implements it, tests use a failing fake.
- `Service` keeps tasks in a `map[string]*Task` plus an ID slice for ordered `List()`.
- `Get`, `List` and `Process` return **copies** (including a copied `Result`), so callers can't modify stored tasks.
- `Service` is not goroutine-safe yet; that belongs to the concurrency checkpoint.

## Run and test

```
go run .
gofmt -l .
go vet ./...
go test ./...
```
