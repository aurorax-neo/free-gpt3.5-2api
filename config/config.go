package config

import (
	"free-gpt3.5-2api/common"
	"github.com/joho/godotenv"
	"net/url"
	"os"
	"strconv"
)

type config struct {
	LogLevel       string
	Bind           string
	Port           string
	Proxy          *url.URL
	AUTHORIZATIONS []string
	PoolMaxCount   int
	AuthED         int
	AuthUseCount   int
}

var CONFIG *config

func init() {
	_ = godotenv.Load()
	CONFIG = &config{}
	// Bind
	CONFIG.Bind = os.Getenv("BIND")
	if CONFIG.Bind == "" {
		CONFIG.Bind = "0.0.0.0"
	}
	// PORT
	CONFIG.Port = os.Getenv("PORT")
	if CONFIG.Port == "" {
		CONFIG.Port = "3040"
	}
	// PROXY
	proxy := os.Getenv("PROXY")
	if proxy == "" {
		CONFIG.Proxy = &url.URL{}
	} else {
		CONFIG.Proxy = common.ParseUrl(proxy)
	}
	// AUTH_TOKEN
	authorizations := os.Getenv("AUTHORIZATIONS")
	if authorizations == "" {
		CONFIG.AUTHORIZATIONS = []string{}
	} else {
		//以,分割 AUTH_TOKEN 并且为每个AUTH_TOKEN前面加上Bearer
		CONFIG.AUTHORIZATIONS = common.SplitAndAddBearer(authorizations)
	}
	// POOL_MAX_COUNT
	poolMaxCount := os.Getenv("POOL_MAX_COUNT")
	var err error
	if poolMaxCount == "" {
		CONFIG.PoolMaxCount = 5
	} else {
		CONFIG.PoolMaxCount, err = strconv.Atoi(poolMaxCount)
		if err != nil {
			CONFIG.PoolMaxCount = 5
		}
	}
	// AUTH_ED
	authED := os.Getenv("AUTH_ED")
	if authED == "" {
		CONFIG.AuthED = 180
	} else {
		CONFIG.AuthED, err = strconv.Atoi(authED)
		if err != nil {
			CONFIG.AuthED = 180
		}
	}
	// AUTH_USE_COUNT
	authUseCount := os.Getenv("AUTH_USE_COUNT")
	if authUseCount == "" {
		CONFIG.AuthUseCount = 5
	} else {
		CONFIG.AuthUseCount, err = strconv.Atoi(authUseCount)
		if err != nil {
			CONFIG.AuthUseCount = 5
		}
	}
}
