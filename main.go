package main

import (
	"context"
    "github.com/gin-gonic/gin"
	"github.com/shurcooL/graphql"
	"net/http"
)

// type address string

type Asset struct {
	Id string `json:"id"`
}

type PoolQ struct {
	Pools []Pool `graphql:"pools(first: 2 where: { token0: $assetId })"`
}

type Pool struct {
	Id graphql.String
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
    var asset Asset
    if err := c.BindJSON(&asset); err != nil {
        c.IndentedJSON(http.StatusInternalServerError, err.Error())
        return
    }

    poolQ, err := queryPools(asset.Id)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, err.Error())
        // replace with other error?
        return
    }

    c.IndentedJSON(http.StatusOK, poolQ)
}

// graphql logic + request(s)
func queryPools(assetId string) (poolQ PoolQ, err error) {
	variables := map[string]interface{}{
		"assetId": graphql.ID(assetId),
	}
	err = client.Query(context.Background(), &poolQ, variables)
    return
}
