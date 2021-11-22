package main

import (
	"context"
    "github.com/gin-gonic/gin"
	"github.com/shurcooL/graphql"
	"net/http"
    "sync"
)

// type address string

type AssetPoolsRequest struct {
	Id string `json:"id"`
    Verbose bool `json:"verbose"`
}

type PoolQ struct {
    Pools []Pool `graphql:"pools(where: { token0: $assetId } orderyBy: volumeUSD orderDirection: desc)" json:"pools"`
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

    poolQ, err := queryPools(request.Id)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, err.Error())
        // replace with other error?
        return
    }

    c.IndentedJSON(http.StatusOK, poolQ)
}

// graphql logic + request(s)
func queryPools(assetId string) (pools PoolQ, err error) {
    var poolQ0 struct {
        Pools []Pool `graphql:"pools(where: { token0: $assetId } orderyBy: volumeUSD orderDirection: desc)"`
    }

    var poolQ1 struct {
        Pools []Pool `graphql:"pools(where: { token1: $assetId } orderyBy: volumeUSD orderDirection: desc)"`
    }
	variables := map[string]interface{}{
		"assetId": graphql.ID(assetId),
	}
    var wg sync.WaitGroup
    wg.Add(2)
    go queryGraph(&poolQ0, variables, &wg)
    go queryGraph(&poolQ1, variables, &wg)
    wg.Wait()
    // chan for appending?
    pools.Pools = append(poolQ0.Pools, poolQ1.Pools...)
    // wait for both queries to complete
    // benchmark this methodology for proof of implementation
    return
}

// add chan to pass errors?
func queryGraph(query interface{}, variables map[string]interface{}, wg *sync.WaitGroup) {
    client.Query(context.Background(), query, variables)
    wg.Done()
}
