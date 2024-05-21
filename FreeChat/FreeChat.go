package FreeChat

import (
	"encoding/json"
	"fmt"
	ProofWork2 "free-gpt3.5-2api/ProofWork"
	"free-gpt3.5-2api/ProxyPool"
	"free-gpt3.5-2api/RequestClient"
	"free-gpt3.5-2api/common"
	"free-gpt3.5-2api/config"
	"free-gpt3.5-2api/constant"
	"github.com/aurorax-neo/go-logger"
	fhttp "github.com/bogdanfinn/fhttp"
	"github.com/google/uuid"
	"io"
	"strings"
)

var (
	BaseUrl          = config.BaseUrl
	FreeAuthUrl      = BaseUrl + "/backend-anon/sentinel/chat-requirements"
	FreeAuthChatUrl  = BaseUrl + "/backend-anon/conversation"
	AccAuthChatUrl   = BaseUrl + "/backend-api/conversation"
	OfficialBaseURLS = []string{"https://chat.openai.com", "https://chatgpt.com"}
)

// NewFreeAuthType 定义一个枚举类型
type NewFreeAuthType int

const (
	NewFreeAuthNormal  NewFreeAuthType = 0 //正常获取
	NewFreeAuthRefresh NewFreeAuthType = 1 // 刷新获取
)

type FreeChat struct {
	RequestClient RequestClient.RequestClient
	Proxy         *ProxyPool.Proxy
	MaxUseCount   int
	ExpiresAt     int64
	FreeAuth      *freeAuth
	AccAuth       string
	ChatUrl       string
	Ua            string
	Cookies       []*fhttp.Cookie
}

type freeAuth struct {
	OaiDeviceId string               `json:"-"`
	Persona     string               `json:"persona"`
	Arkose      arkose               `json:"arkose"`
	Turnstile   turnstile            `json:"turnstile"`
	ProofWork   ProofWork2.ProofWork `json:"proofofwork"`
	Token       string               `json:"token"`
	ForceLogin  bool                 `json:"force_login"`
}

type arkose struct {
	Required bool   `json:"required"`
	Dx       string `json:"dx"`
}

type turnstile struct {
	Required bool `json:"required"`
}

// NewFreeChat 创建 FreeChat 实例 0 无论网络是否被标记限制都获取 1 在网络未标记时才能获取
func NewFreeChat(newType NewFreeAuthType, maxUseCount int, expiresAt int64, accAuth string) *FreeChat {
	// 创建 FreeChat 实例
	freeChat := &FreeChat{
		MaxUseCount: maxUseCount,
		ExpiresAt:   expiresAt,
		FreeAuth:    &freeAuth{},
		Ua:          constant.Ua,
	}
	// ChatUrl
	if strings.HasPrefix(accAuth, "Bearer eyJhbGciOiJSUzI1NiI") {
		freeChat.ChatUrl = AccAuthChatUrl
		freeChat.AccAuth = accAuth
	} else {
		freeChat.ChatUrl = FreeAuthChatUrl
	}
	// 获取请求客户端
	err := freeChat.newRequestClient()
	if err != nil {
		logger.Logger.Debug(err.Error())
		return nil
	}
	// 获取并设置代理
	err = freeChat.getProxy(newType)
	if err != nil {
		logger.Logger.Debug(err.Error())
		return nil
	}
	// 获取cookies
	if common.IsStrInArray(BaseUrl, OfficialBaseURLS) {
		err = freeChat.getCookies()
		if err != nil {
			logger.Logger.Debug(err.Error())
			return nil
		}
	}
	// 获取 FreeAuth
	err = freeChat.newFreeAuth(newType)
	if err != nil {
		logger.Logger.Debug(err.Error())
		return nil
	}
	return freeChat
}

func (FG *FreeChat) NewRequest(method, url string, body io.Reader) (*fhttp.Request, error) {
	request, err := RequestClient.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if FG.AccAuth != "" {
		request.Header.Set("Authorization", FG.AccAuth)
	}
	request.Header.Set("accept", "*/*")
	request.Header.Set("accept-language", "zh-CN,zh;q=0.9,zh-Hans;q=0.8,en;q=0.7")
	for _, cookie := range FG.Cookies {
		request.AddCookie(cookie)
	}
	request.Header.Set("oai-language", "en-US")
	request.Header.Set("origin", common.GetOrigin(url))
	request.Header.Set("referer", common.GetOrigin(url))
	request.Header.Set("sec-ch-ua", `"Microsoft Edge";v="123", "Not:A-Brand";v="8", "Chromium";v="123"`)
	request.Header.Set("sec-ch-ua-mobile", "?0")
	request.Header.Set("sec-ch-ua-platform", `"Windows"`)
	request.Header.Set("sec-fetch-dest", "empty")
	request.Header.Set("sec-fetch-mode", "cors")
	request.Header.Set("sec-fetch-site", "same-origin")
	request.Header.Set("user-agent", FG.Ua)
	request.Header.Set("Connection", "close")
	return request, nil
}

func (FG *FreeChat) newRequestClient() error {
	// 请求客户端
	FG.RequestClient = RequestClient.NewTlsClient(300, constant.ClientProfile)
	if FG.RequestClient == nil {
		errStr := fmt.Sprint("RequestClient is nil")
		logger.Logger.Debug(errStr)
		return fmt.Errorf(errStr)
	}
	return nil
}

func (FG *FreeChat) getProxy(newFreeAuthType NewFreeAuthType) error {
	// 获取代理池
	ProxyPoolInstance := ProxyPool.GetProxyPoolInstance()
	// 获取代理
	FG.Proxy = ProxyPoolInstance.GetProxy()
	// 判断代理是否可用
	if FG.Proxy.CanUseAt > common.GetTimestampSecond(0) && newFreeAuthType == NewFreeAuthRefresh {
		errStr := fmt.Sprint(FG.Proxy.Link, ": Proxy restricted, Reuse at ", FG.Proxy.CanUseAt)
		return fmt.Errorf(errStr)
	}
	// 补全cookies
	FG.Cookies = append(FG.Cookies, FG.Proxy.Cookies...)
	// 设置代理
	err := FG.RequestClient.SetProxy(FG.Proxy.Link.String())
	if err != nil {
		errStr := fmt.Sprint("SetProxy Error: ", err)
		logger.Logger.Debug(errStr)
	}
	return nil
}

func (FG *FreeChat) getCookies() error {
	// 获取cookies
	request, err := FG.NewRequest("GET", fmt.Sprint(BaseUrl, "/?oai-dm=1"), nil)
	if err != nil {
		return err
	}
	// 设置请求头
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	if FG.AccAuth != "" {
		request.Header.Set("Authorization", FG.AccAuth)
	}
	// 发送 GET 请求
	response, err := FG.RequestClient.Do(request)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	if response.StatusCode != 200 {
		return fmt.Errorf("StatusCode: %d", response.StatusCode)
	}
	// 获取cookies
	cookies := response.Cookies()
	for i, cookie := range cookies {
		if cookie.Name == "oai-did" {
			FG.FreeAuth.OaiDeviceId = cookie.Value
			cookies = append(cookies[:i], cookies[i+1:]...)
		}
		if cookie.Name == "__Secure-next-auth.callback-url" {
			cookie.Value = BaseUrl
		}
	}
	// 设置cookies
	FG.Cookies = append(FG.Cookies, cookies...)
	return nil
}

func (FG *FreeChat) newFreeAuth(newFreeAuthType NewFreeAuthType) error {
	// 生成新的设备 ID
	if FG.FreeAuth.OaiDeviceId == "" {
		FG.FreeAuth.OaiDeviceId = uuid.New().String()
	}
	// 创建请求
	request, err := FG.NewRequest("POST", FreeAuthUrl, nil)
	if err != nil {
		return err
	}
	// 设置请求头
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("oai-device-id", FG.FreeAuth.OaiDeviceId)
	// 发送 POST 请求
	response, err := FG.RequestClient.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		logger.Logger.Debug(fmt.Sprint("newFreeAuth: StatusCode: ", response.StatusCode))
		if (response.StatusCode == 429 || response.StatusCode == 403) && newFreeAuthType == NewFreeAuthRefresh {
			FG.Proxy.CanUseAt = common.GetTimestampSecond(300)
			logger.Logger.Debug(fmt.Sprint("newFreeAuth: Proxy(", FG.Proxy.Link, ")restricted, Reuse at ", FG.Proxy.CanUseAt))
		}
		return fmt.Errorf("StatusCode: %d", response.StatusCode)
	} else if newFreeAuthType == 0 {
		// 成功后更新代理的可用时间
		FG.Proxy.CanUseAt = common.GetTimestampSecond(0)
		logger.Logger.Debug(fmt.Sprint("newFreeAuth: Proxy(", FG.Proxy.Link, ")Reuse at ", FG.Proxy.CanUseAt))
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	if err := json.NewDecoder(response.Body).Decode(&FG.FreeAuth); err != nil {
		return err
	}
	if FG.FreeAuth.ForceLogin {
		errStr := fmt.Sprint("ForceLogin: ", FG.FreeAuth.ForceLogin)
		return fmt.Errorf(errStr)
	}
	if strings.HasPrefix(FG.FreeAuth.ProofWork.Difficulty, "00003") {
		errStr := fmt.Sprint("Too Difficulty: ", FG.FreeAuth.ProofWork.Difficulty)
		return fmt.Errorf(errStr)
	}
	// ProofWork
	if FG.FreeAuth.ProofWork.Required {
		FG.FreeAuth.ProofWork.Ospt = ProofWork2.CalcProofToken(FG.FreeAuth.ProofWork.Seed, FG.FreeAuth.ProofWork.Difficulty, request.Header.Get("User-Agent"))
	}
	return nil
}

// SubFreeChatMaxUseCount 减少 FreeChat 实例的最大使用次数
func (FG *FreeChat) SubFreeChatMaxUseCount() *FreeChat {
	FG.MaxUseCount--
	return FG
}
