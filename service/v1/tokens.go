package v1

import (
	"fmt"
	"free-gpt3.5-2api/FreeGpt35Pool"
	"github.com/aurorax-neo/go-logger"
	"github.com/gin-gonic/gin"
)

type TokensResp struct {
	Count int `json:"count"`
}

func Tokens(c *gin.Context) {
	resp := &TokensResp{
		Count: FreeGpt35Pool.GetGpt35PoolInstance().GetSize(),
	}
	logger.Logger.Info(fmt.Sprint("FreeGpt35Pool Tokens: ", resp.Count))
	c.JSON(200, resp)
}
