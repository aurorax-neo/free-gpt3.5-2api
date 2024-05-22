package v1

import (
	"fmt"
	"free-gpt3.5-2api/AccAuthPool"
	"github.com/aurorax-neo/go-logger"
	"github.com/gin-gonic/gin"
)

type AccTokensResp struct {
	Count       int `json:"count"`
	CanUseCount int `json:"canUseCount"`
}

func AccTokens(c *gin.Context) {
	resp := &AccTokensResp{
		Count:       AccAuthPool.GetAccAuthPoolInstance().Size(),
		CanUseCount: AccAuthPool.GetAccAuthPoolInstance().CanUseSize(),
	}
	logger.Logger.Info(fmt.Sprint("AccAuthPool Tokens: ", resp.Count))
	c.JSON(200, resp)
}
