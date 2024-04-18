package router

import (
	"free-gpt3.5-2api/middleware"
	v1 "free-gpt3.5-2api/v1"
	"free-gpt3.5-2api/v1/chat"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetRouter(router *gin.Engine) {

	router.GET("/", Index)
	router.GET("/ping", middleware.Ping)
	v1Router := router.Group("/v1")
	v1Router.Use(middleware.V1Cors)
	v1Router.Use(middleware.V1Request)
	v1Router.Use(middleware.V1Response)
	v1Router.Use(middleware.V1Auth)
	v1Router.GET("/tokens", v1.Tokens)
	v1Router.OPTIONS("/chat/completions", nil)
	v1Router.POST("/chat/completions", chat.Completions)
}

func Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello,This is free-gpt3.5-2api.")
}
