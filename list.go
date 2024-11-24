package quantainer

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
