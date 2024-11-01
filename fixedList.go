package quantainer

/* FixedList */

type FixedList[T any] struct {
	l       List[T]
	maxSize func() int
}

func NewFixedList[T any](size int) *FixedList[T] {
	return &FixedList[T]{
		maxSize: func() int { return size },
	}
}

func NewFixedListConfigurable[T any](configLoader func() int) *FixedList[T] {
	return &FixedList[T]{
		maxSize: configLoader,
	}
}

func (me *FixedList[T]) AddFirst(v T) {
	me.l.AddFirst(v)
	for me.l.count > me.maxSize() {
		me.l.PopLast()
	}
}

func (me *FixedList[T]) AddLast(v T) {
	me.l.AddLast(v)
	for me.l.count > me.maxSize() {
		me.l.PopFirst()
	}
}

func (me *FixedList[T]) Remove(node *Node[T]) (next *Node[T]) {
	return me.l.Remove(node)
}

func (me *FixedList[T]) PopFirst() *Node[T] {
	return me.l.PopFirst()
}

func (me *FixedList[T]) PopLast() *Node[T] {
	return me.l.PopLast()
}

func (me *FixedList[T]) ToSlice() []T {
	return me.l.ToSlice()
}

func (me *FixedList[T]) Count() int {
	return me.l.count
}

func (me *FixedList[T]) First() *Node[T] {
	return me.l.First()
}

func (me *FixedList[T]) Last() *Node[T] {
	return me.l.Last()
}

func (me *FixedList[T]) Clear() {
	me.l.Clear()
}

func (me *FixedList[T]) MaxSize() int {
	return me.maxSize()
}

func (me *FixedList[T]) Full() bool {
	return me.l.count == me.maxSize()
}
