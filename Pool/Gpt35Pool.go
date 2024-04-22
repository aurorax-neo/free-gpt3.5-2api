package Pool

import (
	"fmt"
	"free-gpt3.5-2api/chat"
	"free-gpt3.5-2api/common"
	"free-gpt3.5-2api/config"
	"github.com/aurorax-neo/go-logger"
	"sync"
	"time"
)

var (
	instance *Gpt35Pool
	once     sync.Once
)

func init() {
	instance = GetGpt35PoolInstance()
}

type Gpt35Pool struct {
	Gpt35s   []*chat.Gpt35
	Index    int
	MaxCount int
	Lock     sync.Mutex
}

func GetGpt35PoolInstance() *Gpt35Pool {
	once.Do(func() {
		instance = &Gpt35Pool{
			Gpt35s:   make([]*chat.Gpt35, config.PoolMaxCount),
			Index:    0,
			MaxCount: config.PoolMaxCount,
		}
		logger.Logger.Info(fmt.Sprint("PoolMaxCount: ", config.PoolMaxCount, ", AuthExpirationDate: ", config.AuthED, ", Init Pool..."))
		// 定时刷新 Pool
		go instance.timingUpdateGpt35Pool(60)
	})
	return instance
}

func (G *Gpt35Pool) GetGpt35(retry int) *chat.Gpt35 {
	// 加锁
	G.Lock.Lock()
	defer G.Lock.Unlock()
	if G.IsLiveGpt35(G.Index) { //有缓存
		// 获取 Gpt35 实例
		gpt35 := G.Gpt35s[G.Index]
		// 可用次数减 1
		gpt35.MaxUseCount--
		// 返回 深拷贝的 Gpt35 实例
		gpt35_ := chat.Gpt35{
			RequestClient: gpt35.RequestClient,
			MaxUseCount:   gpt35.MaxUseCount,
			ExpiresIn:     gpt35.ExpiresIn,
			Session:       gpt35.Session,
			Ua:            gpt35.Ua,
			Language:      gpt35.Language,
		}
		// 更新 index 的 Gpt35 实例
		go G.updateGpt35AtIndex(G.Index)
		// 索引加 1，采用取模运算实现循环
		G.Index = (G.Index + 1) % G.MaxCount
		return &gpt35_
	} else if retry > 0 { //无缓存或者缓存无效
		// 释放锁 防止死锁
		G.Lock.Unlock()
		defer G.Lock.Lock()
		// 更新 index 的 Gpt35 实例
		G.updateGpt35AtIndex(G.Index)
		// 等待 Gpt35 实例刷新完成
		G.waitGpt35AtIndexUpdated(G.Index)
		// 保证不会死循环
		if retry == 1 {
			// 索引加 1，采用取模运算实现循环
			G.Index = (G.Index + 1) % G.MaxCount
		}
		// 递归获取 Gpt35 实例
		return G.GetGpt35(retry - 1)
	}
	return nil
}

func (G *Gpt35Pool) timingUpdateGpt35Pool(sec int) {
	ticker := time.NewTicker(time.Duration(sec) * time.Second)
	defer ticker.Stop()
	G.updateGpt35Pool()
	for {
		select {
		case <-ticker.C:
			G.updateGpt35Pool()
		}
	}
}

func (G *Gpt35Pool) IsLiveGpt35(index int) bool {
	//判断是否为空
	if G.Gpt35s[index] == nil || //空的
		G.Gpt35s[index].MaxUseCount <= 0 || //无可用次数
		G.Gpt35s[index].ExpiresIn <= common.GetTimestampSecond(0) ||
		G.Gpt35s[index].IsUpdating {
		return false
	}
	return true
}

func (G *Gpt35Pool) updateGpt35AtIndex(index int) bool {
	if index < 0 || index >= len(G.Gpt35s) {
		return false
	}
	if G.Gpt35s[index] != nil && G.Gpt35s[index].IsUpdating {
		return false
	}
	if !G.IsLiveGpt35(index) {
		// 标志 Gpt35 实例正在刷新
		if G.Gpt35s[index] != nil {
			G.Gpt35s[index].IsUpdating = true
		}
		G.Gpt35s[index] = chat.NewGpt35()
		// 标志 Gpt35 没有正在刷新
		if G.Gpt35s[index] != nil {
			G.Gpt35s[index].IsUpdating = false
		}
		return true
	}
	// 标志 Gpt35 没有正在刷新
	if G.Gpt35s[index] != nil {
		G.Gpt35s[index].IsUpdating = false
	}
	return false
}

func (G *Gpt35Pool) waitGpt35AtIndexUpdated(index int) {
	// 加锁
	G.Lock.Lock()
	defer G.Lock.Unlock()
	// 等待 Gpt35 实例刷新完成
	for {
		if G.Gpt35s[index] == nil || !G.Gpt35s[index].IsUpdating {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func (G *Gpt35Pool) updateGpt35Pool() {
	for i := 0; i < G.MaxCount; i++ {
		G.updateGpt35AtIndex(i)
	}
}
