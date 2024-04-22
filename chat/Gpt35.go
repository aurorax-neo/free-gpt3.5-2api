package chat

import (
	"encoding/json"
	"fmt"
	"free-gpt3.5-2api/ProxyPool"
	"free-gpt3.5-2api/RequestClient"
	"free-gpt3.5-2api/common"
	"free-gpt3.5-2api/config"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/aurorax-neo/go-logger"
	fhttp "github.com/bogdanfinn/fhttp"
	"github.com/bogdanfinn/tls-client/profiles"
	"github.com/google/uuid"
	"io"
)

const BaseUrl = "https://chat.openai.com"
const ApiUrl = BaseUrl + "/backend-anon/conversation"
const SessionUrl = BaseUrl + "/backend-anon/sentinel/chat-requirements"

type Gpt35 struct {
	RequestClient RequestClient.RequestClient
	MaxUseCount   int
	ExpiresIn     int64
	Session       *session
	Ua            string
	Language      string
	IsUpdating    bool
}

type session struct {
	OaiDeviceId string           `json:"-"`
	Persona     string           `json:"persona"`
	Arkose      arkose           `json:"arkose"`
	Turnstile   turnstile        `json:"turnstile"`
	ProofWork   common.ProofWork `json:"proofofwork"`
	Token       string           `json:"token"`
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
		MaxUseCount: 1,
		ExpiresIn:   common.GetTimestampSecond(config.AuthED),
		Session:     &session{},
		Ua:          browser.Firefox(),
		Language:    common.RandomLanguage(),
		IsUpdating:  false,
	}
	// 获取代理池
	ProxyPoolInstance := ProxyPool.GetProxyPoolInstance()
	// 如果代理池中有代理数大于 1 则使用 各自requestClient
	if len(ProxyPoolInstance.Proxies) > 1 {
		instance.RequestClient = RequestClient.NewTlsClient(300, profiles.Firefox_102)
	} else {
		instance.RequestClient = RequestClient.GetInstance()
	}
	err := instance.RequestClient.SetProxy(ProxyPoolInstance.GetProxy().String())
	if err != nil {
		logger.Logger.Error(fmt.Sprint("SetProxy Error: ", err))
	}
	// 获取新的 session
	err = instance.getNewSession()
	if err != nil {
		return &Gpt35{
			MaxUseCount: 0,
			ExpiresIn:   0,
			IsUpdating:  true,
		}
	}
	return instance
}

func (G *Gpt35) getNewSession() error {
	// 生成新的设备 ID
	G.Session.OaiDeviceId = uuid.New().String()
	// 创建请求
	request, err := G.NewRequest("POST", SessionUrl, nil)
	if err != nil {
		return err
	}
	// 设置请求头
	request.Header.Set("oai-device-id", G.Session.OaiDeviceId)
	// 发送 POST 请求
	response, err := G.RequestClient.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("StatusCode: %d", response.StatusCode)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	if err := json.NewDecoder(response.Body).Decode(&G.Session); err != nil {
		return err
	}
	if G.Session.ProofWork.Required {
		G.Session.ProofWork.Ospt = common.CalcProofToken(G.Session.ProofWork.Seed, G.Session.ProofWork.Difficulty, request.Header.Get("User-Agent"))
	}
	return nil
}

func (G *Gpt35) NewRequest(method, url string, body io.Reader) (*fhttp.Request, error) {
	request, err := G.RequestClient.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("origin", common.GetOrigin(BaseUrl))
	request.Header.Set("referer", common.GetOrigin(BaseUrl))
	request.Header.Set("accept", "*/*")
	request.Header.Set("cache-control", "no-cache")
	request.Header.Set("content-type", "application/json")
	request.Header.Set("pragma", "no-cache")
	request.Header.Set("sec-ch-ua-mobile", "?0")
	request.Header.Set("sec-fetch-dest", "empty")
	request.Header.Set("sec-fetch-mode", "cors")
	request.Header.Set("sec-fetch-site", "same-origin")
	request.Header.Set("oai-language", "en-US")
	request.Header.Set("accept-language", G.Language)
	request.Header.Set("User-Agent", G.Ua)
	return request, nil
}
