package config

import (
	"free-gpt3.5-2api/common"
	"github.com/joho/godotenv"
	"net/url"
	"os"
)

type config struct {
	LogLevel   string
	Bind       string
	Port       string
	Proxy      *url.URL
	AuthTokens []string
}

var CONFIG *config

func init() {
	_ = godotenv.Load()
	CONFIG = &config{}
	// Bind
	CONFIG.Bind = os.Getenv("BIND")
	if CONFIG.Bind == "" {
		CONFIG.Bind = "127.0.0.1"
	}
	// PORT
	CONFIG.Port = os.Getenv("PORT")
	if CONFIG.Port == "" {
		CONFIG.Port = "3040"
	}
	// PROXY
	proxy := os.Getenv("PROXY")
	if proxy == "" {
		CONFIG.Proxy = nil
	} else {
		CONFIG.Proxy = common.ParseUrl(proxy)
	}
	// AUTH_TOKEN
	authTokens := os.Getenv("AUTH_TOKENS")
	if authTokens == "" {
		CONFIG.AuthTokens = []string{}
	} else {
		//以,分割 AUTH_TOKEN 并且为每个AUTH_TOKEN前面加上Bearer
		CONFIG.AuthTokens = common.SplitAndAddBearer(authTokens)
	}
}
