package quantainer

import (
	"github.com/szmcdull/treemap/v2"
	"golang.org/x/exp/constraints"
)

/* SortedFixedList */

type SortedFixedList[T constraints.Ordered] struct {
	*FixedList[T]
	*treemap.TreeMap[T, int]
}

func NewSortedFixedList[T constraints.Ordered](size int) *SortedFixedList[T] {
	return &SortedFixedList[T]{
		FixedList: NewFixedList[T](size),
		TreeMap:   treemap.New[T, int](),
	}
}

func NewSortedFixedListConfigurable[T constraints.Ordered](configLoader func() int) *SortedFixedList[T] {
	return &SortedFixedList[T]{
		FixedList: NewFixedListConfigurable[T](configLoader),
		TreeMap:   treemap.New[T, int](),
	}
}

func (me *SortedFixedList[T]) AddFirst(v T) {
	me.l.AddFirst(v)
	me.addToTreeMap(v)
	for me.l.count > me.maxSize() {
		me.PopLast()
	}
}

func (me *SortedFixedList[T]) AddLast(v T) {
	me.l.AddLast(v)
	me.addToTreeMap(v)
	for me.l.count > me.maxSize() {
		me.PopFirst()
	}
}

func (me *SortedFixedList[T]) removeFromTreeMap(v T) {
	ref, ok := me.TreeMap.GetRef(v)
	if ok {
		r := *ref
		r--
		if r == 0 {
			me.TreeMap.Del(v)
		} else {
			*ref = r
		}
	}
}

func (me *SortedFixedList[T]) addToTreeMap(v T) {
	ref, ok := me.TreeMap.GetRef(v)
	if !ok {
		me.TreeMap.Set(v, 1)
		return
	}
	r := *ref
	r++
	*ref = r
}

func (me *SortedFixedList[T]) Remove(node *Node[T]) (next *Node[T]) {
	me.removeFromTreeMap(node.Value)
	return me.FixedList.Remove(node)
}

func (me *SortedFixedList[T]) PopFirst() *Node[T] {
	result := me.FixedList.PopFirst()
	if result != nil {
		me.removeFromTreeMap(result.Value)
	}
	return result
}

func (me *SortedFixedList[T]) PopLast() *Node[T] {
	result := me.FixedList.PopLast()
	if result != nil {
		me.removeFromTreeMap(result.Value)
	}
	return result
}

func (me *SortedFixedList[T]) ToSlice() []T {
	return me.l.ToSlice()
}

func (me *SortedFixedList[T]) SortedSlice() []T {
	result := make([]T, me.l.count)
	ii := 0
	for i := me.TreeMap.Iterator(); i.Valid(); i.Next() {
		count := i.Value()
		key := i.Key()
		for j := 0; j < count; j++ {
			result[ii] = key
			ii++
		}
	}
	return result
}

func (me *SortedFixedList[T]) Count() int {
	return me.l.count
}

func (me *SortedFixedList[T]) First() *Node[T] {
	return me.l.First()
}

func (me *SortedFixedList[T]) Last() *Node[T] {
	return me.l.Last()
}

func (me *SortedFixedList[T]) Clear() {
	me.l.Clear()
	me.TreeMap.Clear()
}

func (me *SortedFixedList[T]) MaxSize() int {
	return me.maxSize()
}

func (me *SortedFixedList[T]) Full() bool {
	return me.l.count == me.maxSize()
}
