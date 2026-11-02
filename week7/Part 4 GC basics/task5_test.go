package task5

import (
	"runtime"
	"testing"
)

// liveHeapMB runs build, forces a GC and returns the live heap in MB. runtime.KeepAlive(x) keeps the result
// reachable until AFTERRR the measurement. W/o it the GC may free x early (it is not used later) and the leaking
// version would show smth close to 0mb
func liveHeapMB(build func() any) uint64 {
	var before, after runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&before)

	x := build()

	runtime.GC()
	runtime.ReadMemStats(&after)
	runtime.KeepAlive(x)

	if after.HeapAlloc < before.HeapAlloc {
		return 0
	}
	return (after.HeapAlloc - before.HeapAlloc) >> 20
}

const n = 100

func TestLeakBytes(t *testing.T) {
	leak := liveHeapMB(func() any { return KeepHeadersSlice(n) })
	fixed := liveHeapMB(func() any { return KeepHeadersCopy(n) })
	t.Logf("subslices: %d MB   copies: %d MB", leak, fixed)
	if leak < 90 {
		t.Errorf("expected the slice version to keep ~100 MB alive, got %d MB", leak)
	}
	if fixed > 1 {
		t.Errorf("expected the copy version to keep ~0 MB, got %d MB", fixed)
	}
}

func TestLeakStrings(t *testing.T) {
	leak := liveHeapMB(func() any { return KeepSubstrings(n) })
	fixed := liveHeapMB(func() any { return KeepClones(n) })
	t.Logf("substrings: %d MB   clones: %d MB", leak, fixed)
	if leak < 90 {
		t.Errorf("expected the substring version to keep ~100 MB alive, got %d MB", leak)
	}
	if fixed > 1 {
		t.Errorf("expected the clone version to keep ~0 MB, got %d MB", fixed)
	}
}

func TestHeaderContent(t *testing.T) {
	a, b := KeepHeadersSlice(3), KeepHeadersCopy(3)
	for i := range a {
		if string(a[i]) != string(b[i]) {
			t.Fatalf("content differs at %d", i)
		}
	}
	if string(b[2]) != "HEADER-0002" {
		t.Fatalf("got %q", b[2])
	}
}
