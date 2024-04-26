package ProxyPool

import (
	"fmt"
	"free-gpt3.5-2api/common"
	"free-gpt3.5-2api/config"
	"github.com/aurorax-neo/go-logger"
	fhttp "github.com/bogdanfinn/fhttp"
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
			cookies, _ := getCookies(px)
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
	cookies, _ := getCookies("")
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

func getCookies(proxy string) ([]*fhttp.Cookie, error) {
	_ = proxy
	return nil, nil
}
