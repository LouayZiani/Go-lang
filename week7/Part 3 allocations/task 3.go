// remove the allocations
// The goal is to compare two functions that produce exactly the same result, but use memory differently.
package main

import (
	"fmt"
	"strconv"
	"strings"
)

func JoinIDsSlow(ids []int) string {
	s := "" //we created empty str that will gradually contain our result

	for i, id := range ids { //iterate thru every num in the slice
		if i > 0 {
			s += "," //we add a comma before every num except the first
		}

		s += fmt.Sprint(id) //convert int to str and append it to s
		// --> this is where its inefficient cuz every time we concat strings, Go may need to create a new string
		// and copy the prev contents into it
	}

	return s
}

func JoinIDsFast1(ids []int) string {
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

// output of the ran test:
// === RUN   TestJoinIDsSlowAllocs
// JoinIDsSlow allocations: 19
// --- PASS: TestJoinIDsSlowAllocs (0.00s)
// PASS
// ok      task3   1.086s

// output of fast id:
// === RUN   TestJoinIDsFastAllocs
// JoinIDsFast allocations: 1
// --- PASS: TestJoinIDsFastAllocs (0.00s)
// PASS
// ok      task3   1.070s

// for slow way: 19 allocations per run for 10 IDs
// for fast way: 1 allocation for 10 IDs
