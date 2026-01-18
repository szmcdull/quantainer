package quantainer

type (
	// BusyPollRingBuffer is a simple fixed-size, write-only ring buffer designed for busy-polling scenarios.
	// It should be used where only one goroutine is writing data and one or more goroutines are reading data.
	// The writer can overwrite old data when the buffer is full.
	// There is no synchronization between the readers and the writer. Readers would miss data if they are too slow.
	BusyPollRingBuffer[T any] struct {
		buf      []T
		from, to int
	}

	BusyPollRingBufferReader[T any] struct {
		buf     *BusyPollRingBuffer[T]
		current int
	}
)

func NewBusyPollRingBuffer[T any](size int) *BusyPollRingBuffer[T] {
	return &BusyPollRingBuffer[T]{
		buf: make([]T, size),
	}
}

func (me *BusyPollRingBuffer[T]) Len() int {
	return len(me.buf)
}

func (me *BusyPollRingBuffer[T]) Write(v T) {
	me.buf[me.to] = v
	me.from = me.to
	me.to++
	if me.to >= len(me.buf) {
		me.to = 0
	}
}

func (me *BusyPollRingBuffer[T]) BeginWrite() *T {
	ptr := &me.buf[me.to]
	return ptr
}

func (me *BusyPollRingBuffer[T]) EndWrite() {
	to := me.to
	me.from = to
	to++
	if to >= len(me.buf) {
		to = 0
	}
	me.to = to
}

func (me *BusyPollRingBuffer[T]) WriteMany(vs []T) {
	from := me.to
	to := from
	size := len(me.buf)
	for _, v := range vs {
		me.buf[to] = v
		to++
		if to >= size {
			to = 0
		}
	}
	me.to = to
	me.from = from
}

// Reader creates a new reader for the ring buffer.
// There is no synchronization between the writer and the reader.
// Usually the reader would busy poll the buffer, or spin-wait until new data is available.
func (me *BusyPollRingBuffer[T]) Reader() *BusyPollRingBufferReader[T] {
	return &BusyPollRingBufferReader[T]{
		buf:     me,
		current: me.from,
	}
}

func (me *BusyPollRingBufferReader[T]) Read(result *T) (ok bool) {
	if me.current == me.buf.to {
		return
	}
	*result = me.buf.buf[me.current]
	ok = true
	me.current++
	if me.current >= len(me.buf.buf) {
		me.current = 0
	}
	return
}
