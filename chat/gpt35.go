package chat

import (
	"encoding/json"
	"fmt"
	"free-gpt3.5-2api/common"
	"free-gpt3.5-2api/config"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/aurorax-neo/go-logger"
	fhttp "github.com/bogdanfinn/fhttp"
	tlsClient "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"github.com/google/uuid"
	"io"
	"strings"
)

const BaseUrl = "https://chat.openai.com"
const ApiUrl = BaseUrl + "/backend-anon/conversation"
const SessionUrl = BaseUrl + "/backend-anon/sentinel/chat-requirements"

type Gpt35 struct {
	Client      tlsClient.HttpClient
	MaxUseCount int
	ExpiresIn   int64
	IsLapse     bool
	Session     *session
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
	jar := tlsClient.NewCookieJar()
	options := []tlsClient.HttpClientOption{
		tlsClient.WithTimeoutSeconds(300),
		tlsClient.WithClientProfile(profiles.Okhttp4Android13),
		tlsClient.WithNotFollowRedirects(),
		tlsClient.WithCookieJar(jar),
		tlsClient.WithProxyUrl(config.CONFIG.Proxy.String()),
	}
	client, err := tlsClient.NewHttpClient(tlsClient.NewNoopLogger(), options...)
	if err != nil {
		return nil
	}
	instance := &Gpt35{
		Client:      client,
		MaxUseCount: config.CONFIG.AuthUseCount,
		ExpiresIn:   common.GetTimestampSecond(config.CONFIG.AuthED),
		IsLapse:     false,
		Session:     &session{},
	}
	// 获取新的 session
	err = instance.getNewSession()
	if err != nil {
		logger.Logger.Error(fmt.Sprint("Failed to get authorization: ", err))
		logger.Logger.Error("If this error persists, your country may not be supported yet.")
		logger.Logger.Error("If your country was the issue, please consider using a U.S. PROXY.")
		return nil
	}
	return instance
}

func (G *Gpt35) getNewSession() error {
	// 生成新的设备 ID
	G.Session.OaiDeviceId = uuid.New().String()
	// 设置请求体
	body := strings.NewReader(`{"conversation_mode_kind":"primary_assistant"}`)
	// 创建请求
	request, err := G.NewRequest("POST", SessionUrl, body)
	if err != nil {
		return nil
	}
	// 设置请求头
	request.Header.Set("oai-device-id", G.Session.OaiDeviceId)
	// 发送 POST 请求
	response, err := G.Client.Do(request)
	if err != nil {
		return nil
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("StatusCode: %d", response.StatusCode)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	if err := json.NewDecoder(response.Body).Decode(&G.Session); err != nil {
		return nil
	}
	return nil
}

func (G *Gpt35) NewRequest(method, url string, body io.Reader) (*fhttp.Request, error) {
	request, err := fhttp.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("origin", BaseUrl)
	request.Header.Set("referer", BaseUrl)
	request.Header.Set("accept", "*/*")
	request.Header.Set("accept-language", common.RandomLanguage())
	request.Header.Set("cache-control", "no-cache")
	request.Header.Set("content-type", "application/json")
	request.Header.Set("oai-language", "en-US")
	request.Header.Set("pragma", "no-cache")
	request.Header.Set("sec-ch-ua-mobile", "?0")
	request.Header.Set("sec-fetch-dest", "empty")
	request.Header.Set("sec-fetch-mode", "cors")
	request.Header.Set("sec-fetch-site", "same-origin")
	request.Header.Set("User-Agent", browser.Random())
	return request, nil
}
