package config

import (
	"free-gpt3.5-2api/AccAuthPool"
	"free-gpt3.5-2api/common"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"strings"
)

var (
	Bind           string
	Port           string
	Proxy          []string
	ACCESS_TOKENS  []string
	AUTHORIZATIONS []string
	BaseUrl        string
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
	// ACCESS_TOKEN
	accessTokens := os.Getenv("ACCESS_TOKENS")
	if accessTokens == "" {
		ACCESS_TOKENS = []string{}
	} else {
		ACCESS_TOKENS = common.SplitAndAddPre("Bearer ", accessTokens, ",")
	}
	AccAuthPool.GetAccAuthPoolInstance().AppendAccAuths(ACCESS_TOKENS)
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
