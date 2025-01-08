package quantainer

import "fmt"

type List[T any] struct {
	front, back *Node[T]
	count       int
}

type Node[T any] struct {
	Value      T
	prev, next *Node[T]
}

func NewList[T any]() *List[T] {
	return &List[T]{}
}

func (me *List[T]) AddLast(v T) *Node[T] {
	node := &Node[T]{
		Value: v,
	}
	node.prev = me.back
	if me.back != nil {
		me.back.next = node
	} else {
		me.front = node
	}
	me.back = node
	me.count++
	return node
}

func (me *List[T]) AddFirst(v T) *Node[T] {
	node := &Node[T]{
		Value: v,
	}
	node.next = me.front
	if me.front != nil {
		me.front.prev = node
	} else {
		me.back = node
	}
	me.front = node
	me.count++
	return node
}

func (me *List[T]) AddAfter(prev *Node[T], v T) *Node[T] {
	node := &Node[T]{
		Value: v,
	}
	node.prev = prev
	node.next = prev.next
	prev.next = node
	if node.next == nil {
		me.back = node
	}
	me.count++
	return node
}

func (me *List[T]) AddBefore(next *Node[T], v T) *Node[T] {
	node := &Node[T]{
		Value: v,
	}
	node.next = next
	node.prev = next.prev
	next.prev = node
	if node.prev == nil {
		me.front = node
	}
	me.count++
	return node
}

// Remove removes the node from the list and returns the next node to it.
// It panics if the node is not in the list.
func (me *List[T]) Remove(node *Node[T]) (next *Node[T]) {
	if node.prev == nil {
		if me.front != node {
			panic("Node not in list")
		}
		me.front = node.next
	} else {
		node.prev.next = node.next
	}

	if node.next == nil {
		if me.back != node {
			panic("Node not in list")
		}
		me.back = node.prev
	} else {
		node.next.prev = node.prev
	}

	node.prev = nil
	next = node.next
	node.next = nil
	me.count--
	return
}

func (me *List[T]) First() *Node[T] {
	return me.front
}

func (me *List[T]) Last() *Node[T] {
	return me.back
}

func (me *List[T]) setFront(n *Node[T]) {
	me.front = n
	if n != nil {
		n.prev = nil
	} else {
		me.back = nil
	}
}

func (me *List[T]) setBack(n *Node[T]) {
	me.back = n
	if n != nil {
		n.next = nil
	} else {
		me.front = nil
	}
}

func (me *List[T]) PopFirst() *Node[T] {
	node := me.front
	if node == nil {
		return nil
	}

	me.setFront(node.next)
	node.next = nil
	me.count--
	return node
}

func (me *List[T]) PopLast() *Node[T] {
	node := me.back
	if node == nil {
		return nil
	}
	me.setBack(node.prev)
	node.prev = nil
	me.count--
	return node
}

func (me *List[T]) Clear() {
	me.front = nil
	me.back = nil
	me.count = 0
}

func (me *List[T]) Count() int {
	return me.count
}

func (me *List[T]) _Count() int {
	first := me.First()
	count := 0
	for ; first != nil; first = first.next {
		count++
	}
	return count
}

func (me *Node[T]) Next() *Node[T] {
	return me.next
}

func (me *Node[T]) Prev() *Node[T] {
	return me.prev
}

func FromSlice[T any](slice []T) *List[T] {
	list := NewList[T]()
	for _, v := range slice {
		list.AddLast(v)
	}
	return list
}

func (me *List[T]) ToSlice() []T {
	slice := make([]T, me.count)
	i := 0
	for n := me.First(); n != nil; n = n.Next() {
		slice[i] = n.Value
		i++
	}
	return slice
}

func (me *List[T]) PopFirstWhen(fn func(v *T) bool) {
	n := me.First()
	if n == nil {
		return
	}

	i := 0
	for ; n != nil; n = n.next {
		if !fn(&n.Value) {
			break
		}
		i++
	}
	me.setFront(n)
	me.count -= i
}

func (me *List[T]) PopLastWhen(fn func(v *T) bool) {
	n := me.Last()
	if n == nil {
		return
	}

	i := 0
	for ; n != nil; n = n.prev {
		if !fn(&n.Value) {
			break
		}
		i++
	}
	me.setBack(n)
	me.count -= i
}

var ErrorIndexOutOfRange = fmt.Errorf("index out of range")

// Trim removes elements from the list starting from the start index to the end index.
// The start index is inclusive and the end index is exclusive.
// Negative indices are counted from the end of the list.
func (me *List[T]) Trim(start, end int) {
	l := me.count
	absStart := start
	absEnd := end
	if absStart < 0 {
		absStart = l + absStart
	}
	if absEnd < 0 {
		absEnd = l + absEnd
	}
	if absStart > absEnd {
		panic(fmt.Errorf("start index %d is greater than end index %d", start, end))
	}
	if absStart < 0 || absEnd < 0 || absStart >= l || absEnd > l {
		panic(ErrorIndexOutOfRange)
	}
	if me.count == 0 {
		return
	}
	if absStart == absEnd {
		me.front = nil
		me.back = nil
		me.count = 0
		return
	}
	front := me.At(start)
	back := me.At(end - 1)
	if front.prev != nil {
		front.prev = nil
	}
	me.front = front
	if back.next != nil {
		back.next = nil
	}
	me.back = back
	me.count = absEnd - absStart
}

// At returns the node at the given index. Negative indices are counted from the end of the list.
func (me *List[T]) At(i int) *Node[T] {
	if i < 0 {
		return me.fromBack(-i - 1)
	}
	if i < 0 || i >= me.count {
		return nil
	}
	n := me.front
	for ; i > 0; i-- {
		n = n.next
	}
	return n
}

func (me *List[T]) fromBack(i int) *Node[T] {
	if i < 0 || i >= me.count {
		return nil
	}
	n := me.back
	for ; i > 0; i-- {
		n = n.prev
	}
	return n
}
