package main

// each heap allocation costs twice: once when memory is allocated, and again later when the GC has to find it and
// free it. Fewer allocs --> less GC work --> faster programs

// methods to count allocations:
// 1. Benchmarking : go test -bench=. -benchmem
// 2. testing.AllocsPerRun: this lets us count allocations for a particular fct inside a test
// 3. runtime.ReadMemStats: lets us inspect memory statistics while a program is running
