package quantainer

import (
	"github.com/szmcdull/treemap/v2"
	"golang.org/x/exp/constraints"
)

/* SortedRingBuffer */

type SortedRingBuffer[T constraints.Ordered] struct {
	rb *RingBuffer[T]
	m  *treemap.TreeMap[T, int]
}

func NewSortedRingBuffer[T constraints.Ordered](size int) *SortedRingBuffer[T] {
	return &SortedRingBuffer[T]{
		rb: NewRingBuffer[T](size),
		m:  treemap.New[T, int](),
	}
}

// Cannot change size at runtime. Because rb.AddLast will potentially drop items if size decrease,
// but currently no way to sync m with rb in such situation.
//
// func NewSortedRingBufferConfigurable[T constraints.Ordered](configLoader func() int) *SortedRingBuffer[T] {
// 	return &SortedRingBuffer[T]{
// 		rb: NewRingBufferConfigurable[T](configLoader),
// 		m:  treemap.New[T, int](),
// 	}
// }

func (me *SortedRingBuffer[T]) AddLast(v T) {
	if me.rb.Full() {
		me.PopFirst()
	}
	me.rb.AddLast(v)
	me.addToTreeMap(v)
}

func (me *SortedRingBuffer[T]) removeFromTreeMap(v T) {
	ref, ok := me.m.GetRef(v)
	if ok {
		r := *ref
		r--
		if r == 0 {
			me.m.Del(v)
		} else {
			*ref = r
		}
	}
}

func (me *SortedRingBuffer[T]) addToTreeMap(v T) {
	ref, ok := me.m.GetRef(v)
	if !ok {
		me.m.Set(v, 1)
		return
	}
	r := *ref
	r++
	*ref = r
}

func (me *SortedRingBuffer[T]) PopFirst() *T {
	result := me.rb.PopFirst()
	if result != nil {
		me.removeFromTreeMap(*result)
	}
	return result
}

func (me *SortedRingBuffer[T]) ToSlice() []T {
	return me.rb.ToSlice()
}

// SortedSlice returns a sorted slice of the elements in the buffer.
// If cachedSlice is provided, it will be reused if large enough.
func (me *SortedRingBuffer[T]) SortedSlice(cachedSlice []T) []T {
	result := cachedSlice
	if result == nil || len(result) < me.rb.count {
		result = make([]T, me.rb.count)
	}
	ii := 0
	for i := me.m.Iterator(); i.Valid(); i.Next() {
		count := i.Value()
		key := i.Key()
		for j := 0; j < count; j++ {
			result[ii] = key
			ii++
		}
	}
	return result
}

func (me *SortedRingBuffer[T]) Count() int {
	return me.rb.count
}

func (me *SortedRingBuffer[T]) First() *T {
	return me.rb.First()
}

func (me *SortedRingBuffer[T]) Last() *T {
	return me.rb.Last()
}

func (me *SortedRingBuffer[T]) Clear() {
	me.rb.Clear()
	me.m.Clear()
}

func (me *SortedRingBuffer[T]) MaxSize() int {
	return me.rb.maxSize()
}

func (me *SortedRingBuffer[T]) Full() bool {
	return me.rb.Full()
}
