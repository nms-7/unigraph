package main

import (
	"context"
    "github.com/gin-gonic/gin"
	"github.com/shurcooL/graphql"
	"net/http"
)

// type address string

type AssetPoolsRequest struct {
	Id string `json:"id"`
}

type AssetPoolsResponse struct {
    Pools []Pool `json:"pools"`
}


type Pool struct {
    Id graphql.String `json:"id"`
    Token0 Token `json:"token0"`
    Token1 Token `json:"token1"`
}

type Token struct {
    Id graphql.String `json:"id"`
    Symbol graphql.String `json:"symbol"`
}

var client *graphql.Client

func init() {
	client = graphql.NewClient("https://api.thegraph.com/subgraphs/name/ianlapham/uniswap-v3-alt", nil)
}

func main() {
    router := gin.Default()
    router.GET("/pools", getPools)
    router.Run("localhost:8080")
}

func getPools(c *gin.Context) {
    var request AssetPoolsRequest
    if err := c.BindJSON(&request); err != nil {
        c.IndentedJSON(http.StatusInternalServerError, err.Error())
        return
    }

    response, err := queryPools(request.Id)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, err.Error())
        // replace with other error?
        return
    }

    c.IndentedJSON(http.StatusOK, response)
}

// graphql logic + request(s)
func queryPools(assetId string) (pools AssetPoolsResponse, err error) {
    var poolQ0 struct {
        Pools []Pool `graphql:"pools(where: { token0: $assetId } orderyBy: volumeUSD orderDirection: desc)"`
    }

    var poolQ1 struct {
        Pools []Pool `graphql:"pools(where: { token1: $assetId } orderyBy: volumeUSD orderDirection: desc)"`
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
    // chan for appending?
    pools.Pools = append(poolQ0.Pools, poolQ1.Pools...)
    // sort?
    // wait for both queries to complete
    // benchmark this methodology for proof of implementation
    return
}

// add chan to pass errors?
func queryGraph(query interface{}, variables map[string]interface{}, errs chan error) {
    err := client.Query(context.Background(), query, variables)
    errs <- err
}
