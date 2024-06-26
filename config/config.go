package config

import (
	"free-gpt3.5-2api/AccessTokenPool"
	"free-gpt3.5-2api/common"
	"github.com/donnie4w/go-logger/logger"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

var (
	LogLevel       string
	LogPath        string
	LogFile        string
	Bind           string
	Port           string
	Proxy          []string
	TokensFile     string
	AUTHORIZATIONS []string
	BaseUrl        string
)

type Tokens struct {
	AccessTokens []*AccessTokenPool.AccessToken `yaml:"access_tokens,omitempty"`
}

func init() {
	_ = godotenv.Load()
	// logger Options
	LogOpts := &logger.Option{
		Level:   logger.LEVEL_INFO,
		Console: true,
	}
	// LOG_PATH
	LogPath = os.Getenv("LOG_PATH")
	if LogPath == "" {
		LogPath = "logs"
	}
	// LOG_FILE
	LogFile = os.Getenv("LOG_FILE")
	if LogFile != "" {
		LogFilename := filepath.Join(LogPath, LogFile)
		LogOpts.FileOption = &logger.FileTimeMode{
			Filename:   LogFilename,
			Maxbuckup:  10,
			IsCompress: true,
			Timemode:   logger.MODE_DAY,
		}
	}
	// LOG_LEVEL
	LogLevel = os.Getenv("LOG_LEVEL")
	switch LogLevel {
	case "DEBUG":
		LogOpts.Level = logger.LEVEL_DEBUG
	case "INFO":
		LogOpts.Level = logger.LEVEL_INFO
	case "WARN":
		LogOpts.Level = logger.LEVEL_WARN
	case "ERROR":
		LogOpts.Level = logger.LEVEL_ERROR
	case "FATAL":
		LogOpts.Level = logger.LEVEL_FATAL
	default:
		LogOpts.Level = logger.LEVEL_INFO
	}
	// set logger options
	logger.SetOption(LogOpts)

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
	// ACCESS_TOKEN_FILE
	accessTokensFile := os.Getenv("TOKENS_FILE")
	if accessTokensFile == "" {
		TokensFile = common.GetAbsPath("tokens.yml")
	} else {
		TokensFile = common.GetAbsPath(accessTokensFile)
	}
	if common.IsFileExist(TokensFile) {
		bytes, err := common.ReadFile(TokensFile)
		if err != nil {
			logger.Error("ReadFile error: ", err)
		}
		var tokens Tokens
		if err = yaml.Unmarshal(bytes, &tokens); err != nil {
			logger.Error("Unmarshal error: ", err)
		}
		for _, token := range tokens.AccessTokens {
			token.Token = "Bearer " + token.Token
		}
		AccessTokenPool.GetAccAuthPoolInstance().AppendAccessTokens(tokens.AccessTokens)
	}
	// AUTH_TOKEN
	authorizations := os.Getenv("AUTHORIZATIONS")
	if authorizations == "" {
		panic("please add AUTHORIZATIONS in environment variable or .env file  (example: AUTHORIZATIONS=authkey1,authkey2)")
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
