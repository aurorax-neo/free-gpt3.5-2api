package router

import (
	v1 "free-gpt3.5-2api/service/v1"
	"free-gpt3.5-2api/service/v1Chat"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetRouter(router *gin.Engine) {
	router.GET("/", Index)
	router.GET("/ping", Ping)
	v1Router := router.Group("/v1")
	v1Router.Use(V1Cors)
	v1Router.Use(V1Request)
	v1Router.Use(V1Response)
	v1Router.Use(V1Auth)
	v1Router.GET("/tokens", v1.Tokens)
	v1Router.OPTIONS("/FreeGpt35/completions", nil)
	v1Router.POST("/FreeGpt35/completions", v1Chat.Completions)
}

func Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello,This is free-gpt3.5-2api.")
}
