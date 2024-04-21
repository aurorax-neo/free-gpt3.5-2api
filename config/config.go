package config

import (
	"free-gpt3.5-2api/common"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"strings"
)

var (
	LogLevel       string
	Bind           string
	Port           string
	Proxy          []string
	AUTHORIZATIONS []string
	PoolMaxCount   int
	AuthED         int
)

func init() {
	_ = godotenv.Load()
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
	// AUTH_TOKEN
	authorizations := os.Getenv("AUTHORIZATIONS")
	if authorizations == "" {
		AUTHORIZATIONS = []string{}
	} else {
		//以,分割 AUTH_TOKEN 并且为每个AUTH_TOKEN前面加上Bearer
		AUTHORIZATIONS = common.SplitAndAddBearer(authorizations)
	}
	// POOL_MAX_COUNT
	poolMaxCount := os.Getenv("POOL_MAX_COUNT")
	var err error
	if poolMaxCount == "" {
		PoolMaxCount = 64
	} else {
		PoolMaxCount, err = strconv.Atoi(poolMaxCount)
		if err != nil {
			PoolMaxCount = 64
		}
	}
	// AUTH_ED
	authED := os.Getenv("AUTH_ED")
	if authED == "" {
		AuthED = 600
	} else {
		AuthED, err = strconv.Atoi(authED)
		if err != nil {
			AuthED = 600
		}
	}
}
