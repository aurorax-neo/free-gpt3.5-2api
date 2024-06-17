package api

import (
	"free-gpt3.5-2api/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

var ginEngine *gin.Engine

func init() {
	// Initialize HTTP server
	gin.SetMode(gin.ReleaseMode)
	ginEngine = gin.Default()
	// Register routes
	ginEngine.GET("/", Index)
	ginEngine.GET("/ping", Ping)
	v1Router := ginEngine.Group("/v1")
	v1Router.Use(V1Cors)
	v1Router.Use(V1Auth)
	v1Router.GET("/accTokens", service.AccTokens)
	v1Router.OPTIONS("/chat/completions", nil)
	v1Router.POST("/chat/completions", service.Completions)
}

// HandlerGinEngine gin http handler
func HandlerGinEngine(w http.ResponseWriter, r *http.Request) {
	ginEngine.ServeHTTP(w, r)
}

// Ping 测试接口
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// Index 首页
func Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello,This is free-gpt3.5-2api.")
}
