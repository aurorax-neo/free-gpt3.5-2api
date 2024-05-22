package v1

import (
	"fmt"
	"free-gpt3.5-2api/FreeChatPool"
	"github.com/aurorax-neo/go-logger"
	"github.com/gin-gonic/gin"
)

type TokensResp struct {
	Count int `json:"count"`
}

func Tokens(c *gin.Context) {
	resp := &TokensResp{
		Count: FreeChatPool.GetFreeChatPoolInstance().GetSize(),
	}
	logger.Logger.Info(fmt.Sprint("FreeChatPool Tokens: ", resp.Count))
	c.JSON(200, resp)
}
