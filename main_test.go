package main

import (
	"github.com/shurcooL/graphql"
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

func TestQueryBlockAssetsForDupes(t *testing.T) {
	randomBlocks := []int64{12774522, 12739081, 13046554, 13162473}
	for _, blockNumber := range randomBlocks {
		response, err := queryTxsForSwapTokens(blockNumber)
		if err != nil {
			t.Error(err)
		}
		seen := make(map[graphql.String]bool)
		for _, asset := range response.Assets {
			if seen[asset.Id] {
				t.Error("Duplicate token:", asset.Symbol, "/", asset.Id)
			} else {
				seen[asset.Id] = true
			}
		}
	}
}
