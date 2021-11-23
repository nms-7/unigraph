package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/shurcooL/graphql"
)

type Token struct {
	Id     graphql.String `json:"id"`
	Symbol graphql.String `json:"symbol"`
}

var client *graphql.Client

func init() {
	client = graphql.NewClient("https://api.thegraph.com/subgraphs/name/ianlapham/uniswap-v3-alt", nil)
}

func main() {
	router := gin.Default()
	router.GET("/asset/pools", getPools)
	router.Run("localhost:8080")
}

func queryGraph(query interface{}, variables map[string]interface{}, errs chan error) {
	err := client.Query(context.Background(), query, variables)
	errs <- err
}
