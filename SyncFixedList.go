package quantainer

import "sync"

/* FixedList */

type SyncFixedList[T any] struct {
	sync.Mutex
	l FixedList[T]
}

func NewSyncFixedList[T any](size int) *SyncFixedList[T] {
	return &SyncFixedList[T]{
		l: FixedList[T]{
			maxSize: func() int { return size },
		},
	}
}

func NewSyncFixedListConfigurable[T any](configLoader func() int) *SyncFixedList[T] {
	return &SyncFixedList[T]{
		l: FixedList[T]{
			maxSize: configLoader,
		},
	}
}

func (me *SyncFixedList[T]) AddFirst(v T) {
	me.Lock()
	defer me.Unlock()
	me.l.AddFirst(v)
}

func (me *SyncFixedList[T]) AddLast(v T) {
	me.Lock()
	defer me.Unlock()
	me.l.AddLast(v)
}

func (me *SyncFixedList[T]) ToSlice() []T {
	me.Lock()
	defer me.Unlock()
	return me.l.ToSlice()
}

func (me *SyncFixedList[T]) Count() int {
	return me.l.l.count
}

func (me *SyncFixedList[T]) First() *Node[T] {
	return me.l.l.front
}

func (me *SyncFixedList[T]) Last() *Node[T] {
	return me.l.l.back
}

func (me *SyncFixedList[T]) Clear() {
	me.Lock()
	defer me.Unlock()
	me.l.Clear()
}
