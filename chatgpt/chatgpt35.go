package chatgpt

import (
	"crypto/tls"
	"fmt"
	"free-gpt3.5-2api/config"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/aurorax-neo/go-logger"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"net/http"
	"sync"
	"time"
)

const BaseUrl = "https://chat.openai.com"
const ApiUrl = BaseUrl + "/backend-api/conversation"
const refreshInterval = 60 // Interval to refresh token in ms
const errorWait = 120      // Wait time in ms after an error

func init() {
	// 启动一个 goroutine 定时刷新 session
	go func() {
		for {
			err := GetGpt35Instance().getNewSession()
			if err != nil {
				logger.Logger.Error(fmt.Sprint("refreshing session failed, retrying in ", errorWait, " second..."))
				logger.Logger.Error("if this error persists, your country may not be supported yet.")
				logger.Logger.Error("if your country was the issue, please consider using a U.S. VPN.")
				time.Sleep(errorWait * time.Second)
				continue
			}
			logger.Logger.Info(fmt.Sprint("refreshed session successfully, next refresh in ", refreshInterval, " second..."))
			time.Sleep(refreshInterval * time.Second)
		}
	}()
}

var (
	instance *gpt35
	once     sync.Once
)

type gpt35 struct {
	Client  *resty.Client
	Session *session
}

type session struct {
	OaiDeviceId string    `json:"-"`
	Persona     string    `json:"persona"`
	Arkose      arkose    `json:"arkose"`
	Turnstile   turnstile `json:"turnstile"`
	Token       string    `json:"token"`
}

type arkose struct {
	Required bool   `json:"required"`
	Dx       string `json:"dx"`
}

type turnstile struct {
	Required bool `json:"required"`
}

func GetGpt35Instance() *gpt35 {
	once.Do(func() {
		instance = &gpt35{
			Client: resty.NewWithClient(&http.Client{
				Transport: &http.Transport{
					// 禁用长连接
					DisableKeepAlives: true,
					// 配置TLS设置，跳过证书验证
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
					//配置代理
					Proxy: http.ProxyURL(config.CONFIG.Proxy),
				},
			}),
			Session: &session{},
		}
		// 设置请求头
		instance.Client.
			SetHeader("origin", BaseUrl).
			SetHeader("referer", BaseUrl).
			SetHeader("accept", "*/*").
			SetHeader("accept-language", "en-US,en;q=0.9").
			SetHeader("cache-control", "no-cache").
			SetHeader("content-type", "application/json").
			SetHeader("oai-language", "en-US").
			SetHeader("pragma", "no-cache").
			SetHeader("sec-ch-ua", `"Google Chrome";v="123", "Not:A-Brand";v="8", "Chromium";v="123"`).
			SetHeader("sec-ch-ua-mobile", "?0").
			SetHeader("sec-ch-ua-platform", "Windows").
			SetHeader("sec-fetch-dest", "empty").
			SetHeader("sec-fetch-mode", "cors").
			SetHeader("sec-fetch-site", "same-origin").
			SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	})
	return instance
}

func (G *gpt35) getNewSession() error {
	UA := browser.Random()
	G.Client.SetHeader("User-Agent", UA)
	// 生成新的设备 ID
	G.Session.OaiDeviceId = uuid.New().String()
	// 发送 POST 请求
	resp, err := G.Client.R().
		SetHeader("oai-device-id", G.Session.OaiDeviceId).
		SetBody(`{"conversation_mode_kind":"primary_assistant"}`).
		SetResult(G.Session).
		Post(BaseUrl + "/backend-anon/sentinel/chat-requirements")
	if err != nil || resp.StatusCode() != 200 {
		logger.Logger.Error(fmt.Sprintf("system: Failed to get new session: %v", err))
		return fmt.Errorf("system: Failed to get new session ID: %v", err)
	}
	return nil
}
