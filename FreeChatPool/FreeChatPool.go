package FreeChatPool

import (
	"fmt"
	"free-gpt3.5-2api/FreeChat"
	"free-gpt3.5-2api/common"
	"free-gpt3.5-2api/config"
	"free-gpt3.5-2api/queue"
	"github.com/aurorax-neo/go-logger"
	"strings"
	"sync"
	"time"
)

var (
	instance *FreeChatPool
	once     sync.Once
)

type FreeChatPool struct {
	queue    *queue.Queue
	capacity int // 队列容量
}

func newFreeChatPool(capacity int) *FreeChatPool {
	return &FreeChatPool{
		queue:    queue.New(),
		capacity: capacity,
	}
}

func GetFreeChatPoolInstance() *FreeChatPool {
	once.Do(func() {
		logger.Logger.Info(fmt.Sprint("Init FreeChatPool..."))
		// 初始化 FreeChatPool
		instance = newFreeChatPool(config.PoolMaxCount)
		// 定时刷新 FreeChatPool
		instance.refreshFreeChatPool(time.Millisecond * 128)
		//
		logger.Logger.Info(fmt.Sprint("Init FreeChatPool Success", ", PoolMaxCount: ", config.PoolMaxCount, ", AuthExpirationDate: ", config.AuthED))
	})
	return instance
}

func (G *FreeChatPool) refreshFreeChatPool(sleep time.Duration) {
	// 检测 FreeChatPool 是否已满
	common.AsyncLoopTask(sleep, func() {
		// 判断 FreeChatPool 是否已满
		if G.IsFull() {
			return
		}
		// 获取新 FreeChat 实例
		freeChat := FreeChat.NewFreeGpt35(FreeChat.NewFreeAuthRefresh, 1, common.GetTimestampSecond(config.AuthED), "")
		// 判断 FreeChat 实例是否有效
		if G.isLiveFreeChat(freeChat) {
			// 入队新 FreeChat 实例
			G.AddFreeGpt35(freeChat)
		}
	})
	// 检测并移除无效 FreeChat 实例
	common.AsyncLoopTask(sleep, func() {
		// 遍历队列中的所有元素
		G.queue.Traverse(func(n *queue.Node) {
			// 判断是否为无效 FreeChat 实例
			if !G.isLiveFreeChat(n.Value.(*FreeChat.FreeChat)) {
				// 移除无效 FreeChat 实例
				G.queue.Remove(n)
			}
		})
	})
}

func (G *FreeChatPool) isLiveFreeChat(freeChat *FreeChat.FreeChat) bool {
	//判断是否为空
	if freeChat == nil ||
		freeChat.MaxUseCount <= 0 || //无可用次数
		freeChat.ExpiresAt <= common.GetTimestampSecond(0) {
		return false
	}
	return true
}

func (G *FreeChatPool) GetFreeChat(accAuth string, retry int) *FreeChat.FreeChat {
	if strings.HasPrefix(accAuth, "Bearer eyJhbGciOiJSUzI1NiI") {
		freeChat := FreeChat.NewFreeGpt35(FreeChat.NewFreeAuthNormal, 1, common.GetTimestampSecond(config.AuthED), accAuth)
		if freeChat == nil && retry > 0 {
			return G.GetFreeChat(accAuth, retry-1)
		}
		return freeChat
	}
	// 获取 FreeChat 实例
	n := G.queue.Peek()
	if n != nil {
		freeChat := n.Value.(*FreeChat.FreeChat)
		// 判断 FreeChat 实例是否有效
		if G.isLiveFreeChat(freeChat) {
			// 减少 FreeChat 实例可用次数
			freeChat.SubFreeGpt35MaxUseCount()
			// 判断 FreeChat 实例是否有效 无效则移除
			if !G.isLiveFreeChat(freeChat) {
				G.queue.Dequeue()
			}
			return freeChat
		} else if retry > 0 {
			time.Sleep(time.Millisecond * 128)
			return G.GetFreeChat(accAuth, retry-1)
		}
	}
	// 缓存内无可用 FreeChat 实例，返回新 FreeChat 实例
	return FreeChat.NewFreeGpt35(FreeChat.NewFreeAuthNormal, 1, common.GetTimestampSecond(config.AuthED), "")
}

// GetSize 获取队列当前元素个数
func (G *FreeChatPool) GetSize() int {
	return G.queue.Len()
}

// GetCapacity 获取队列容量
func (G *FreeChatPool) GetCapacity() int {
	return G.capacity
}

// IsFull 检查队列是否已满
func (G *FreeChatPool) IsFull() bool {
	return G.queue.Len() == G.capacity
}

// AddFreeGpt35 入队
func (G *FreeChatPool) AddFreeGpt35(v *FreeChat.FreeChat) bool {
	if G.IsFull() || v == nil {
		return false
	}
	G.queue.Enqueue(v)
	return true
}
