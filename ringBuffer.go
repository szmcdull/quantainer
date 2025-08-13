package quantainer

import (
	"cmp"
	"slices"

	"github.com/szmcdull/glinq/garray"
)

type (
	// RingBuffer is a circular buffer with a fixed maximum size.
	// It supports adding elements to the end and removing elements from the front.
	// If the buffer is full, adding a new element will overwrite the oldest element.
	RingBuffer[T cmp.Ordered] struct {
		data    []T
		tail    int // Points to the next available position
		count   int // Number of elements in the buffer
		maxSize func() int
	}
)

func NewRingBuffer[T cmp.Ordered](size int) *RingBuffer[T] {
	return &RingBuffer[T]{
		data: make([]T, size),
		maxSize: func() int {
			return size
		},
	}
}

// NewRingBufferConfigurable creates a new ring buffer with a configurable maximum size.
// If maxSize is changed at runtime, the buffer will not resize until you add an element.
func NewRingBufferConfigurable[T cmp.Ordered](maxSize func() int) *RingBuffer[T] {
	size := maxSize()
	return &RingBuffer[T]{
		data:    make([]T, size),
		maxSize: maxSize,
	}
}

func (me *RingBuffer[T]) Count() int {
	return me.count
}

func (me *RingBuffer[T]) Full() bool {
	return me.count >= me.maxSize()
}

func (me *RingBuffer[T]) head() int {
	result := me.tail - me.count
	if result < 0 {
		result += len(me.data)
	}
	return result
}

// AddLast adds an element to the end of the ring buffer.
// If the buffer is full, it will overwrite the oldest element.
func (me *RingBuffer[T]) AddLast(v T) {
	size := me.maxSize()
	oldSize := len(me.data)
	head := me.head()
	tail := me.tail
	if sizeDiff := size - oldSize; sizeDiff > 0 { // maxSize enlarged
		data := append(me.data, make([]T, sizeDiff)...)
		if tail < head || me.count == oldSize {
			if tail > sizeDiff {
				copy(data[oldSize:], me.data[:sizeDiff])
				copy(data, me.data[sizeDiff:tail])
				tail -= sizeDiff
			} else {
				copy(data[oldSize:], me.data[:tail]) // move wrapped tail data to new end
				tail = (tail + oldSize) % size
			}
		}
		me.data = data
	} else if sizeDiff < 0 { // maxSize reduced
		count := me.count
		me.count = size
		me.data = me.ToSlice()
		if count > size {
			count = size
		}
		me.count = count
		tail = 0
	}

	me.data[tail] = v
	tail++
	if me.count < size {
		me.count++
	}

	if tail >= size {
		tail -= size
	}
	me.tail = tail
}

func (me *RingBuffer[T]) PopFirst() (result *T) {
	if me.count == 0 {
		return nil
	}
	head := me.head()
	result = &me.data[head]
	me.count--
	return
}

func (me *RingBuffer[T]) First() (result *T) {
	if me.count == 0 {
		return nil
	}
	return &me.data[me.head()]
}

func (me *RingBuffer[T]) Last() (result *T) {
	if me.count == 0 {
		return nil
	}
	i := me.tail - 1
	if i < 0 {
		i += len(me.data)
	}
	return &me.data[i]
}

func (me *RingBuffer[T]) At(i int) (result *T) {
	l := len(me.data)

	if i < 0 {
		i = me.tail + i
	} else {
		i = me.tail - me.count + i // i += head
	}
	if i < 0 {
		i += l
	}

	if i < 0 || i >= l { // allow to access unfilled elements?
		return nil
	}

	return &me.data[i]
}

func (me *RingBuffer[T]) toSlice(slice []T) {
	if len(slice) != me.count {
		panic("length of provided slice does not match count of elements in ring buffer")
	}

	// Nothing to copy for empty buffer
	if me.count == 0 {
		return
	}

	head := me.head()
	tail := me.tail
	data := me.data

	if tail > head {
		copy(slice, data[head:tail])
	} else {
		copy(slice, data[head:])
		copy(slice[len(data[head:]):], data[:tail])
	}
}

func (me *RingBuffer[T]) ToSlice() []T {
	result := make([]T, me.count)
	me.toSlice(result)
	return result
}

// ToSliceAndSort copies the elements of the ring buffer to a slice (optional nil at first, and reuse after) and sorts it.
// s must not contain NaN, or the result is undefined.
func (me *RingBuffer[T]) ToSliceAndSort(s []T) []T {
	if s == nil || len(s) < me.count {
		s = me.ToSlice()
	} else {
		me.toSlice(s)
	}
	garray.Sort(s)
	return s
}

// ToSliceAndSort copies the elements of the ring buffer to a slice (optional nil at first, and reuse after) and sorts it.
func (me *RingBuffer[T]) ToSliceAndSortNaN(s []T) []T {
	if s == nil || len(s) < me.count {
		s = me.ToSlice()
	} else {
		me.toSlice(s)
	}
	slices.Sort(s)
	return s
}

func (me *RingBuffer[T]) Clear() {
	me.tail = 0
	me.count = 0
}
