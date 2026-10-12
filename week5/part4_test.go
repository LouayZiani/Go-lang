package main

import (
	"testing"
)

func SumInts(nums []int) int {
	var s int
	for _, n := range nums {
		s += n
	}
	return s
}

func SumGeneric[T Number](nums []T) T {
	var s T
	for _, n := range nums {
		s += n
	}
	return s
}

func SumAny(nums []any) int {
	total := 0
	for _, n := range nums {
		total += n.(int)
	}
	return total
}

func ToAny(nums []int) []any {
	anySlice := make([]any, len(nums))
	for i, v := range nums {
		anySlice[i] = v
	}
	return anySlice
}

func BenchmarkSumInts(b *testing.B) {
	nums := make([]int, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SumInts(nums)
	}
}

func BenchmarkSumGeneric(b *testing.B) {
	nums := make([]int, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SumGeneric(nums)
	}
}

func BenchmarkSumAny(b *testing.B) {
	nums := make([]int, 1000)
	anyNums := ToAny(nums)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SumAny(anyNums)
	}
}

func BenchmarkSumAnyWithBoxing(b *testing.B) {
	nums := make([]int, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SumAny(ToAny(nums))
	}
}
