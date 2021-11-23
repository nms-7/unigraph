package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shurcooL/graphql"
	"net/http"
)

type AssetPoolsRequest struct {
	Id string `json:"id"`
}

type AssetPoolsResponse struct {
	Pools []Pool `json:"pools"`
}

type Pool struct {
	Id     graphql.String `json:"id"`
	Token0 Token          `json:"token0"`
	Token1 Token          `json:"token1"`
}

func getPools(c *gin.Context) {
	var request AssetPoolsRequest
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	response, err := queryPools(request.Id)
	if err != nil {
		c.IndentedJSON(http.StatusBadGateway, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, response)
}

func queryPools(assetId string) (pools AssetPoolsResponse, err error) {
	var poolQ0 struct {
		Pools []Pool `graphql:"pools(where: { token0: $assetId })"`
	}

	var poolQ1 struct {
		Pools []Pool `graphql:"pools(where: { token1: $assetId })"`
	}
	variables := map[string]interface{}{
		"assetId": graphql.ID(assetId),
	}
	errs := make(chan error, 2)
	go queryGraph(&poolQ0, variables, errs)
	go queryGraph(&poolQ1, variables, errs)
	for i := 0; i < 2; i++ {
		err = <-errs
		if err != nil {
			return
		}
	}
	pools.Pools = append(poolQ0.Pools, poolQ1.Pools...)
	return
}


