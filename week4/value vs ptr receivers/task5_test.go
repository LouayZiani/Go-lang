package main

import "testing"

// package-level sink so the compiler can't throw the result away
var sink float64

func BenchmarkFirstValue(b *testing.B) {
	r := Report{Title: "x"}
	var s float64
	for i := 0; i < b.N; i++ {
		s += r.FirstValue() // 8 KB copied on every call
	}
	sink = s
}

func BenchmarkFirstPointer(b *testing.B) {
	r := &Report{Title: "x"}
	var s float64
	for i := 0; i < b.N; i++ {
		s += r.FirstPointer() // only an 8-byte pointer is passed
	}
	sink = s
}

func BenchmarkPointValue(b *testing.B) {
	p := Point{1, 2}
	var s float64
	for i := 0; i < b.N; i++ {
		s += p.SumValue()
	}
	sink = s
}

func BenchmarkPointPointer(b *testing.B) {
	p := &Point{1, 2}
	var s float64
	for i := 0; i < b.N; i++ {
		s += p.SumPointer()
	}
	sink = s
}
