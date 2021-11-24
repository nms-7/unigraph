package main

import (
    "fmt"
    "context"
	"github.com/gin-gonic/gin"
	"github.com/shurcooL/graphql"
	"net/http"
    "strconv"
)

type AssetVolumeRequest struct {
    Id string `json:"id"`
    StartTimestamp int `json:"start_timestamp"`
    EndTimestamp int `json:"end_timestamp"`
}

type AssetVolumeResponse struct {
    Volume float64 `json:"volume"`
}

// note: hour data from subgraph always returns 0
type TokenDayData struct {
    VolumeUSD string`graphql:"volumeUSD"`
}

func getVolume(c *gin.Context) {
    var request AssetVolumeRequest
    if err := c.BindJSON(&request); err != nil {
        c.IndentedJSON(http.StatusInternalServerError, err.Error())
        return
    }

    response, err := queryVolume(request.Id, request.StartTimestamp, request.EndTimestamp)

    if err!= nil {
        c.IndentedJSON(http.StatusBadGateway, err.Error())
        return
    }

    c.IndentedJSON(http.StatusOK, response)
}

func queryVolume(assetId string, startTimestamp, endTimestamp int) (response AssetVolumeResponse, err error) {
    // optional tokendaydata untracked
    var volumeQ struct {
        TokenDayDatas []TokenDayData `graphql:"tokenDayDatas(where: {token: $assetId date_gte: $startTimestamp date_lt: $endTimestamp })"`
    }
    variables := map[string]interface{}{
        "assetId": graphql.ID(assetId),
        "startTimestamp": graphql.Int(startTimestamp),
        "endTimestamp":graphql.Int(endTimestamp),
    }
    // concurrent chunks?
	err = client.Query(context.Background(), &volumeQ, variables)
    if err != nil {

        fmt.Println(err)
        fmt.Println(volumeQ.TokenDayDatas)
        return
    }
    // sum volumes
    for _, day := range volumeQ.TokenDayDatas {
        dayVolume, parseErr := strconv.ParseFloat(day.VolumeUSD, 64)
        if parseErr != nil {
            err = parseErr
            return
        }
        response.Volume += dayVolume 
        fmt.Println(dayVolume, response.Volume)
    }
    return
}
