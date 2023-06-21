package quantainer

/* FixedList */

type FixedList[T any] struct {
	l       List[T]
	maxSize int
}

func NewFixedList[T any](size int) *FixedList[T] {
	return &FixedList[T]{
		maxSize: size,
	}
}

func (me *FixedList[T]) AddFirst(v T) {
	me.l.AddFirst(v)
	for me.l.count > me.maxSize {
		me.l.PopLast()
	}
}

func (me *FixedList[T]) AddLast(v T) {
	me.l.AddLast(v)
	for me.l.count > me.maxSize {
		me.l.PopFirst()
	}
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
