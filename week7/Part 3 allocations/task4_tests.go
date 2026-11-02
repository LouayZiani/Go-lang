package main

import "testing"

var (
	sinkP  Point
	sinkPP *Point
	sinkI  []int
)

func BenchmarkNewPointValue(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkP = NewPointValue(i, i)
	}
}

func BenchmarkNewPointPointer(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkPP = NewPointPointer(i, i)
	}
}

func BenchmarkCollectNoCap(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkI = CollectNoCap(1000)
	}
}

func BenchmarkCollectWithCap(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkI = CollectWithCap(1000)
	}
}

func TestAllocCounts(t *testing.T) {
	check := func(name string, want float64, f func()) {
		t.Helper()
		got := testing.AllocsPerRun(1000, f)
		t.Logf("%-20s %v allocs/op", name, got)
		if got != want {
			t.Errorf("%s: %v allocs, want %v", name, got, want)
		}
	}
	check("NewPointValue", 0, func() { sinkP = NewPointValue(1, 2) })
	check("NewPointPointer", 1, func() { sinkPP = NewPointPointer(1, 2) })
	check("CollectNoCap(1000)", 12, func() { sinkI = CollectNoCap(1000) })
	check("CollectWithCap(1000)", 1, func() { sinkI = CollectWithCap(1000) })
}

func TestCollectSame(t *testing.T) {
	a, b := CollectNoCap(1000), CollectWithCap(1000)
	if len(a) != 1000 || len(b) != 1000 {
		t.Fatal("wrong length")
	}
	for i := range a {
		if a[i] != b[i] || a[i] != i {
			t.Fatalf("mismatch at %d", i)
		}
	}
}

func TestCapGrowth(t *testing.T) {
	var s []int
	last := -1
	total := 0
	n := 0
	for i := 0; i < 1000; i++ {
		s = append(s, i)
		if cap(s) != last {
			last = cap(s)
			n++
			total += last * 8
			t.Logf("alloc #%2d: cap=%4d  (%5d bytes)", n, last, last*8)
		}
	}
	t.Logf("total: %d allocations, %d bytes", n, total)
}
