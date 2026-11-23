package main

import (
	"slices"
	"testing"
	"time"
)

func TestPass(t *testing.T) {
	for _, c := range cases {
		if got := c.f(); got != N*(N-1)/2 {
			t.Fatal(c.name, got)
		}
	}
}

func TestMerge(t *testing.T) {
	var got []int
	for v := range merge(feed(time.Millisecond, 1, 2, 3), feed(time.Millisecond, 100, 101)) {
		got = append(got, v)
	}
	slices.Sort(got)
	if !slices.Equal(got, []int{1, 2, 3, 100, 101}) {
		t.Fatal(got)
	}
}

func BenchmarkPass(b *testing.B) {
	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				c.f()
			}
		})
	}
}
