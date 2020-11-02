package main

import "testing"

func BenchmarkConcurrentPortScanner(b *testing.B) {
	b.Run("Scan ports from 1 to 1024", func(b *testing.B) {
		scanPortsConcurrently("httpbin.org", 1, 1024)
	})
}
