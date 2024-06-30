package FreeChat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"free-gpt3.5-2api/AccessTokenPool"
	"free-gpt3.5-2api/ProofWork"
	"free-gpt3.5-2api/ProxyPool"
	"free-gpt3.5-2api/common"
	"free-gpt3.5-2api/config"
	"github.com/aurorax-neo/tls_client_httpi"
	"github.com/aurorax-neo/tls_client_httpi/tls_client"
	"github.com/donnie4w/go-logger/logger"
	"github.com/google/uuid"
	"io"
	"strings"
)

var (
	BaseUrl            = config.BaseUrl
	FreeAuthUrl        = BaseUrl + "/backend-anon/sentinel/chat-requirements"
	FreeAuthChatUrl    = BaseUrl + "/backend-anon/conversation"
	AccessTokenChatUrl = BaseUrl + "/backend-api/conversation"
	OfficialBaseURLS   = []string{"https://chat.openai.com", "https://chatgpt.com"}
)

// NewFreeAuthType 定义一个枚举类型
type NewFreeAuthType int

type FreeChat struct {
	Http     tls_client_httpi.TCHI
	Proxy    *ProxyPool.Proxy
	FreeAuth *freeAuth
	AccAuth  string
	ChatUrl  string
	Ua       string
	Cookies  tls_client_httpi.Cookies
}

type freeAuth struct {
	OaiDeviceId string              `json:"-"`
	Persona     string              `json:"persona"`
	Arkose      arkose              `json:"arkose"`
	Turnstile   turnstile           `json:"turnstile"`
	ProofWork   ProofWork.ProofWork `json:"proofofwork"`
	Token       string              `json:"token"`
	ForceLogin  bool                `json:"force_login"`
}

type arkose struct {
	Required bool   `json:"required"`
	Dx       string `json:"dx"`
}

type turnstile struct {
	Required bool `json:"required"`
}

func GetFreeChat(token string, retry int) (*FreeChat, error) {
	// 判断是否为指定账号
	if strings.HasPrefix(token, "Bearer eyJhbGciOiJSUzI1NiI") {
		auth := common.GetStrPreOrSuf(token, "#", 1)
		if !common.IsStrInArray("Bearer "+auth, config.AUTHORIZATIONS) {
			return nil, fmt.Errorf("unauthorized, please add authkey in access_tokens (example: access_tokens#authkey)")
		}
		at := common.GetStrPreOrSuf(token, "#", -1)
		freeChat, err := newFreeChat(at)
		if freeChat == nil && retry > 0 {
			return GetFreeChat(token, retry-1)
		}
		return freeChat, err
	}
	// 判断是否使用 AccessTokenPool
	if strings.HasPrefix(token, AccessTokenPool.AccAuthAuthorizationPre) && !AccessTokenPool.GetAccAuthPoolInstance().IsEmpty() {
		at := AccessTokenPool.GetAccAuthPoolInstance().GetToken()
		if at == "" {
			return nil, fmt.Errorf("AccessTokenPool is Empty")
		}
		freeChat, err := newFreeChat(at)
		if freeChat == nil && retry > 0 {
			return GetFreeChat(token, retry-1)
		}
		return freeChat, err
	}
	// 返回免登 FreeChat 实例
	freeChat, err := newFreeChat(token)
	if freeChat == nil && retry > 0 {
		return GetFreeChat(token, retry-1)
	}
	return freeChat, err
}

// newFreeChat 创建 FreeChat 实例
func newFreeChat(token string) (*FreeChat, error) {
	// 创建 FreeChat 实例
	freeChat := &FreeChat{
		FreeAuth: &freeAuth{},
		Ua:       common.GetUa(),
		ChatUrl:  FreeAuthChatUrl,
	}
	// ChatUrl
	if strings.HasPrefix(token, "Bearer eyJhbGciOiJSUzI1NiI") {
		freeChat.ChatUrl = AccessTokenChatUrl
		freeChat.AccAuth = token
	}
	// 获取请求客户端
	err := freeChat.newRequestClient()
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	// 获取并设置代理
	err = freeChat.getProxy()
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	// 获取cookies
	if common.IsStrInArray(BaseUrl, OfficialBaseURLS) {
		err = freeChat.getCookies()
		if err != nil {
			logger.Debug(err.Error())
			return nil, err
		}
	}
	// 获取 FreeAuth
	err = freeChat.newFreeAuth()
	if err != nil {
		logger.Debug(err.Error())
		return nil, err
	}
	return freeChat, nil
}

func (f *FreeChat) GetHC(url string) (tls_client_httpi.Headers, tls_client_httpi.Cookies) {
	headers := tls_client_httpi.Headers{}
	headers.Set(strings.ToLower("accept"), "*/*")
	headers.Set(strings.ToLower("accept-language"), "zh-CN,zh;q=0.9,zh-Hans;q=0.8,en;q=0.7")
	headers.Set(strings.ToLower("oai-language"), "en-US")
	headers.Set(strings.ToLower("origin"), common.GetOrigin(url))
	headers.Set(strings.ToLower("referer"), common.GetOrigin(url))
	headers.Set(strings.ToLower("sec-ch-ua"), `"Microsoft Edge";v="123", "Not:A-Brand";v="8", "Chromium";v="123"`)
	headers.Set(strings.ToLower("sec-ch-ua-mobile"), "?0")
	headers.Set(strings.ToLower("sec-ch-ua-platform"), `"Windows"`)
	headers.Set(strings.ToLower("sec-fetch-dest"), "empty")
	headers.Set(strings.ToLower("sec-fetch-mode"), "cors")
	headers.Set(strings.ToLower("sec-fetch-site"), "same-origin")
	headers.Set(strings.ToLower("user-agent"), f.Ua)
	if f.AccAuth != "" {
		headers.Set(strings.ToLower("Authorization"), f.AccAuth)
	}
	return headers, f.Cookies
}

func (f *FreeChat) newRequestClient() error {
	// 请求客户端
	f.Http = tls_client.NewClient(tls_client.NewClientOptions(300, common.GetClientProfile()))
	if f.Http == nil {
		errStr := fmt.Sprint("Http is nil")
		logger.Debug(errStr)
		return fmt.Errorf(errStr)
	}
	return nil
}

func (f *FreeChat) getProxy() error {
	// 获取代理池
	ProxyPoolInstance := ProxyPool.GetProxyPoolInstance()
	// 获取代理
	f.Proxy = ProxyPoolInstance.GetProxy()
	// 补全cookies
	f.Cookies = append(f.Cookies, f.Proxy.Cookies...)
	// 设置代理
	err := f.Http.SetProxy(f.Proxy.Link.String())
	if err != nil {
		errStr := fmt.Sprint("SetProxy Error: ", err)
		logger.Debug(errStr)
	}
	return nil
}

func (f *FreeChat) getCookies() error {
	// 获取请求头和cookies
	headers, cookies := f.GetHC(BaseUrl)
	// 设置请求头
	headers.Set(strings.ToLower("Accept"), "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	// 发送 GET 请求 获取cookies
	response, err := f.Http.Request(tls_client_httpi.GET, fmt.Sprint(BaseUrl, "/?oai-dm=1"), headers, cookies, nil)
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
	cks := response.Cookies()
	for i, cookie := range cks {
		if strings.ToLower(cookie.Name) == strings.ToLower("oai-did") {
			f.FreeAuth.OaiDeviceId = cookie.Value
			cookies = append(cks[:i], cks[i+1:]...)
		}
		if strings.ToLower(cookie.Name) == strings.ToLower("__Secure-next-auth.callback-url") {
			cookie.Value = BaseUrl
		}
	}
	// 设置cookies
	f.Cookies = append(f.Cookies, cks...)
	return nil
}

func (f *FreeChat) newFreeAuth() error {
	// 生成新的设备 ID
	if f.FreeAuth.OaiDeviceId == "" {
		f.FreeAuth.OaiDeviceId = uuid.New().String()
	}
	// 请求体
	body := bytes.NewBufferString(`{"p":"gAAAAACWzI0MTIsIlRodSBNYXkgMjMgMjAyNCAxNjozNjoyNyBHTVQrMDgwMCAoR01UKzA4OjAwKSIsNDI5NDcwNTE1MiwwLCJNb3ppbGxhLzUuMCAoV2luZG93cyBOVCAxMC4wOyBXaW42NDsgeDY0KSBBcHBsZVdlYktpdC81MzcuMzYgKEtIVE1MLCBsaWtlIEdlY2tvKSBDaHJvbWUvMTIyLjAuMC4wIFNhZmFyaS81MzcuMzYiLCJodHRwczovL2Nkbi5vYWlzdGF0aWMuY29tL19uZXh0L3N0YXRpYy9jaHVua3Mvd2VicGFjay01YzQ4NDI4ZTBlZTgxMTBlLmpzP2RwbD00ODExZmQxYzk0YjU1MGM4ZjAzZmNjODYzZWU2YzFhOTk5NDBlZmM1IiwiZHBsPTQ4MTFmZDFjOTRiNTUwYzhmMDNmY2M4NjNlZTZjMWE5OTk0MGVmYzUiLCJ6aC1DTiIsInpoLUNOLHpoLHpoLUhhbnMsZW4iLDIzNiwiZGV2aWNlTWVtb3J54oiSOCIsIl9yZWFjdExpc3RlbmluZzh6MmcweHF4M2Z4IiwiX19SRUFDVF9JTlRMX0NPTlRFWFRfXyIsNzIzLjM5OTk5OTk5ODUwOTld"}`)
	headers, cookies := f.GetHC(FreeAuthUrl)
	// 设置请求头
	headers.Set(strings.ToLower("Content-Type"), "application/json")
	headers.Set(strings.ToLower("oai-device-id"), f.FreeAuth.OaiDeviceId)
	// 发送 POST 请求
	response, err := f.Http.Request(tls_client_httpi.POST, FreeAuthUrl, headers, cookies, body)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		logger.Debug(fmt.Sprint("newFreeAuth: StatusCode: ", response.StatusCode))
		return fmt.Errorf("StatusCode: %d", response.StatusCode)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	if err := json.NewDecoder(response.Body).Decode(&f.FreeAuth); err != nil {
		return err
	}
	if f.FreeAuth.ForceLogin {
		common.SubUpdateThreshold()
		errStr := fmt.Sprint("ForceLogin: ", f.FreeAuth.ForceLogin)
		return fmt.Errorf(errStr)
	}
	if strings.HasPrefix(f.FreeAuth.ProofWork.Difficulty, "00003") {
		errStr := fmt.Sprint("Too Difficulty: ", f.FreeAuth.ProofWork.Difficulty)
		return fmt.Errorf(errStr)
	}
	// ProofWork
	if f.FreeAuth.ProofWork.Required {
		f.FreeAuth.ProofWork.Ospt = ProofWork.CalcProofToken(f.FreeAuth.ProofWork.Seed, f.FreeAuth.ProofWork.Difficulty, headers.Get(strings.ToLower("User-Agent")))
		if f.FreeAuth.ProofWork.Ospt == "" {
			errStr := fmt.Sprint("ProofWork Failed")
			return fmt.Errorf(errStr)
		}
	}
	return nil
}
