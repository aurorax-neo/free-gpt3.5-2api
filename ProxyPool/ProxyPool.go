package ProxyPool

import (
	"free-gpt3.5-2api/common"
	"free-gpt3.5-2api/config"
	"net/url"
	"sync"
)

func init() {
	Instance = GetProxyPoolInstance()
}

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
}

func GetProxyPoolInstance() *ProxyPool {
	Once.Do(func() {
		Instance = NewProxyPool(nil)
		for _, px := range config.Proxy {
			Instance.AddProxy(&Proxy{
				Link:     common.ParseUrl(px),
				CanUseAt: common.GetTimestampSecond(0),
			})
		}
	})
	return Instance
}

func NewProxyPool(proxies []*Proxy) *ProxyPool {
	return &ProxyPool{
		Proxies: append([]*Proxy{
			{
				Link:     &url.URL{},
				CanUseAt: common.GetTimestampSecond(0),
			},
		}, proxies...),
		Index: 0,
	}
}

func (PP *ProxyPool) GetProxy() *Proxy {
	// 获取代理
	proxy := PP.Proxies[PP.Index]
	PP.Index = (PP.Index + 1) % len(PP.Proxies)
	// 如果配置了代理 不会使用无代理
	if PP.Index == 0 && len(PP.Proxies) > 1 {
		PP.Index = 1
	}
	return proxy
}

func (PP *ProxyPool) AddProxy(proxy *Proxy) {
	PP.Proxies = append(PP.Proxies, proxy)
}
