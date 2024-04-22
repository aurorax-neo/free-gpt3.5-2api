package v1Chat

import (
	"free-gpt3.5-2api/common"
	"free-gpt3.5-2api/service/v1Chat/reqmodel"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Completions(c *gin.Context) {
	// 从请求中获取参数
	apiReq := &reqmodel.ApiReq{}
	err := c.BindJSON(apiReq)
	if err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, "Invalid parameter", nil)
		return
	}
	Gpt35(c, apiReq)
}
