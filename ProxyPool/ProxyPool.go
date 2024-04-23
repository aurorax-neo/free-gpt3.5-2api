package ProxyPool

import (
	"free-gpt3.5-2api/common"
	"free-gpt3.5-2api/config"
	browser "github.com/EDDYCJY/fake-useragent"
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
	Ua       string
	Language string
}

func GetProxyPoolInstance() *ProxyPool {
	Once.Do(func() {
		Instance = NewProxyPool()
		for _, px := range config.Proxy {
			proxy := &Proxy{
				Link:     common.ParseUrl(px),
				Ua:       browser.Random(),
				Language: common.RandomLanguage(),
			}
			Instance.AddProxy(proxy)
		}
	})
	return Instance
}

func NewProxyPool() *ProxyPool {
	return &ProxyPool{
		Proxies: []*Proxy{},
		Index:   0,
	}
}

func (PP *ProxyPool) GetProxy() *Proxy {
	// 如果没有代理则返回空代理
	if len(PP.Proxies) == 0 {
		return &Proxy{
			Link:     &url.URL{},
			Ua:       browser.Safari(),
			Language: common.RandomLanguage(),
		}
	}
	// 获取代理
	proxy := PP.Proxies[PP.Index]
	PP.Index = (PP.Index + 1) % len(PP.Proxies)
	return proxy
}

func (PP *ProxyPool) AddProxy(proxy *Proxy) {
	PP.Proxies = append(PP.Proxies, proxy)
}
