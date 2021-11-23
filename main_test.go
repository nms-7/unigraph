package main

import (
    "testing"
)

func BenchmarkQueryPools(b *testing.B) {
    for i := 0; i < b.N; i++ {
        queryPools("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48")
    }
}
