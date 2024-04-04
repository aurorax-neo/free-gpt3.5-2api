package pool

import (
	"fmt"
	"free-gpt3.5-2api/chat"
	"free-gpt3.5-2api/config"
	"github.com/aurorax-neo/go-logger"
	"sync"
	"time"
)

const refreshInterval = 60

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
	mutex    sync.Mutex
}

func GetGpt35PoolInstance() *Gpt35Pool {
	once.Do(func() {
		logger.Logger.Info(fmt.Sprint("Gpt35Pool init, PoolMaxCount: ", config.CONFIG.PoolMaxCount))
		gpt35PoolInstance = &Gpt35Pool{
			Gpt35s:   make([]*chat.Gpt35, 0),
			Index:    -1,
			MaxCount: config.CONFIG.PoolMaxCount,
		}
		// 启动一个 goroutine 定时刷新 Gpt35Pool
		go func() {
			for {
				gpt35PoolInstance.flushGpt35Pool()
				time.Sleep(refreshInterval * time.Second)
			}
		}()
	})
	return gpt35PoolInstance
}

func (G *Gpt35Pool) flushGpt35Pool() {
	G.mutex.Lock()
	defer G.mutex.Unlock()
	for i := 0; i < G.MaxCount; i++ {
		gpt35 := chat.NewGpt35()
		if gpt35 != nil {
			G.Gpt35s = append(G.Gpt35s, gpt35)
			logger.Logger.Info(fmt.Sprint("Gpt35Pool flush Gpt35 success, index: ", i))
			continue
		}
		logger.Logger.Error(fmt.Sprint("Gpt35Pool flush Gpt35 fail, index: ", i))
		i--
		time.Sleep(1 * time.Second)
	}
}

func (G *Gpt35Pool) GetGpt35(retry int) *chat.Gpt35 {
	G.mutex.Lock()
	defer G.mutex.Unlock()
	// 索引加 1
	G.Index++
	// 如果索引等于最大数量，则重置为 0
	if G.Index >= G.MaxCount {
		G.Index = 0
	}
	// 返回索引对应的 Gpt35 实例
	gpt35 := G.Gpt35s[G.Index]
	// 如果 Gpt35 实例为空且重试次数大于 0，则重新获取 Gpt35 实例
	if gpt35 == nil && retry > 0 {
		retry--
		go G.RAGpt35AtIndex(G.Index)
		return G.GetGpt35(retry)
	}
	return gpt35
}

func (G *Gpt35Pool) RAGpt35AtIndex(index int) {
	G.mutex.Lock()
	defer G.mutex.Unlock()
	if index < 0 || index >= len(G.Gpt35s) {
		return
	}
	gpt35 := chat.NewGpt35()
	if gpt35 != nil {
		G.Gpt35s[index] = gpt35
		return
	}
	G.Gpt35s[index] = nil
}
