package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/shurcooL/graphql"
	"net/http"
)

type BlockSwapsRequest struct {
	BlockNumber int64 `json:"block_number"`
}

type BlockSwapsResponse struct {
	Swaps []Swap `json:"swaps"`
}

type Swap struct {
	Id        graphql.String `json:"id"`
	Pool      `json:"pool"`
	AmountUSD string `json:"amountUSD" graphql:"amountUSD"`
}

func getSwaps(c *gin.Context) {
	var request BlockSwapsRequest
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	response, err := queryTxsForSwaps(request.BlockNumber)
	if err != nil {
		c.IndentedJSON(http.StatusBadGateway, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, response)
}

func queryTxsForSwaps(blockNumber int64) (response BlockSwapsResponse, err error) {
	var txsQ struct {
		Transactions []BlockSwapsResponse `graphql:"transactions(where : { blockNumber: $blockNumber })`
	}
	variables := map[string]interface{}{
		"blockNumber": graphql.Int(blockNumber),
	}
	err = client.Query(context.Background(), &txsQ, variables)
	if err != nil {
		return
	}
	for _, tx := range txsQ.Transactions {
		response.Swaps = append(response.Swaps, tx.Swaps...)
	}
	return
}
