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
	mutex    sync.Mutex
}

func GetGpt35PoolInstance() *Gpt35Pool {
	once.Do(func() {
		logger.Logger.Info(fmt.Sprint("Gpt35Pool init, PoolMaxCount: ", config.CONFIG.PoolMaxCount))
		gpt35PoolInstance = &Gpt35Pool{
			Gpt35s:   make([]*chat.Gpt35, config.CONFIG.PoolMaxCount),
			Index:    -1,
			MaxCount: config.CONFIG.PoolMaxCount,
		}
		gpt35PoolInstance.initGpt35Pool()
		// 启动一个 goroutine 定时刷新 Gpt35Pool
		go func() {
			// 遍历池子 g.Gpt35s
			for i := 0; i < gpt35PoolInstance.MaxCount; i++ {
				//判断是否为空
				if gpt35PoolInstance.Gpt35s[i] == nil || //空的
					gpt35PoolInstance.Gpt35s[i].MaxUseCount <= 0 || //可使用次数为0
					gpt35PoolInstance.Gpt35s[i].IsLapse || //失效
					gpt35PoolInstance.Gpt35s[i].ExpiresIn <= common.GetTimestampSecond(0) { //过期
					//重新初始化
					gpt35PoolInstance.raGpt35AtIndex(i)
				}
				if i == gpt35PoolInstance.MaxCount-1 {
					i = -1
					time.Sleep(1 * time.Second)
				}
			}
		}()
	})
	return gpt35PoolInstance
}

func (G *Gpt35Pool) initGpt35Pool() {
	G.mutex.Lock()
	defer G.mutex.Unlock()
	for i := 0; i < G.MaxCount; i++ {
		gpt35 := chat.NewGpt35()
		if gpt35 != nil {
			G.Gpt35s[i] = gpt35
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
		go G.raGpt35AtIndex(G.Index)
		return G.GetGpt35(retry)
	}
	gpt35.MaxUseCount--
	return gpt35
}

func (G *Gpt35Pool) raGpt35AtIndex(index int) {
	G.mutex.Lock()
	defer G.mutex.Unlock()
	if index < 0 || index >= len(G.Gpt35s) {
		return
	}
	G.Gpt35s[index] = chat.NewGpt35()
	logger.Logger.Info(fmt.Sprint("Gpt35Pool update Gpt35, index: ", index))
}
