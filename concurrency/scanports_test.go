package main

import "testing"

func BenchmarkPortScanner(b *testing.B) {
	b.Run("Scan ports from 1 to 10", func(b *testing.B) {
		scanPorts("httpbin.org", 1, 10)
	})
}
