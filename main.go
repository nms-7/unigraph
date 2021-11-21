package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/shurcooL/graphql"
	"net/http"
)

type address string

type Asset struct {
	id address `json:"id"`
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
	fmt.Println(client)
}

func main() {
	/*
	   server := http.Server{
	       Addr: "127.0.0.1:8080",
	   }
	   http.HandleFunc("/pools/", getPools)
	   server.ListenAndServe()
	*/
	queryPools()
}

func getPools(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var asset Asset
	err := decoder.Decode(&asset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// graphql logic + request(s)
func queryPools() {
	var poolQ PoolQ
	variables := map[string]interface{}{
		// temp constant
		"assetId": graphql.ID("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"),
	}
	fmt.Println(client)
	err := client.Query(context.Background(), &poolQ, variables)
	if err != nil {
		fmt.Println("error mofo:", err)
	}
	fmt.Println(poolQ)
	// encode and write response
}
