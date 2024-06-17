package main

import (
	"fmt"
	"free-gpt3.5-2api/api"
	"free-gpt3.5-2api/config"
	"github.com/donnie4w/go-logger/logger"
	"net/http"
)

func main() {
	host := config.Bind
	if config.Bind == "0.0.0.0" {
		host = "127.0.0.1"
	}
	// HTTP 服务
	http.HandleFunc("/", api.HandlerGinEngine)
	// 提示服务启动
	logger.Info(fmt.Sprint("Server started on http://", host, ":", config.Port))
	// 启动 HTTP 服务器
	err := http.ListenAndServe(fmt.Sprint(config.Bind, ":", config.Port), nil)
	if err != nil {
		panic(err)
	}
}
