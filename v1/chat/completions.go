package chat

import (
	v1 "free-gpt3.5-2api/v1"
	"free-gpt3.5-2api/v1/chat/reqmodel"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Completions(c *gin.Context) {
	// 从请求中获取参数
	apiReq := &reqmodel.ApiReq{}
	err := c.BindJSON(apiReq)
	if err != nil {
		v1.ErrorResponse(c, http.StatusBadRequest, "Invalid parameter", nil)
		return
	}
	gpt35(c, apiReq)
}
