package service

import (
	"fmt"
	"free-gpt3.5-2api/AccessTokenPool"
	"github.com/donnie4w/go-logger/logger"
	"github.com/gin-gonic/gin"
)

type AccTokensResp struct {
	Count       int `json:"count"`
	CanUseCount int `json:"canUseCount"`
}

func AccTokens(c *gin.Context) {
	resp := &AccTokensResp{
		Count:       AccessTokenPool.GetAccAuthPoolInstance().Size(),
		CanUseCount: AccessTokenPool.GetAccAuthPoolInstance().CanUseSize(),
	}
	logger.Info(fmt.Sprint("AccessTokenPool Tokens: ", resp.Count))
	c.JSON(200, resp)
}
