package main

import "math/rand"

type Tree struct {
	Left, Right *Tree
	Value       int
}

func insert(t *Tree, v int) *Tree {
	if t == nil {
		return &Tree{Value: v}
	}
	if v < t.Value {
		t.Left = insert(t.Left, v)
	} else {
		t.Right = insert(t.Right, v)
	}
	return t
}

func New(k int) *Tree {
	var t *Tree
	for _, v := range rand.Perm(10) {
		t = insert(t, (v+1)*k)
	}
	return t
}

func walk(t *Tree, ch chan<- int, done <-chan struct{}) bool {
	if t == nil {
		return true
	}
	if !walk(t.Left, ch, done) {
		return false
	}
	select {
	case ch <- t.Value:
	case <-done:
		return false
	}
	return walk(t.Right, ch, done)
}

func Walk(t *Tree, ch chan<- int, done <-chan struct{}) {
	walk(t, ch, done)
	close(ch)
}

func Same(t1, t2 *Tree, stop bool) bool {
	var done chan struct{}
	if stop {
		done = make(chan struct{})
		defer close(done)
	}
	c1, c2 := make(chan int), make(chan int)
	go Walk(t1, c1, done)
	go Walk(t2, c2, done)
	for {
		v1, ok1 := <-c1
		v2, ok2 := <-c2
		if v1 != v2 || ok1 != ok2 {
			return false
		}
		if !ok1 {
			return true
		}
	}
}
