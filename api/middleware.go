package api

import (
	"free-gpt3.5-2api/common"
	"free-gpt3.5-2api/config"
	"github.com/gin-gonic/gin"
	"strings"
)

// V1Cors 跨域中间件
func V1Cors(c *gin.Context) {
	// 允许跨域
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Token, Content-Type, Accept")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	// 如果是OPTIONS请求，直接返回
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}
	c.Next()
}

// V1Auth 验证v1 api 的token
func V1Auth(c *gin.Context) {
	authToken := c.Request.Header.Get("Authorization")
	if strings.HasPrefix(authToken, "Bearer eyJhbGciOiJSUzI1NiI") {
		c.Next()
		return
	}
	if authToken == "" && len(config.AUTHORIZATIONS) > 0 {
		common.ErrorResponse(c, 401, "You didn't provide an API key. You need to provide your API key in an Token header using Bearer auth (i.e. Token: Bearer YOUR_KEY)", nil)
		return
	}
	// 判断 authToken 是否在 config.CONFIG.AUTHORIZATIONS 列表
	if !common.IsStrInArray(authToken, config.AUTHORIZATIONS) {
		common.ErrorResponse(c, 401, "Incorrect API key provided: sk-4yNZz***************************************6mjw.", nil)
		return
	}
	c.Next()
}
