package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/shurcooL/graphql"
	"net/http"
	"strconv"
	"time"
)

const CHUNK_SIZE = 86400 * 10 // 10 days
const UNISWAP_V3_ALT_GRAPH_EPOCH = 1620086400

type AssetVolumeRequest struct {
	Id             string `json:"id"`
	StartTimestamp int    `json:"start_timestamp"`
	EndTimestamp   int    `json:"end_timestamp"`
}

type AssetVolumeResponse struct {
	Volume float64 `json:"volume"`
}

type TokenDataDaysQ struct {
	TokenDayDatas []TokenDayData `graphql:"tokenDayDatas(where: {token: $assetId date_gte: $startTimestamp date_lt: $endTimestamp })"`
}

// note: hour data from subgraph always returns 0
type TokenDayData struct {
	VolumeUSD string `graphql:"volumeUSD"`
}

func getVolume(c *gin.Context) {
	var request AssetVolumeRequest
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	response, err := queryVolume(request.Id, request.StartTimestamp, request.EndTimestamp)

	if err != nil {
		c.IndentedJSON(http.StatusBadGateway, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, response)
}

func queryVolume(assetId string, startTimestamp, endTimestamp int) (response AssetVolumeResponse, err error) {
	// optional tokendaydata untracked

	startTimestamp = max(startTimestamp, UNISWAP_V3_ALT_GRAPH_EPOCH)
	endTimestamp = min(endTimestamp, int(time.Now().Unix()))

	timeDiff := endTimestamp - startTimestamp
	if timeDiff <= 0 {
		err = errors.New("Invalid input: end_timestamp must be greater than start_timestamp")
		return
	}

	numChunks := timeDiff/CHUNK_SIZE + 1
	chunkChan := make(chan struct {
		TokenDataDaysQ
		error
	}, numChunks)
	var end int
	for start := startTimestamp; start < endTimestamp; start = end {
		end = min(start+CHUNK_SIZE, endTimestamp)
		variables := map[string]interface{}{
			"assetId":        graphql.ID(assetId),
			"startTimestamp": graphql.Int(start),
			"endTimestamp":   graphql.Int(end),
		}
		go queryVolumeGraph(variables, chunkChan)
	}
	for i := 0; i < numChunks; i++ {
		pair := <-chunkChan
		if pair.error != nil {
			err = pair.error
			return
		}
		chunkDaysQ := pair.TokenDataDaysQ
		for _, day := range chunkDaysQ.TokenDayDatas {
			dayVolume, parseErr := strconv.ParseFloat(day.VolumeUSD, 64)
			if parseErr != nil {
				err = parseErr
				return
			}
			response.Volume += dayVolume
		}
	}
	return
}

func queryVolumeGraph(variables map[string]interface{}, chunkChan chan struct {
	TokenDataDaysQ
	error
}) {
	var chunkDaysQ TokenDataDaysQ
	err := client.Query(context.Background(), &chunkDaysQ, variables)
	chunkChan <- struct {
		TokenDataDaysQ
		error
	}{chunkDaysQ, err}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
