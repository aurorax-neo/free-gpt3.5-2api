package chat

import (
	"encoding/json"
	"free-gpt3.5-2api/v1/chat/reqmodel"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Completions(c *gin.Context) {
	// 从请求中获取参数
	apiReq := &reqmodel.ApiReq{}
	err := json.NewDecoder(c.Request.Body).Decode(&apiReq)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"detail": "Invalid parameter",
		})
		return
	}
	gpt35(c, apiReq)
}