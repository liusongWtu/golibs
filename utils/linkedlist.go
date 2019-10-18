package utils

import (
	"fmt"
	"github.com/cheekybits/genny/generic"
	"sync"
)

type Item generic.Type

type Node struct {
	content Item
	next    *Node
}

type LinkedList struct {
	head *Node
	size int
	lock sync.RWMutex
}

// 在链表结尾追加元素
func (list *LinkedList) Append(t Item) {
	list.lock.Lock()
	newNode := Node{t, nil}
	// 查找并追加
	if list.head == nil { // 空链表第一次追加元素
		list.head = &newNode
	} else {
		curNode := list.head // 遍历链表，找到尾部结点
		for {
			if curNode.next == nil {
				break
			}
			curNode = curNode.next
		}
		curNode.next = &newNode
	}
	// 追加后链表长度+1
	list.size++
	list.lock.Unlock()
}

// 在链表指定位置插入指定元素
func (list *LinkedList) Insert(i int, t Item) error {
	list.lock.Lock()
	defer list.lock.Unlock()
	if i < 0 || i > list.size {
		return fmt.Errorf("Index %d out of bonuds", i)
	}
	newNode := Node{t, nil}

	if i == 0 { // 插入到链表头部
		newNode.next = list.head
		list.head = &newNode
		list.size++
		return nil
	}

	preNode := list.head
	preIndex := 0
	for preIndex < i-2 {
		preIndex++
		preNode = preNode.next
	}
	// 执行插入
	newNode.next = preNode.next
	preNode.next = &newNode
	list.size++
	return nil
}

// 删除指定位置的元素
func (list *LinkedList) RemoveAt(i int) (*Item, error) {
	list.lock.Lock()
	defer list.lock.Unlock()

	if i < 0 || i > list.size {
		return nil, fmt.Errorf("Index %d out of bonuds", i)
	}

	curNode := list.head
	preIndex := 0
	for preIndex < i-1 {
		preIndex++
		curNode = curNode.next
	}
	item := curNode.content
	curNode.next = curNode.next.next
	list.size--
	return &item, nil
}

// 获取指定元素在链表中的索引
func (list *LinkedList) IndexOf(t Item) int {
	list.lock.RLock()
	defer list.lock.RUnlock()
	curNode := list.head
	locIndex := 0
	for {
		if curNode.content == t {
			return locIndex
		}
		if curNode.next == nil {
			return -1
		}
		curNode = curNode.next
		locIndex++
	}
}

// 检查链表是否为空
func (list *LinkedList) IsEmpty() bool {
	list.lock.RLock()
	defer list.lock.RUnlock()
	if list.head == nil {
		return true
	}
	return false
}

// 获取链表的长度
func (list *LinkedList) Size() int {
	list.lock.RLock()
	defer list.lock.RUnlock()
	size := 1
	nextNode := list.head
	for {
		if nextNode == nil || nextNode.next == nil { // 单结点链表的 nextNode == nil
			break
		}
		size++
		nextNode = nextNode.next
	}
	return size
}

// 格式化打印链表
func (list *LinkedList) String() {
	list.lock.RLock()
	defer list.lock.RUnlock()
	curNode := list.head
	for {
		if curNode == nil {
			break
		}
		print(curNode.content)
		print(" ")
		curNode = curNode.next
	}
	println()
}

// 获取链表的头结点
func (list *LinkedList) Head() *Node {
	list.lock.RLock()
	defer list.lock.RUnlock()
	return list.head
}
