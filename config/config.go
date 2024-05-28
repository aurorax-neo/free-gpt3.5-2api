package config

import (
	"free-gpt3.5-2api/AccAuthPool"
	"free-gpt3.5-2api/common"
	"github.com/donnie4w/go-logger/logger"
	"github.com/joho/godotenv"
	"os"
	"strings"
)

var (
	LogLevel       string
	LogPath        string
	LogFile        string
	Bind           string
	Port           string
	Proxy          []string
	AccessTokens   []string
	AUTHORIZATIONS []string
	BaseUrl        string
)

func init() {
	_ = godotenv.Load()
	// LOG_LEVEL
	LogLevel = os.Getenv("LOG_LEVEL")
	if LogLevel == "" {
		LogLevel = "INFO"
	}
	switch LogLevel {
	case "DEBUG":
		_ = logger.SetLevel(logger.LEVEL_DEBUG)
	case "INFO":
		_ = logger.SetLevel(logger.LEVEL_INFO)
	case "WARN":
		_ = logger.SetLevel(logger.LEVEL_WARN)
	case "ERROR":
		_ = logger.SetLevel(logger.LEVEL_ERROR)
	case "FATAL":
		_ = logger.SetLevel(logger.LEVEL_FATAL)
	default:
		_ = logger.SetLevel(logger.LEVEL_INFO)
	}

	// LOG_PATH
	LogPath = os.Getenv("LOG_PATH")

	// LOG_FILE
	LogFile = os.Getenv("LOG_FILE")
	_, _ = logger.SetRollingDaily(LogPath, LogFile)
	// Bind
	Bind = os.Getenv("BIND")
	if Bind == "" {
		Bind = "0.0.0.0"
	}
	// PORT
	Port = os.Getenv("PORT")
	if Port == "" {
		Port = "3040"
	}
	// PROXY
	proxy := os.Getenv("PROXY")
	if proxy != "" {
		Proxy = strings.Split(proxy, ",")
	}
	// ACCESS_TOKEN
	accessTokens := os.Getenv("ACCESS_TOKENS")
	if accessTokens == "" {
		AccessTokens = []string{}
	} else {
		AccessTokens = common.SplitAndAddPre("Bearer ", accessTokens, ",")
	}
	AccAuthPool.GetAccAuthPoolInstance().AppendAccAuths(AccessTokens)
	// AUTH_TOKEN
	authorizations := os.Getenv("AUTHORIZATIONS")
	if authorizations == "" {
		AUTHORIZATIONS = []string{}
	} else {
		//以,分割 AUTH_TOKEN 并且为每个AUTH_TOKEN前面加上Bearer
		AUTHORIZATIONS = common.SplitAndAddPre("Bearer ", authorizations, ",")
	}
	// BASE_URL
	BaseUrl = os.Getenv("BASE_URL")
	if BaseUrl == "" {
		BaseUrl = "https://chatgpt.com"
	} else {
		BaseUrl = strings.TrimRight(BaseUrl, "/")
	}
}
