//go:build go1.23
// +build go1.23

package quantainer

import "iter"

func (me *RingBuffer[T]) All() iter.Seq[T] {
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

func (me *SortedRingBuffer[T]) All() iter.Seq[T] {
	return me.rb.All()
}
