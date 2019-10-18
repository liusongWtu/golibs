package utils

import (
	"github.com/cheekybits/genny/generic"
	"sync"
)

type QueueItem generic.Type

type Queue struct {
	items []QueueItem
	lock  sync.RWMutex
}

// 创建队列
func (q *Queue) New() *Queue {
	q.items = []QueueItem{}
	return q
}

// 入队列
func (q *Queue) Enqueue(t QueueItem) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.items = append(q.items, t)
}

// 出队列
func (q *Queue) Dequeue() *QueueItem {
	q.lock.Lock()
	defer q.lock.Unlock()

	item := q.items[0]
	q.items = q.items[1:len(q.items)]
	return &item
}

// 获取队列的第一个元素，不移除
func (q *Queue) Front() *QueueItem {
	q.lock.Lock()
	defer q.lock.Unlock()

	item := q.items[0]
	return &item
}

// 判空
func (q *Queue) IsEmpty() bool {
	q.lock.Lock()
	defer q.lock.Unlock()

	return len(q.items) == 0
}

// 获取队列的长度
func (q *Queue) Size() int {
	q.lock.Lock()
	defer q.lock.Unlock()

	return len(q.items)
}
