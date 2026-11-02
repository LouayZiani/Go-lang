// Task 7: reuse buffers with sync.Pool
package task7

import (
	"bytes"
	"strconv"
	"sync"
)

// writeIDs writes "id=<n>\n" for every id into buf
func writeIDs(buf *bytes.Buffer, ids []int) {
	var tmp [20]byte
	for _, id := range ids {
		buf.WriteString("id=")
		buf.Write(strconv.AppendInt(tmp[:0], int64(id), 10))
		buf.WriteByte('\n')
	}
}

// RenderNew alloc a 4kb buffer on every call
func RenderNew(ids []int) int {
	buf := new(bytes.Buffer)
	buf.Grow(4096)
	writeIDs(buf, ids)
	return buf.Len()
}

var bufPool = sync.Pool{
	New: func() any { b := new(bytes.Buffer); b.Grow(4096); return b },
}

// RenderPool takes a buffer from the pool and gives it back afterwards
func RenderPool(ids []int) int {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset() // the buffer still holds the prev callers data
	writeIDs(buf, ids)
	n := buf.Len()
	bufPool.Put(buf) //after this line buf belongs to the pool, cannt be touched
	return n
}
