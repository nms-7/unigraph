package main

import (
	"context"
	"encoding/json"
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
    server := http.Server{
       Addr: "127.0.0.1:8080",
    }
    http.HandleFunc("/pools/", getPools)
    server.ListenAndServe()
}

func getPools(w http.ResponseWriter, r *http.Request) {
    /* decode */
	decoder := json.NewDecoder(r.Body)
	var asset Asset
	err := decoder.Decode(&asset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
    poolQ, err := queryPools(asset.Id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError) // replace with other error?
        return
    }
    encoder := json.NewEncoder(w)
    err = encoder.Encode(poolQ)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    return
}

// graphql logic + request(s)
func queryPools(assetId string) (poolQ PoolQ, err error) {
	variables := map[string]interface{}{
		"assetId": graphql.ID(assetId),
	}
	err = client.Query(context.Background(), &poolQ, variables)
    return
}
