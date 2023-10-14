package msgpool

import (
	"container/ring"
	"sync"
)

// ConcurrentQueue 是一个并发安全的队列数据结构
type ConcurrentQueue[T any] struct {
	data *ring.Ring
	size int
	mu   sync.Mutex
}

type QueueResult[T any] struct {
	IsNil bool
	Value T
}

// NewConcurrentQueue 创建一个新的并发安全队列
func NewConcurrentQueue[T any]() *ConcurrentQueue[T] {
	return &ConcurrentQueue[T]{
		data: ring.New(0),
		size: 0,
	}
}

// Enqueue 将元素添加到队列尾部
func (q *ConcurrentQueue[T]) Enqueue(value T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.data = q.data.Link(&ring.Ring{Value: value})
	q.size++
}

// Dequeue 从队列头部移除并返回元素
func (q *ConcurrentQueue[T]) Dequeue() QueueResult[T] {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.size == 0 {
		return QueueResult[T]{IsNil: true}
	}
	value := q.data.Value
	q.data = q.data.Next()
	q.size--
	return QueueResult[T]{IsNil: false, Value: value.(T)}
}

// Size 返回队列的大小
func (q *ConcurrentQueue[T]) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.size
}
