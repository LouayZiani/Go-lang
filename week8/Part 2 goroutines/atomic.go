package main

import "sync/atomic"

type atomicInt64 struct{ value atomic.Int64 }

func (a *atomicInt64) Add(delta int64) { a.value.Add(delta) }
func (a *atomicInt64) Load() int64     { return a.value.Load() }
