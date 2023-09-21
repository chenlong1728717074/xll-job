package util

import (
	"sync"
)

type LinkedList[T any] struct {
	lock  sync.RWMutex
	len   int
	arr   []T
	index int
}

func NewLinkedList[T any]() *LinkedList[T] {
	return &LinkedList[T]{
		lock: sync.RWMutex{},
		arr:  make([]T, 0),
	}
}

func (do *LinkedList[T]) put(n T) {
	do.lock.Lock()
	defer do.lock.Unlock()
	do.arr = append(do.arr, n)
	do.len++
}
