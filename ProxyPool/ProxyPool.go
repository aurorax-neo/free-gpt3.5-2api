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
	Proxies []*url.URL
	Index   int
}

func GetProxyPoolInstance() *ProxyPool {
	Once.Do(func() {
		Instance = NewProxyPool([]*url.URL{})
		for _, px := range config.Proxy {
			Instance.AddProxy(common.ParseUrl(px))
		}
	})
	return Instance
}

func NewProxyPool(proxies []*url.URL) *ProxyPool {
	return &ProxyPool{
		Proxies: proxies,
		Index:   0,
	}
}

func (p *ProxyPool) GetProxy() *url.URL {
	if len(p.Proxies) == 0 {
		return &url.URL{}
	}
	proxy := p.Proxies[p.Index]
	p.Index = (p.Index + 1) % len(p.Proxies)
	return proxy
}

func (p *ProxyPool) AddProxy(proxy *url.URL) {
	p.Proxies = append(p.Proxies, proxy)
}
