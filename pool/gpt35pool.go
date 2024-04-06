package pool

import (
	"fmt"
	"free-gpt3.5-2api/chat"
	"free-gpt3.5-2api/common"
	"free-gpt3.5-2api/config"
	"github.com/aurorax-neo/go-logger"
	"sync"
)

func init() {
	GetGpt35PoolInstance()
}

var (
	gpt35PoolInstance *Gpt35Pool
	once              sync.Once
)

type Gpt35Pool struct {
	Gpt35s    []*chat.Gpt35
	Index     int
	MaxCount  int
	LiveCount int
	mutex     sync.Mutex
}

func GetGpt35PoolInstance() *Gpt35Pool {
	once.Do(func() {
		gpt35PoolInstance = &Gpt35Pool{
			Gpt35s:    make([]*chat.Gpt35, config.CONFIG.PoolMaxCount),
			Index:     -1,
			LiveCount: 0,
			MaxCount:  config.CONFIG.PoolMaxCount,
		}
		logger.Logger.Info(fmt.Sprint("PoolMaxCount: ", config.CONFIG.PoolMaxCount, ", AuthUseCount: ", config.CONFIG.AuthUseCount, ", AuthExpirationDate: ", config.CONFIG.AuthED, ", Init Gpt35Pool..."))
		// 启动一个 goroutine 刷新 Gpt35Pool
		go gpt35PoolInstance.flushGpt35Pool()
	})
	return gpt35PoolInstance
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
		go G.raGpt35AtIndex(G.Index)
		return G.GetGpt35(retry)
	}
	gpt35.MaxUseCount--
	return gpt35
}

func (G *Gpt35Pool) flushGpt35Pool() {
	for i := 0; i < G.MaxCount; i++ {
		//判断是否为空
		if G.Gpt35s[i] == nil || //空的
			G.Gpt35s[i].MaxUseCount <= 0 || //可使用次数为0
			G.Gpt35s[i].IsLapse || //失效
			G.Gpt35s[i].ExpiresIn <= common.GetTimestampSecond(0) { //过期
			//重新初始化
			G.raGpt35AtIndex(i)
		} else {
			G.LiveCount++
		}
		if i == G.MaxCount-1 {
			G.LiveCount = 0
			i = -1
		}
	}
}

func (G *Gpt35Pool) raGpt35AtIndex(index int) {
	G.mutex.Lock()
	defer G.mutex.Unlock()
	if index < 0 || index >= len(G.Gpt35s) {
		return
	}
	G.Gpt35s[index] = chat.NewGpt35()
}
