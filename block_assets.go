package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/shurcooL/graphql"
	"net/http"
)

type BlockAssetsRequest struct {
	BlockNumber int64 `json:"block_number"`
}

type BlockAssetsResponse struct {
	Assets []Token `json:"assets"`
}

type SwapTokens struct {
	Token0 Token
	Token1 Token
}

func getBlockAssets(c *gin.Context) {
	var request BlockAssetsRequest
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	response, err := queryTxsForSwapTokens(request.BlockNumber)
	if err != nil {
		c.IndentedJSON(http.StatusBadGateway, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, response)
}

func queryTxsForSwapTokens(blockNumber int64) (response BlockAssetsResponse, err error) {
	var txsQ struct {
		Transactions []struct{ Swaps []SwapTokens } `graphql:"transactions(where : { blockNumber: $blockNumber })`
	}
	variables := map[string]interface{}{
		"blockNumber": graphql.Int(blockNumber),
	}
	err = client.Query(context.Background(), &txsQ, variables)
	if err != nil {
		return
	}
	seen := make(map[graphql.String]bool)
	filterTokens := func(token Token) {
		if !seen[token.Id] {
			response.Assets = append(response.Assets, token)
			seen[token.Id] = true
		}
	}
	for _, tx := range txsQ.Transactions {
		for _, swap := range tx.Swaps {
			filterTokens(swap.Token0)
			filterTokens(swap.Token1)
		}
	}
	return
}
