package main

import (
	"fmt"
	"free-gpt3.5-2api/config"
	"free-gpt3.5-2api/router"
	"github.com/aurorax-neo/go-logger"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize HTTP server
	gin.SetMode(gin.ReleaseMode)
	server := gin.New()
	server.Use(gin.Recovery())
	// 设置路由
	router.SetRouter(server)
	// 提示服务启动
	host := config.Bind
	if config.Bind == "0.0.0.0" {
		host = "127.0.0.1"
	}
	logger.Logger.Info(fmt.Sprint("Server started on http://", host, ":", config.Port))
	// 启动 HTTP 服务器
	_ = server.Run(fmt.Sprint(config.Bind, ":", config.Port))
}
