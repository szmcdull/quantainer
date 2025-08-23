//go:build go1.23
// +build go1.23

package quantainer

import "iter"

func (me *RingBuffer[T]) Values() iter.Seq[T] {
	return func(yield func(T) bool) {
		p := me.head()
		for range me.count {
			if !yield(me.data[p]) {
				return
			}
			p++
			if p >= len(me.data) {
				p = 0
			}
		}
	}
}

func (me *RingBuffer[T]) All() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		p := me.head()
		for i := 0; i < me.count; i++ {
			if !yield(i, me.data[p]) {
				return
			}
			p++
			if p >= len(me.data) {
				p = 0
			}
		}
	}
}

func (me *SortedRingBuffer[T]) Values() iter.Seq[T] {
	return me.rb.Values()
}

func (me *SortedRingBuffer[T]) All() iter.Seq2[int, T] {
	return me.rb.All()
}
