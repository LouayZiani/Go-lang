package main

import (
	"bytes"
	"strings"
	"testing"
)

// Original version (with conversions)
func processPayloadOriginal(data []byte) string {
	trimmed := string(bytes.TrimSpace(data))
	clean := strings.ReplaceAll(trimmed, "\r", "")
	return clean
}

// Optimized version (no conversions)
func processPayloadOptimized(data []byte) []byte {
	trimmed := bytes.TrimSpace(data)
	clean := bytes.ReplaceAll(trimmed, []byte("\r"), []byte(""))
	return clean
}

func BenchmarkProcessPayloadOriginal(b *testing.B) {
	data := []byte("   hello\rworld   ")
	for i := 0; i < b.N; i++ {
		_ = processPayloadOriginal(data)
	}
}

func BenchmarkProcessPayloadOptimized(b *testing.B) {
	data := []byte("   hello\rworld   ")
	for i := 0; i < b.N; i++ {
		_ = processPayloadOptimized(data)
	}
}
