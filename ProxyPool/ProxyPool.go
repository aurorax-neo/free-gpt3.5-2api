package ProxyPool

import (
	"free-gpt3.5-2api/common"
	"free-gpt3.5-2api/config"
	browser "github.com/EDDYCJY/fake-useragent"
	"net/url"
	"sync"
	"time"
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
	Ua       string
	Language string
}

func GetProxyPoolInstance() *ProxyPool {
	Once.Do(func() {
		Instance = NewProxyPool(nil)
		for _, px := range config.Proxy {
			Instance.AddProxy(&Proxy{
				Link:     common.ParseUrl(px),
				CanUseAt: common.GetTimestampSecond(0),
				Ua:       browser.Safari(),
				Language: common.RandomLanguage(),
			})
		}
		// 定时刷新代理
		Instance.timingUpdateProxy(time.Duration(config.AuthED) * time.Minute)
	})
	return Instance
}

func NewProxyPool(proxies []*Proxy) *ProxyPool {
	return &ProxyPool{
		Proxies: append([]*Proxy{
			{
				Link:     &url.URL{},
				CanUseAt: common.GetTimestampSecond(0),
				Ua:       browser.Random(),
				Language: common.RandomLanguage(),
			},
		}, proxies...),
		Index: 0,
	}
}

func (PP *ProxyPool) GetProxy() *Proxy {
	// 如果没有代理则返回空代理
	if len(PP.Proxies) == 1 {
		return PP.Proxies[0]
	}
	// 获取代理
	proxy := PP.Proxies[PP.Index]
	PP.Index = (PP.Index + 1) % len(PP.Proxies)
	if PP.Index == 0 {
		PP.Index = 1
	}
	return proxy
}

func (PP *ProxyPool) AddProxy(proxy *Proxy) {
	PP.Proxies = append(PP.Proxies, proxy)
}

func (PP *ProxyPool) timingUpdateProxy(nanosecond time.Duration) {
	common.TimingTask(nanosecond, func() {
		for _, px := range PP.Proxies {
			px.Ua = browser.Safari()
			px.Language = common.RandomLanguage()
		}
	})
}
