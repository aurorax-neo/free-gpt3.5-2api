package ProxyPool

import (
	"fmt"
	"free-gpt3.5-2api/common"
	"free-gpt3.5-2api/config"
	"github.com/aurorax-neo/tls_client_httpi"
	"github.com/donnie4w/go-logger/logger"
	"net/url"
	"sync"
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
	Link    *url.URL
	Cookies tls_client_httpi.Cookies
}

func GetProxyPoolInstance() *ProxyPool {
	Once.Do(func() {
		logger.Debug(fmt.Sprint("Init ProxyPool..."))
		// 初始化 ProxyPool
		Instance = NewProxyPool(nil)
		// 遍历配置文件中的代理 添加到代理池
		for _, px := range config.Proxy {
			proxy := NewProxy(px)
			Instance.AddProxy(proxy)
		}
		logger.Debug(fmt.Sprint("Init ProxyPool Success"))
	})
	return Instance
}

func NewProxyPool(proxies []*Proxy) *ProxyPool {
	proxy := NewProxy("")
	return &ProxyPool{
		Proxies: append([]*Proxy{proxy}, proxies...),
		Index:   0,
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

func NewProxy(link string) *Proxy {
	return &Proxy{
		Link: common.ParseUrl(link),
	}
}
