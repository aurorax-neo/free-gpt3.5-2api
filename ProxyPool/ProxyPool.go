package ProxyPool

import (
	"fmt"
	"free-gpt3.5-2api/common"
	"free-gpt3.5-2api/config"
	"free-gpt3.5-2api/constant"
	"github.com/aurorax-neo/go-logger"
	fhttp "github.com/bogdanfinn/fhttp"
	"net/url"
	"sync"
	"time"
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
		// 初始化 ProxyPool
		Instance = NewProxyPool(nil)
		// 遍历配置文件中的代理 添加到代理池
		for _, px := range config.Proxy {
			proxy := NewProxy(px, common.GetTimestampSecond(0), constant.Ua)
			_ = proxy.getCookies()
			Instance.AddProxy(proxy)
		}
		//定时刷新代理cookies
		common.AsyncLoopTask(1*time.Minute, func() {
			for _, proxy := range Instance.Proxies {
				_ = proxy.getCookies()
			}
		})
		logger.Logger.Info(fmt.Sprint("Init ProxyPool Success"))
	})
	return Instance
}

func NewProxyPool(proxies []*Proxy) *ProxyPool {
	proxy := NewProxy("", common.GetTimestampSecond(0), constant.Ua)
	_ = proxy.getCookies()
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

func NewProxy(link string, cannotUseTime int64, ua string) *Proxy {
	return &Proxy{
		Link:     common.ParseUrl(link),
		CanUseAt: cannotUseTime,
		Ua:       ua,
	}
}

func (P *Proxy) getCookies() error {
	return nil
}
