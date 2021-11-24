package main

import (
	"testing"
)

var USDC = "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"

func BenchmarkQueryPools(b *testing.B) {
	for i := 0; i < b.N; i++ {
		queryPools(USDC)
	}
}

// for benchmarking chunk span
func BenchmarkQueryVolume(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for i := 0; i < b.N; i++ {
			queryVolume(USDC, 1600000000, 1635724800)
		}
	}
}
