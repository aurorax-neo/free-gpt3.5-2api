package chat

import (
	"crypto/tls"
	"fmt"
	"free-gpt3.5-2api/config"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"net/http"
)

const BaseUrl = "https://chat.openai.com"
const ApiUrl = BaseUrl + "/backend-anon/conversation"
const SessionUrl = BaseUrl + "/backend-anon/sentinel/chat-requirements"

type Gpt35 struct {
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

func NewGpt35() *Gpt35 {
	instance := &Gpt35{
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
		SetHeader("User-Agent", browser.Random())
	// 获取新的 session
	err := instance.getNewSession()
	if err != nil {
		return nil
	}
	return instance
}

func (G *Gpt35) getNewSession() error {
	// 生成新的设备 ID
	G.Session.OaiDeviceId = uuid.New().String()
	// 发送 POST 请求
	resp, err := G.Client.R().
		SetHeader("oai-device-id", G.Session.OaiDeviceId).
		SetBody(`{"conversation_mode_kind":"primary_assistant"}`).
		SetResult(G.Session).
		Post(SessionUrl)
	if err != nil || resp.StatusCode() != 200 {
		return fmt.Errorf("system: Failed to get new session: %v", err)
	}
	return nil
}
