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
	data     []*chat.Gpt35
	head     int // 队头指针
	tail     int // 队尾指针
	size     int // 队列当前元素个数
	capacity int // 队列容量
}

func newGpt35Pool(capacity int) *Gpt35Pool {
	return &Gpt35Pool{
		data:     make([]*chat.Gpt35, capacity),
		capacity: config.PoolMaxCount,
	}
}

func GetGpt35PoolInstance() *Gpt35Pool {
	once.Do(func() {
		instance = newGpt35Pool(config.PoolMaxCount)
		logger.Logger.Info(fmt.Sprint("PoolMaxCount: ", config.PoolMaxCount, ", AuthExpirationDate: ", config.AuthED, ", Init Pool..."))
		// 定时刷新 Pool
		instance.updateGpt35Pool(time.Millisecond * 200)
	})
	return instance
}

func (G *Gpt35Pool) updateGpt35Pool(nanosecond time.Duration) {
	common.TimingTask(nanosecond, func() {
		// 遍历队列中的所有元素
		G.Traverse(func(index int, gpt35 *chat.Gpt35) {
			// 判断是否为无效 Gpt35 实例
			if !G.isLiveGpt35(gpt35) {
				// 移除无效 Gpt35 实例
				G.RemoveAt(index)
			}
		})
		if !G.IsFull() {
			G.Enqueue(chat.NewGpt35(1))
		}
	})
}

func (G *Gpt35Pool) isLiveGpt35(gpt35 *chat.Gpt35) bool {
	//判断是否为空
	if gpt35 == nil ||
		gpt35.MaxUseCount <= 0 || //无可用次数
		gpt35.ExpiresIn <= common.GetTimestampSecond(0) {
		return false
	}
	return true
}

func (G *Gpt35Pool) GetGpt35(retry int) *chat.Gpt35 {
	// 获取 Gpt35 实例
	gpt35 := G.Dequeue()
	if gpt35 != nil { //有缓存
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
		return &gpt35_
	} else if retry > 0 {
		// 递归获取 Gpt35 实例
		time.Sleep(time.Millisecond * 200)
		return G.GetGpt35(retry - 1)
	}
	// 缓存内无可用 Gpt35 实例，返回新 Gpt35 实例
	return chat.NewGpt35(0)
}

// GetSize 获取队列当前元素个数
func (G *Gpt35Pool) GetSize() int {
	return G.size
}

// GetCapacity 获取队列容量
func (G *Gpt35Pool) GetCapacity() int {
	return G.capacity
}

// IsFull 检查队列是否已满
func (G *Gpt35Pool) IsFull() bool {
	return G.size == G.capacity
}

// Enqueue 入队
func (G *Gpt35Pool) Enqueue(v *chat.Gpt35) bool {
	if G.IsFull() || v == nil {
		return false
	}
	G.data[G.tail] = v
	G.tail = (G.tail + 1) % G.capacity
	G.size++
	return true
}

// IsEmpty 检查队列是否为空
func (G *Gpt35Pool) IsEmpty() bool {
	return G.size == 0
}

// Dequeue 出队
func (G *Gpt35Pool) Dequeue() *chat.Gpt35 {
	// 判断是否为空
	if G.IsEmpty() {
		return nil
	}
	// 获取 Gpt35 实例
	gpt35 := G.data[G.head]
	// 判断是否为无效 Gpt35 实例
	if !G.isLiveGpt35(gpt35) {
		G.head = (G.head + 1) % G.capacity
		G.size--
		return nil
	}
	return gpt35
}

// RemoveAt 移除指定位置的元素
func (G *Gpt35Pool) RemoveAt(index int) (*chat.Gpt35, bool) {
	if index < 0 || index >= G.size {
		return nil, false
	}
	// 计算要移除的元素在数组中的索引
	removeIndex := (G.head + index) % G.capacity
	removedValue := G.data[removeIndex]

	// 移动队列中被移除元素后面的元素
	for i := index; i < G.size-1; i++ {
		currentIndex := (G.head + i) % G.capacity
		nextIndex := (currentIndex + 1) % G.capacity
		G.data[currentIndex] = G.data[nextIndex]
	}
	// 将最后一个元素置为空
	emptyIndex := (G.head + G.size - 1) % G.capacity
	G.data[emptyIndex] = nil

	// 更新队尾指针和元素个数
	G.tail = (G.tail - 1 + G.capacity) % G.capacity
	G.size--
	return removedValue, true
}

// Traverse 遍历队列中的所有元素，并对每个元素执行指定操作
func (G *Gpt35Pool) Traverse(callback func(int, *chat.Gpt35)) {
	if G.IsEmpty() {
		return
	}
	// 从队头开始遍历到队尾
	for i := 0; i < G.size; i++ {
		index := (G.head + i) % G.capacity
		callback(index, G.data[index])
	}
}
