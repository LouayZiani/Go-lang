package main

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

var sinkS string

func JoinIDsFast(ids []int) string {
	var sb strings.Builder
	sb.Grow(len(ids) * 8)

	var tmp [20]byte

	for i, id := range ids {
		if i > 0 {
			sb.WriteByte(',')
		}

		sb.Write(strconv.AppendInt(tmp[:0], int64(id), 10))
	}

	return sb.String()
}

func TestJoinIDsFastAllocs(t *testing.T) {
	ids := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	allocs := testing.AllocsPerRun(1000, func() {
		sinkS = JoinIDsFast(ids)
	})

	fmt.Println("JoinIDsFast allocations:", allocs)

	if allocs > 1 {
		t.Errorf("got %v allocations, want at most 1", allocs)
	}
}

// package main

// import (
// 	"fmt"
// 	"testing"
// )

// var sinkS string

// func JoinIDsSlow(ids []int) string {
// 	s := ""

// 	for i, id := range ids {
// 		if i > 0 {
// 			s += ","
// 		}

// 		s += fmt.Sprint(id)
// 	}

// 	return s
// }

// func TestJoinIDsSlowAllocs(t *testing.T) {
// 	ids := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

// 	allocs := testing.AllocsPerRun(1000, func() {
// 		sinkS = JoinIDsSlow(ids)
// 	})

// 	fmt.Println("JoinIDsSlow allocations:", allocs)
// }
