package quantainer

import "time"

type FixedDurationSlice[T any] struct {
	l           []T
	t           []int64
	maxDuration int64
}

func NewFixedDurationSlice[T any](maxDuration time.Duration) *FixedDurationSlice[T] {
	return &FixedDurationSlice[T]{
		maxDuration: maxDuration.Nanoseconds(),
	}
}

func (me *FixedDurationSlice[T]) Add(v T) {
	now := time.Now().UnixNano()
	me.l = append(me.l, v)
	me.t = append(me.t, now)
	cutoff := now - me.maxDuration
	i := 0
	l := len(me.t)
	for i < l && me.t[i] < cutoff {
		i++
	}
	if i > 0 {
		me.l = me.l[i-1:]
		me.t = me.t[i-1:]
	}
}

func (me *FixedDurationSlice[T]) Head() (t T, tm time.Time, ok bool) {
	if len(me.l) == 0 {
		return
	}
	return me.l[0], time.Unix(0, me.t[0]), true
}

func (me *FixedDurationSlice[T]) Values() []T {
	return me.l
}

func (me *FixedDurationSlice[T]) Clear() {
	me.l = nil
	me.t = nil
}
