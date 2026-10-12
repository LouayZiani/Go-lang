package main

import (
	"cmp"
	"fmt"
)

type Number interface {
	~int | ~int64 | ~float64
}

// task 5-1: Sum and Average
func Sum[T Number](nums []T) T {
	var s T
	for _, n := range nums {
		s += n
	}
	return s
}

func Average[T Number](nums []T) (float64, bool) {
	if len(nums) == 0 {
		return 0, false
	}
	total := Sum(nums)
	// we convrt total to float64 for div
	var sumFloat float64
	switch v := any(total).(type) {
	case int:
		sumFloat = float64(v)
	case int64:
		sumFloat = float64(v)
	case float64:
		sumFloat = v
	}
	return sumFloat / float64(len(nums)), true
}

// task5-2: MaxOf
func MaxOf[T cmp.Ordered](items []T) (T, bool) {
	var zero T
	if len(items) == 0 {
		return zero, false
	}
	m := items[0]
	for _, v := range items[1:] {
		if v > m {
			m = v
		}
	}
	return m, true
}

// taks5-3: IndexOf
func IndexOf[T comparable](items []T, target T) int {
	for i, v := range items {
		if v == target {
			return i
		}
	}
	return -1
}

// task5-4: Filter
func Filter[T any](items []T, keep func(T) bool) []T {
	var result []T
	for _, item := range items {
		if keep(item) {
			result = append(result, item)
		}
	}
	return result
}

type Score int

// task7: Tour of Go -> Generic types List[T]
type Node[T any] struct {
	val  T
	next *Node[T]
}

type List[T any] struct {
	head *Node[T]
}

func (l *List[T]) Push(v T) {
	newNode := &Node[T]{val: v}
	if l.head == nil {
		l.head = newNode
		return
	}
	curr := l.head
	for curr.next != nil {
		curr = curr.next
	}
	curr.next = newNode
}

func (l *List[T]) Len() int {
	count := 0
	curr := l.head
	for curr != nil {
		count++
		curr = curr.next
	}
	return count
}

func (l *List[T]) String() string {
	var sb stringsBuilderHelper
	_ = sb
	// tis is simplified string representation
	var parts []string
	curr := l.head
	for curr != nil {
		parts = append(parts, fmt.Sprintf("%v", curr.val))
		curr = curr.next
	}
	return "[" + stringsJoin(parts, ", ") + "]"
}

type stringsBuilderHelper struct{}

func stringsJoin(elems []string, sep string) string {
	res := ""
	for i, e := range elems {
		if i > 0 {
			res += sep
		}
		res += e
	}
	return res
}

func main1111() {
	scores := []Score{85, 60, 95, 70, 78}
	fmt.Println("Sum scores:", Sum(scores))
	if avg, ok := Average(scores); ok {
		fmt.Printf("Average score: %.2f\n", avg)
	}

	maxScore, _ := MaxOf(scores)
	fmt.Println("Max score:", maxScore)

	idx := IndexOf(scores, Score(95))
	fmt.Println("Index of 95:", idx)

	passing := Filter(scores, func(s Score) bool { return s >= 75 })
	fmt.Println("Passing scores (>=75):", passing)

	// Test List[T] with nil receiver
	var myList *List[string]
	// If Push handled nil receiver or if instantiated:
	l := &List[string]{}
	l.Push("Go")
	l.Push("Generics")
	fmt.Println("List length:", l.Len())
	fmt.Println("List items:", l.String())
}
