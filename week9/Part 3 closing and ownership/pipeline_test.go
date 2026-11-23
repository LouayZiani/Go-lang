package main

import "testing"

func TestCollect(t *testing.T) {
	if got := sum(collect(3)); got != 189 {
		t.Fatal(got)
	}
}

func TestPipeline(t *testing.T) {
	if got := sum(square(generator(1, 2, 3, 4, 5, 6, 7, 8, 9, 10))); got != 385 {
		t.Fatal(got)
	}
	if got := sum(filterEven(square(generator(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)))); got != 220 {
		t.Fatal(got)
	}
}

func TestSame(t *testing.T) {
	if !Same(New(1), New(1), true) || Same(New(1), New(2), true) {
		t.Fatal("same")
	}
}
