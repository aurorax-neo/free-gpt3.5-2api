package FreeGpt35Pool

import (
	"fmt"
	"free-gpt3.5-2api/FreeGpt35"
	"free-gpt3.5-2api/common"
	"free-gpt3.5-2api/config"
	"free-gpt3.5-2api/queue"
	"github.com/aurorax-neo/go-logger"
	"sync"
	"time"
)

var (
	instance *FreeGpt35Pool
	once     sync.Once
)

type FreeGpt35Pool struct {
	queue    *queue.Queue
	capacity int // 队列容量
}

func newGpt35Pool(capacity int) *FreeGpt35Pool {
	return &FreeGpt35Pool{
		queue:    queue.New(),
		capacity: capacity,
	}
}

func GetGpt35PoolInstance() *FreeGpt35Pool {
	once.Do(func() {
		logger.Logger.Info(fmt.Sprint("PoolMaxCount: ", config.PoolMaxCount, ", AuthExpirationDate: ", config.AuthED, ", Init FreeGpt35Pool..."))
		// 初始化 FreeGpt35Pool
		instance = newGpt35Pool(config.PoolMaxCount)
		// 定时刷新 FreeGpt35Pool
		instance.updateGpt35Pool(time.Millisecond * 256)
	})
	return instance
}

func (G *FreeGpt35Pool) updateGpt35Pool(sleep time.Duration) {
	// 检测 FreeGpt35Pool 是否已满
	common.AsyncLoopTask(sleep, func() {
		// 判断 FreeGpt35Pool 是否已满
		if G.IsFull() {
			return
		}
		// 获取新 FreeGpt35 实例
		gpt35 := FreeGpt35.NewGpt35(1)
		// 判断 FreeGpt35 实例是否有效
		if G.isLiveGpt35(gpt35) {
			// 入队新 FreeGpt35 实例
			G.Enqueue(gpt35)
		}
	})
	// 检测并移除无效 FreeGpt35 实例
	common.AsyncLoopTask(sleep, func() {
		// 遍历队列中的所有元素
		G.queue.Traverse(func(n *queue.Node) {
			// 判断是否为无效 FreeGpt35 实例
			if !G.isLiveGpt35(n.Value.(*FreeGpt35.Gpt35)) {
				// 移除无效 FreeGpt35 实例
				G.queue.Remove(n)
			}
		})
	})
}

func (G *FreeGpt35Pool) isLiveGpt35(gpt35 *FreeGpt35.Gpt35) bool {
	//判断是否为空
	if gpt35 == nil ||
		gpt35.MaxUseCount <= 0 || //无可用次数
		gpt35.ExpiresAt <= common.GetTimestampSecond(0) {
		return false
	}
	return true
}

func (G *FreeGpt35Pool) GetGpt35(retry int) *FreeGpt35.Gpt35 {
	// 获取 FreeGpt35 实例
	n := G.queue.Peek()
	if n != nil {
		gpt35 := n.Value.(*FreeGpt35.Gpt35)
		if G.isLiveGpt35(gpt35) { //有缓存
			// 深拷贝
			gpt35_ := common.DeepCopyStruct(gpt35).(*FreeGpt35.Gpt35)
			// 减少 FreeGpt35 实例的最大使用次数
			gpt35.MaxUseCount--
			// 判断 FreeGpt35 实例是否有效 无效则移除
			if !G.isLiveGpt35(gpt35) {
				G.queue.Dequeue()
			}
			return gpt35_
		} else if retry > 0 {
			time.Sleep(time.Millisecond * 128)
			return G.GetGpt35(retry - 1)
		}
	}
	// 缓存内无可用 FreeGpt35 实例，返回新 FreeGpt35 实例
	return FreeGpt35.NewGpt35(0)
}

// GetSize 获取队列当前元素个数
func (G *FreeGpt35Pool) GetSize() int {
	return G.queue.Len()
}

// GetCapacity 获取队列容量
func (G *FreeGpt35Pool) GetCapacity() int {
	return G.capacity
}

// IsFull 检查队列是否已满
func (G *FreeGpt35Pool) IsFull() bool {
	return G.GetSize() == G.capacity
}

// Enqueue 入队
func (G *FreeGpt35Pool) Enqueue(v *FreeGpt35.Gpt35) bool {
	if G.IsFull() || v == nil {
		return false
	}
	G.queue.Enqueue(v)
	return true
}
