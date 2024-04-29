package FreeGpt35Pool

import (
	"fmt"
	"free-gpt3.5-2api/FreeGpt35"
	"free-gpt3.5-2api/common"
	"free-gpt3.5-2api/config"
	"github.com/aurorax-neo/go-logger"
	"sync"
	"time"
)

var (
	instance *FreeGpt35Pool
	once     sync.Once
)

type FreeGpt35Pool struct {
	data     []*FreeGpt35.Gpt35
	head     int // 队头指针
	tail     int // 队尾指针
	size     int // 队列当前元素个数
	capacity int // 队列容量
}

func newGpt35Pool(capacity int) *FreeGpt35Pool {
	return &FreeGpt35Pool{
		data:     make([]*FreeGpt35.Gpt35, capacity),
		capacity: config.PoolMaxCount,
		size:     0,
		head:     0,
		tail:     0,
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
		if !G.isLiveGpt35(gpt35) {
			// 入队新 FreeGpt35 实例
			G.Enqueue(gpt35)
		}
	})
	// 检测并移除无效 FreeGpt35 实例
	common.AsyncLoopTask(sleep, func() {
		// 遍历队列中的所有元素
		G.Traverse(func(index int, gpt35 *FreeGpt35.Gpt35) {
			// 判断是否为无效 FreeGpt35 实例
			if !G.isLiveGpt35(gpt35) {
				// 移除无效 FreeGpt35 实例
				G.RemoveAt(index)
			}
		})
	})
}

func (G *FreeGpt35Pool) isLiveGpt35(gpt35 *FreeGpt35.Gpt35) bool {
	//判断是否为空
	if gpt35 == nil ||
		gpt35.MaxUseCount <= 0 || //无可用次数
		gpt35.ExpiresIn <= common.GetTimestampSecond(0) {
		return false
	}
	return true
}

func (G *FreeGpt35Pool) GetGpt35(retry int) *FreeGpt35.Gpt35 {
	// 获取 FreeGpt35 实例
	gpt35 := G.Front()
	if G.isLiveGpt35(gpt35) { //有缓存
		// 减少 FreeGpt35 实例的最大使用次数
		gpt35.MaxUseCount--
		// 判断 FreeGpt35 实例是否有效 无效则移除
		if G.isLiveGpt35(gpt35) {
			G.Dequeue()
		}
		// 深拷贝 FreeGpt35 实例
		return common.DeepCopyStruct(gpt35).(*FreeGpt35.Gpt35)
	} else if retry > 0 {
		return G.GetGpt35(retry - 1)
	}
	// 缓存内无可用 FreeGpt35 实例，返回新 FreeGpt35 实例
	return FreeGpt35.NewGpt35(0)
}

// GetSize 获取队列当前元素个数
func (G *FreeGpt35Pool) GetSize() int {
	return G.size
}

// GetCapacity 获取队列容量
func (G *FreeGpt35Pool) GetCapacity() int {
	return G.capacity
}

// IsFull 检查队列是否已满
func (G *FreeGpt35Pool) IsFull() bool {
	return G.size == G.capacity
}

// Enqueue 入队
func (G *FreeGpt35Pool) Enqueue(v *FreeGpt35.Gpt35) bool {
	if G.IsFull() || v == nil {
		return false
	}
	G.data[G.tail] = v
	G.tail = (G.tail + 1) % G.capacity
	G.size++
	return true
}

// IsEmpty 检查队列是否为空
func (G *FreeGpt35Pool) IsEmpty() bool {
	return G.size == 0
}

// Dequeue 出队
func (G *FreeGpt35Pool) Dequeue() *FreeGpt35.Gpt35 {
	// 判断是否为空
	if G.IsEmpty() {
		return nil
	}
	value := G.data[G.head]
	G.head = (G.head + 1) % G.capacity
	G.size--
	return value
}

// Front 获取队首元素
func (G *FreeGpt35Pool) Front() *FreeGpt35.Gpt35 {
	if G.IsEmpty() {
		return nil
	}
	return G.data[G.head]
}

// Rear 获取队尾元素
func (G *FreeGpt35Pool) Rear() *FreeGpt35.Gpt35 {
	if G.IsEmpty() {
		return nil
	}
	// 需要计算tail的上一个位置
	tailIndex := (G.tail - 1 + G.capacity) % G.capacity
	return G.data[tailIndex]
}

// RemoveAt 移除指定位置的元素
func (G *FreeGpt35Pool) RemoveAt(index int) (*FreeGpt35.Gpt35, bool) {
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
func (G *FreeGpt35Pool) Traverse(callback func(int, *FreeGpt35.Gpt35)) {
	if G.IsEmpty() {
		return
	}
	// 从队头开始遍历到队尾
	for i := 0; i < G.size; i++ {
		index := (G.head + i) % G.capacity
		callback(index, G.data[index])
	}
}
