package queue

type (
	Queue struct {
		start, end *Node
		length     int
	}
	Node struct {
		Value interface{}
		next  *Node
	}
)

// New 新建一个队列
func New() *Queue {
	return &Queue{nil, nil, 0}
}

// Dequeue 出队
func (Q *Queue) Dequeue() *Node {
	if Q.length == 0 {
		return nil
	}
	n := Q.start
	if Q.length == 1 {
		Q.start = nil
		Q.end = nil
	} else {
		Q.start = Q.start.next
	}
	Q.length--
	return n
}

// Enqueue 入队
func (Q *Queue) Enqueue(value interface{}) {
	n := &Node{value, nil}
	if Q.length == 0 {
		Q.start = n
		Q.end = n
	} else {
		Q.end.next = n
		Q.end = n
	}
	Q.length++
}

// Len 获取队列长度
func (Q *Queue) Len() int {
	return Q.length
}

// Peek 返回队列的第一个元素
func (Q *Queue) Peek() *Node {
	if Q.length == 0 {
		return nil
	}
	return Q.start
}

// Remove 移除指定节点
func (Q *Queue) Remove(n *Node) {
	if Q.length == 0 || n == nil {
		return
	}

	// 如果移除的是队列的第一个元素
	if n == Q.start {
		Q.start = Q.start.next
		if Q.start == nil {
			// 如果移除后队列为空，则end也应该设置为nil
			Q.end = nil
		}
		Q.length--
		return
	}

	// 找到n的前一个节点
	prevNode := Q.start
	for prevNode != nil && prevNode.next != n {
		prevNode = prevNode.next
	}

	if prevNode == nil {
		// 没有找到n的前一个节点（n不在队列中）
		return
	}

	// 移除节点n
	prevNode.next = n.next
	// 如果移除的是最后一个元素，更新end指针
	if n.next == nil {
		Q.end = prevNode
	}
	Q.length--
}

// Traverse 遍历队列
func (Q *Queue) Traverse(cb func(n *Node)) {
	if Q.length == 0 {
		return
	}
	for n := Q.start; n != nil; n = n.next {
		cb(n)
	}
}
