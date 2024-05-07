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

func newFreeGpt35Pool(capacity int) *FreeGpt35Pool {
	return &FreeGpt35Pool{
		queue:    queue.New(),
		capacity: capacity,
	}
}

func GetFreeGpt35PoolInstance() *FreeGpt35Pool {
	once.Do(func() {
		logger.Logger.Info(fmt.Sprint("Init FreeGpt35Pool..."))
		// 初始化 FreeGpt35Pool
		instance = newFreeGpt35Pool(config.PoolMaxCount)
		// 定时刷新 FreeGpt35Pool
		instance.refreshFreeGpt35Pool(time.Millisecond * 256)
		//
		logger.Logger.Info(fmt.Sprint("Init FreeGpt35Pool Success", ", PoolMaxCount: ", config.PoolMaxCount, ", AuthExpirationDate: ", config.AuthED))
	})
	return instance
}

func (G *FreeGpt35Pool) refreshFreeGpt35Pool(sleep time.Duration) {
	// 检测 FreeGpt35Pool 是否已满
	common.AsyncLoopTask(sleep, func() {
		// 判断 FreeGpt35Pool 是否已满
		if G.IsFull() {
			return
		}
		// 获取新 FreeGpt35 实例
		gpt35 := FreeGpt35.NewFreeGpt35(FreeGpt35.NewFreeAuthRefresh, 1, common.GetTimestampSecond(config.AuthED))
		// 判断 FreeGpt35 实例是否有效
		if G.isLiveGpt35(gpt35) {
			// 入队新 FreeGpt35 实例
			G.AddFreeGpt35(gpt35)
		}
	})
	// 检测并移除无效 FreeGpt35 实例
	common.AsyncLoopTask(sleep, func() {
		// 遍历队列中的所有元素
		G.queue.Traverse(func(n *queue.Node) {
			// 判断是否为无效 FreeGpt35 实例
			if !G.isLiveGpt35(n.Value.(*FreeGpt35.FreeGpt35)) {
				// 移除无效 FreeGpt35 实例
				G.queue.Remove(n)
			}
		})
	})
}

func (G *FreeGpt35Pool) isLiveGpt35(gpt35 *FreeGpt35.FreeGpt35) bool {
	//判断是否为空
	if gpt35 == nil ||
		gpt35.MaxUseCount <= 0 || //无可用次数
		gpt35.ExpiresAt <= common.GetTimestampSecond(0) {
		return false
	}
	return true
}

func (G *FreeGpt35Pool) GetFreeGpt35(retry int) *FreeGpt35.FreeGpt35 {
	// 获取 FreeGpt35 实例
	n := G.queue.Peek()
	if n != nil {
		gpt35 := n.Value.(*FreeGpt35.FreeGpt35)
		if G.isLiveGpt35(gpt35) { //有缓存
			// 深拷贝
			gpt35_ := common.DeepCopyStruct(gpt35).(*FreeGpt35.FreeGpt35)
			// 减少 FreeGpt35 实例的最大使用次数
			gpt35.MaxUseCount--
			// 判断 FreeGpt35 实例是否有效 无效则移除
			if !G.isLiveGpt35(gpt35) {
				G.queue.Dequeue()
			}
			return gpt35_
		} else if retry > 0 {
			time.Sleep(time.Millisecond * 128)
			return G.GetFreeGpt35(retry - 1)
		}
	}
	// 缓存内无可用 FreeGpt35 实例，返回新 FreeGpt35 实例
	return FreeGpt35.NewFreeGpt35(FreeGpt35.NewFreeAuthNormal, 1, common.GetTimestampSecond(config.AuthED))
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

// AddFreeGpt35 入队
func (G *FreeGpt35Pool) AddFreeGpt35(v *FreeGpt35.FreeGpt35) bool {
	if G.IsFull() || v == nil {
		return false
	}
	G.queue.Enqueue(v)
	return true
}
