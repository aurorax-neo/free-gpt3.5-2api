package v1

import (
	"fmt"
	"free-gpt3.5-2api/Pool"
	"github.com/aurorax-neo/go-logger"
	"github.com/gin-gonic/gin"
)

type TokensResp struct {
	Count int `json:"count"`
}

func Tokens(c *gin.Context) {
	resp := &TokensResp{
		Count: 0,
	}
	instance := Pool.GetGpt35PoolInstance()
	for i := 0; i < instance.MaxCount; i++ {
		if instance.IsLiveGpt35(i) {
			resp.Count++
		}
	}
	logger.Logger.Info(fmt.Sprint("Pool Tokens: ", resp.Count))
	c.JSON(200, resp)
}
