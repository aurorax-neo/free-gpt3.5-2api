package router

import (
	"fmt"
	"free-gpt3.5-2api/common"
	"free-gpt3.5-2api/config"
	"github.com/aurorax-neo/go-logger"
	"github.com/gin-gonic/gin"
)

// Ping 测试接口
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// V1Cors 跨域中间件
func V1Cors(c *gin.Context) {
	// 允许跨域
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	// 如果是OPTIONS请求，直接返回
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}
	c.Next()
}

// V1Request 请求中间件
func V1Request(c *gin.Context) {
	// 打印请求摘要 方法 url ip - user-agent 格式化输出
	infoStr := fmt.Sprint(" -> ", c.Request.Method, " ", c.Request.URL.String(), " - ", c.ClientIP(), " - ", c.Request.Header.Get("User-Agent"))
	logger.Logger.Info(infoStr)
	c.Next()
}

func inArray(user string, list []string) bool {
	// 如果 list 为空，直接返回 true
	if len(list) == 0 {
		return true
	}
	for _, v := range list {
		if v == user {
			return true
		}
	}
	return false
}

// V1Auth 验证v1 api 的token
func V1Auth(c *gin.Context) {
	authToken := c.Request.Header.Get("Authorization")
	if authToken == "" && len(config.AUTHORIZATIONS) > 0 {
		common.ErrorResponse(c, 401, "You didn't provide an API key. You need to provide your API key in an Authorization header using Bearer auth (i.e. Authorization: Bearer YOUR_KEY)", nil)
		return
	}
	// 判断 authToken 是否在 config.CONFIG.AUTHORIZATIONS 列表
	if !inArray(authToken, config.AUTHORIZATIONS) {
		common.ErrorResponse(c, 401, "Incorrect API key provided: sk-4yNZz***************************************6mjw.", nil)
		return
	}
	c.Next()
}

// V1Response 响应中间件
func V1Response(c *gin.Context) {
	c.Next()
	// 打印响应摘要 方法 url 状态码
	infoStr := fmt.Sprint(" <- ", c.Request.Method, " ", c.Request.URL.String(), " - ", c.Writer.Status())
	logger.Logger.Info(infoStr)
}
