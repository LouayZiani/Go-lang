package task7

import (
	"bytes"
	"testing"
)

var (
	ids200 = func() []int {
		s := make([]int, 200)
		for i := range s {
			s[i] = i + 1
		}
		return s
	}()
	sinkN int
)

func TestRenderSame(t *testing.T) {
	want := RenderNew(ids200)
	for i := 0; i < 10; i++ {
		if got := RenderPool(ids200); got != want {
			t.Fatalf("call %d: got %d, want %d", i, got, want)
		}
	}

	var b bytes.Buffer
	writeIDs(&b, []int{7, 42})
	if b.String() != "id=7\nid=42\n" {
		t.Fatalf("got %q", b.String())
	}
}

func TestRenderPoolAllocs(t *testing.T) {
	RenderPool(ids200) // warm up the pool
	pool := testing.AllocsPerRun(1000, func() { sinkN = RenderPool(ids200) })
	fresh := testing.AllocsPerRun(1000, func() { sinkN = RenderNew(ids200) })
	t.Logf("RenderNew: %v allocs/op, RenderPool: %v allocs/op", fresh, pool)
	if pool >= fresh {
		t.Errorf("pool (%v) should allocate less than new (%v)", pool, fresh)
	}
}

func BenchmarkRenderNew(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkN = RenderNew(ids200)
	}
}

func BenchmarkRenderPool(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkN = RenderPool(ids200)
	}
}

func BenchmarkRenderNewParallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sinkN = RenderNew(ids200)
		}
	})
}

func BenchmarkRenderPoolParallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sinkN = RenderPool(ids200)
		}
	})
}
