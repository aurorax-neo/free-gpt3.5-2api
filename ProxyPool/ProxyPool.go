package ProxyPool

import (
	"fmt"
	"free-gpt3.5-2api/RequestClient"
	"free-gpt3.5-2api/common"
	"free-gpt3.5-2api/config"
	"github.com/aurorax-neo/go-logger"
	fhttp "github.com/bogdanfinn/fhttp"
	"github.com/bogdanfinn/tls-client/profiles"
	"io"
	"net/url"
	"sync"
)

var (
	Ua = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"
)

var (
	Instance *ProxyPool
	Once     sync.Once
)

type ProxyPool struct {
	Proxies []*Proxy
	Index   int
}

type Proxy struct {
	Link     *url.URL
	CanUseAt int64
	Ua       string
	Cookies  []*fhttp.Cookie
}

func GetProxyPoolInstance() *ProxyPool {
	Once.Do(func() {
		logger.Logger.Info(fmt.Sprint("Init ProxyPool..."))
		Instance = NewProxyPool(nil)
		for _, px := range config.Proxy {
			cookies, _ := getCookies(px, Ua)
			Instance.AddProxy(&Proxy{
				Link:     common.ParseUrl(px),
				CanUseAt: common.GetTimestampSecond(0),
				Ua:       Ua,
				Cookies:  cookies,
			})

		}
		logger.Logger.Info(fmt.Sprint("Init ProxyPool Success"))
	})
	return Instance
}

func NewProxyPool(proxies []*Proxy) *ProxyPool {
	cookies, _ := getCookies("", Ua)
	return &ProxyPool{
		Proxies: append([]*Proxy{
			{
				Link:     &url.URL{},
				CanUseAt: common.GetTimestampSecond(0),
				Ua:       Ua,
				Cookies:  cookies,
			},
		}, proxies...),
		Index: 0,
	}
}

func (PP *ProxyPool) GetProxy() *Proxy {
	PP.Index = (PP.Index + 1) % len(PP.Proxies)
	// 如果配置了代理 不会使用无代理
	if PP.Index == 0 && len(PP.Proxies) > 1 {
		PP.Index = 1
	}
	// 返回代理
	return PP.Proxies[PP.Index]
}

func (PP *ProxyPool) AddProxy(proxy *Proxy) {
	PP.Proxies = append(PP.Proxies, proxy)
}

func getCookies(proxy string, ua string) ([]*fhttp.Cookie, error) {
	// 获取cookies
	request, err := RequestClient.NewRequest("GET", "https://chat.openai.com", nil)
	if err != nil {
		return nil, err
	}
	// 设置请求头
	request.Header.Set("User-Agent", ua)
	// 获取请求客户端
	client := RequestClient.NewTlsClient(60, profiles.Okhttp4Android13)
	// 设置代理
	_ = client.SetProxy(proxy)
	// 发送 GET 请求
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("StatusCode: %d", response.StatusCode)
	}
	// 获取cookies
	cookies := response.Cookies()
	for i, cookie := range cookies {
		if cookie.Name == "oai-did" {
			cookies = append(cookies[:i], cookies[i+1:]...)
		}
		if cookie.Name == "__Secure-next-auth.callback-url" {
			cookie.Value = "https://chat.openai.com"
		}
	}
	return cookies, nil
}
