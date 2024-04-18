package pool

import (
	"fmt"
	"free-gpt3.5-2api/chat"
	"free-gpt3.5-2api/common"
	"free-gpt3.5-2api/config"
	"github.com/aurorax-neo/go-logger"
	"sync"
	"time"
)

func init() {
	GetGpt35PoolInstance()
}

var (
	gpt35PoolInstance *Gpt35Pool
	once              sync.Once
)

type Gpt35Pool struct {
	Gpt35s   []*chat.Gpt35
	Index    int
	MaxCount int
}

func GetGpt35PoolInstance() *Gpt35Pool {
	once.Do(func() {
		gpt35PoolInstance = &Gpt35Pool{
			Gpt35s:   make([]*chat.Gpt35, config.CONFIG.PoolMaxCount),
			Index:    -1,
			MaxCount: config.CONFIG.PoolMaxCount,
		}
		logger.Logger.Info(fmt.Sprint("PoolMaxCount: ", config.CONFIG.PoolMaxCount, ", AuthUseCount: ", config.CONFIG.AuthUseCount, ", AuthExpirationDate: ", config.CONFIG.AuthED, ", Init Gpt35Pool..."))
		// 定时刷新 Gpt35Pool
		go gpt35PoolInstance.timingFlushGpt35Pool(60)
	})
	return gpt35PoolInstance
}

func (G *Gpt35Pool) GetGpt35(retry int) *chat.Gpt35 {
	// 索引加 1，采用取模运算实现循环
	G.Index = (G.Index + 1) % G.MaxCount

	// 处理索引为负数的情况
	if G.Index < 0 {
		G.Index = G.MaxCount - 1
	}

	// 返回索引对应的 Gpt35 实例
	if G.IsLive(G.Index) {
		gpt35 := G.Gpt35s[G.Index]
		gpt35.MaxUseCount--
		return gpt35
	} else if retry > 0 { // 如果 Gpt35 实例为空且重试次数大于 0，则重新获取 Gpt35 实例
		G.raGpt35AtIndex(G.Index)
		retry--
		return G.GetGpt35(retry)
	}
	return nil
}

func (G *Gpt35Pool) timingFlushGpt35Pool(sec int) {
	ticker := time.NewTicker(time.Duration(sec) * time.Second)
	defer ticker.Stop()
	G._flushGpt35Pool()
	for {
		select {
		case <-ticker.C:
			G._flushGpt35Pool()
		}
	}
}

func (G *Gpt35Pool) raGpt35AtIndex(index int) {
	if index < 0 || index >= len(G.Gpt35s) {
		return
	}
	G.Gpt35s[index] = chat.NewGpt35()
}

func (G *Gpt35Pool) _flushGpt35Pool() {
	for i := 0; i < G.MaxCount; i++ {
		if !G.IsLive(i) { //过期
			G.raGpt35AtIndex(i)
		}
	}
}

func (G *Gpt35Pool) IsLive(index int) bool {
	//判断是否为空
	if G.Gpt35s[index] == nil || //空的
		G.Gpt35s[index].MaxUseCount <= 0 || //可使用次数为0
		G.Gpt35s[index].IsLapse || //失效
		G.Gpt35s[index].ExpiresIn <= common.GetTimestampSecond(0) {
		return false
	}
	return true
}
