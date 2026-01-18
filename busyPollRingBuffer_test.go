package quantainer

import (
	"testing"
)

func TestNewBusyPollRingBuffer(t *testing.T) {
	buf := NewBusyPollRingBuffer[int](5)
	if buf == nil {
		t.Fatal("NewBusyPollRingBuffer returned nil")
	}
	if buf.Len() != 5 {
		t.Errorf("Expected length 5, got %d", buf.Len())
	}
}

func TestBusyPollRingBuffer_Write(t *testing.T) {
	buf := NewBusyPollRingBuffer[int](3)
	reader := buf.Reader()

	// Write a single value
	buf.Write(10)
	val, ok := reader.Read()
	if !ok {
		t.Fatal("Expected to read a value")
	}
	if val != 10 {
		t.Errorf("Expected 10, got %d", val)
	}

	// No more data available
	_, ok = reader.Read()
	if ok {
		t.Error("Expected no more data")
	}
}

func TestBusyPollRingBuffer_WriteMultiple(t *testing.T) {
	buf := NewBusyPollRingBuffer[int](5)
	reader := buf.Reader()

	// Write multiple values
	buf.Write(1)
	buf.Write(2)
	buf.Write(3)

	// Read them back
	val, ok := reader.Read()
	if !ok || val != 1 {
		t.Errorf("Expected 1, got %d (ok=%v)", val, ok)
	}
	val, ok = reader.Read()
	if !ok || val != 2 {
		t.Errorf("Expected 2, got %d (ok=%v)", val, ok)
	}
	val, ok = reader.Read()
	if !ok || val != 3 {
		t.Errorf("Expected 3, got %d (ok=%v)", val, ok)
	}

	// No more data
	_, ok = reader.Read()
	if ok {
		t.Error("Expected no more data")
	}
}

func TestBusyPollRingBuffer_WrapAround(t *testing.T) {
	buf := NewBusyPollRingBuffer[int](3)
	reader := buf.Reader()

	// Fill the buffer
	buf.Write(1)
	buf.Write(2)
	buf.Write(3)

	// Write more, causing wrap-around
	buf.Write(4)
	buf.Write(5)

	// Reader starts from 'from' position (4's position) and reads until 'to'
	// After writes: from=1 (where 4 was written), to=2 (next position after 5)
	// Reader should see: 4, then 5
	val, ok := reader.Read()
	if !ok || val != 4 {
		t.Errorf("Expected 4, got %d (ok=%v)", val, ok)
	}

	val, ok = reader.Read()
	if !ok || val != 5 {
		t.Errorf("Expected 5, got %d (ok=%v)", val, ok)
	}

	// No more data
	_, ok = reader.Read()
	if ok {
		t.Error("Expected no more data")
	}
}

func TestBusyPollRingBuffer_WriteMany(t *testing.T) {
	buf := NewBusyPollRingBuffer[int](5)
	reader := buf.Reader()

	// Write multiple values at once
	buf.WriteMany([]int{10, 20, 30})

	// Due to WriteMany implementation: from stays at initial position (0)
	// Values written at positions 0, 1, 2; to=3, from=0
	// Reader reads from position 0 to position 3 (exclusive)
	val, ok := reader.Read()
	if !ok || val != 10 {
		t.Errorf("Expected 10, got %d (ok=%v)", val, ok)
	}
	val, ok = reader.Read()
	if !ok || val != 20 {
		t.Errorf("Expected 20, got %d (ok=%v)", val, ok)
	}
	val, ok = reader.Read()
	if !ok || val != 30 {
		t.Errorf("Expected 30, got %d (ok=%v)", val, ok)
	}

	// No more data
	_, ok = reader.Read()
	if ok {
		t.Error("Expected no more data")
	}
}

func TestBusyPollRingBuffer_WriteManyWrapAround(t *testing.T) {
	buf := NewBusyPollRingBuffer[int](3)

	// Write many values that exceed buffer size
	// Buffer size is 3, writing 5 values: 1,2,3,4,5
	// Positions: 0=1, 1=2, 2=3, 0=4, 1=5
	// Final buffer: [4, 5, 3], from=0, to=2
	buf.WriteMany([]int{1, 2, 3, 4, 5})

	reader := buf.Reader()

	// Reader reads from position 0 to position 2 (exclusive)
	val, ok := reader.Read()
	if !ok || val != 4 {
		t.Errorf("Expected 4, got %d (ok=%v)", val, ok)
	}
	val, ok = reader.Read()
	if !ok || val != 5 {
		t.Errorf("Expected 5, got %d (ok=%v)", val, ok)
	}

	// No more data (to=2, current=2)
	_, ok = reader.Read()
	if ok {
		t.Error("Expected no more data")
	}
}

func TestBusyPollRingBuffer_MultipleReaders(t *testing.T) {
	buf := NewBusyPollRingBuffer[int](5)

	// Create readers before writing
	reader1 := buf.Reader()
	reader2 := buf.Reader()

	// Write data
	buf.Write(100)
	buf.Write(200)

	// Both readers should see the same data
	val1, ok1 := reader1.Read()
	val2, ok2 := reader2.Read()
	if !ok1 || !ok2 || val1 != 100 || val2 != 100 {
		t.Errorf("Expected both readers to see 100, got %d and %d", val1, val2)
	}

	val1, ok1 = reader1.Read()
	val2, ok2 = reader2.Read()
	if !ok1 || !ok2 || val1 != 200 || val2 != 200 {
		t.Errorf("Expected both readers to see 200, got %d and %d", val1, val2)
	}
}

func TestBusyPollRingBuffer_ReaderWrapsAround(t *testing.T) {
	buf := NewBusyPollRingBuffer[int](3)
	reader := buf.Reader()

	// Write values that will cause wrap-around
	buf.Write(1)         // pos 0, from=0, to=1
	buf.Write(2)         // pos 1, from=1, to=2
	_, _ = reader.Read() // Read 1 from pos 0, current=1
	_, _ = reader.Read() // Read 2 from pos 1, current=2

	buf.Write(3) // pos 2, from=2, to=0
	buf.Write(4) // pos 0, from=0, to=1
	buf.Write(5) // pos 1, from=1, to=2

	// After the writes: from=1, to=2
	// Reader current=2, but it missed value at pos 0 (which is 4)
	// Reader can't read because current(2) == to(2)
	// This demonstrates the "readers would miss data if they are too slow" behavior
	_, ok := reader.Read()
	if ok {
		t.Error("Expected no data - reader was too slow and missed data")
	}
}

func TestBusyPollRingBuffer_EmptyRead(t *testing.T) {
	buf := NewBusyPollRingBuffer[int](3)
	reader := buf.Reader()

	// Try to read from empty buffer
	_, ok := reader.Read()
	if ok {
		t.Error("Expected no data from empty buffer")
	}
}

func TestBusyPollRingBuffer_StringType(t *testing.T) {
	buf := NewBusyPollRingBuffer[string](3)
	reader := buf.Reader()

	buf.Write("hello")
	buf.Write("world")

	val, ok := reader.Read()
	if !ok || val != "hello" {
		t.Errorf("Expected 'hello', got '%s' (ok=%v)", val, ok)
	}
	val, ok = reader.Read()
	if !ok || val != "world" {
		t.Errorf("Expected 'world', got '%s' (ok=%v)", val, ok)
	}
}

func TestBusyPollRingBuffer_StructType(t *testing.T) {
	type TestStruct struct {
		ID   int
		Name string
	}

	buf := NewBusyPollRingBuffer[TestStruct](3)
	reader := buf.Reader()

	buf.Write(TestStruct{ID: 1, Name: "first"})
	buf.Write(TestStruct{ID: 2, Name: "second"})

	val, ok := reader.Read()
	if !ok || val.ID != 1 || val.Name != "first" {
		t.Errorf("Expected {1, 'first'}, got {%d, '%s'} (ok=%v)", val.ID, val.Name, ok)
	}
	val, ok = reader.Read()
	if !ok || val.ID != 2 || val.Name != "second" {
		t.Errorf("Expected {2, 'second'}, got {%d, '%s'} (ok=%v)", val.ID, val.Name, ok)
	}
}

func TestBusyPollRingBuffer_Len(t *testing.T) {
	sizes := []int{1, 5, 10, 100}
	for _, size := range sizes {
		buf := NewBusyPollRingBuffer[int](size)
		if buf.Len() != size {
			t.Errorf("Expected length %d, got %d", size, buf.Len())
		}
	}
}

func TestBusyPollRingBuffer_ReaderAfterWrites(t *testing.T) {
	buf := NewBusyPollRingBuffer[int](3)

	// Write data first
	buf.Write(1) // pos 0, from=0, to=1
	buf.Write(2) // pos 1, from=1, to=2
	buf.Write(3) // pos 2, from=2, to=0

	// Create reader after writes - starts at from=2
	reader := buf.Reader()

	buf.Write(4) // pos 0, from=0, to=1

	// Reader is at current=2, can read from pos 2 (value 3) then wrap to pos 0 (value 4)
	val, ok := reader.Read()
	if !ok || val != 3 {
		t.Errorf("Expected 3, got %d (ok=%v)", val, ok)
	}

	val, ok = reader.Read()
	if !ok || val != 4 {
		t.Errorf("Expected 4, got %d (ok=%v)", val, ok)
	}
}

func TestBusyPollRingBuffer_BeginEndWrite(t *testing.T) {
	buf := NewBusyPollRingBuffer[int](3)
	reader := buf.Reader()

	// Use BeginWrite/EndWrite to write a value
	ptr := buf.BeginWrite()
	*ptr = 42
	buf.EndWrite()

	// Read the value back
	val, ok := reader.Read()
	if !ok || val != 42 {
		t.Errorf("Expected 42, got %d (ok=%v)", val, ok)
	}

	// No more data
	_, ok = reader.Read()
	if ok {
		t.Error("Expected no more data")
	}
}

func TestBusyPollRingBuffer_BeginEndWriteMultiple(t *testing.T) {
	buf := NewBusyPollRingBuffer[int](5)
	reader := buf.Reader()

	// Write multiple values using BeginWrite/EndWrite
	for i := 1; i <= 3; i++ {
		ptr := buf.BeginWrite()
		*ptr = i * 10
		buf.EndWrite()
	}

	// Read them back
	expected := []int{10, 20, 30}
	for i, exp := range expected {
		val, ok := reader.Read()
		if !ok || val != exp {
			t.Errorf("Read %d: expected %d, got %d (ok=%v)", i, exp, val, ok)
		}
	}

	// No more data
	_, ok := reader.Read()
	if ok {
		t.Error("Expected no more data")
	}
}

func TestBusyPollRingBuffer_BeginEndWriteWrapAround(t *testing.T) {
	buf := NewBusyPollRingBuffer[int](3)
	reader := buf.Reader()

	// Fill the buffer using BeginWrite/EndWrite
	for i := 1; i <= 3; i++ {
		ptr := buf.BeginWrite()
		*ptr = i
		buf.EndWrite()
	}

	// Write more, causing wrap-around
	ptr := buf.BeginWrite()
	*ptr = 4
	buf.EndWrite()

	ptr = buf.BeginWrite()
	*ptr = 5
	buf.EndWrite()

	// After writes: from=1 (where 4 was written), to=2 (next position after 5)
	// Reader should see: 4, then 5
	val, ok := reader.Read()
	if !ok || val != 4 {
		t.Errorf("Expected 4, got %d (ok=%v)", val, ok)
	}

	val, ok = reader.Read()
	if !ok || val != 5 {
		t.Errorf("Expected 5, got %d (ok=%v)", val, ok)
	}

	// No more data
	_, ok = reader.Read()
	if ok {
		t.Error("Expected no more data")
	}
}

func TestBusyPollRingBuffer_BeginEndWriteStruct(t *testing.T) {
	type TestStruct struct {
		ID   int
		Name string
	}

	buf := NewBusyPollRingBuffer[TestStruct](3)
	reader := buf.Reader()

	// Write using BeginWrite/EndWrite
	ptr := buf.BeginWrite()
	ptr.ID = 100
	ptr.Name = "test"
	buf.EndWrite()

	// Read back
	val, ok := reader.Read()
	if !ok || val.ID != 100 || val.Name != "test" {
		t.Errorf("Expected {100, 'test'}, got {%d, '%s'} (ok=%v)", val.ID, val.Name, ok)
	}
}

func TestBusyPollRingBuffer_BeginWritePointerStability(t *testing.T) {
	buf := NewBusyPollRingBuffer[int](3)

	// Get pointer from BeginWrite
	ptr1 := buf.BeginWrite()
	*ptr1 = 99

	// Before EndWrite, the pointer should point to the same location
	if *ptr1 != 99 {
		t.Errorf("Expected pointer to still contain 99, got %d", *ptr1)
	}

	buf.EndWrite()

	// After EndWrite, the value should be in the buffer
	reader := buf.Reader()
	val, ok := reader.Read()
	if !ok || val != 99 {
		t.Errorf("Expected 99, got %d (ok=%v)", val, ok)
	}
}

func BenchmarkBusyPollRingBuffer_Write(b *testing.B) {
	buf := NewBusyPollRingBuffer[int](1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Write(i)
	}
}

func BenchmarkBusyPollRingBuffer_WriteMany(b *testing.B) {
	buf := NewBusyPollRingBuffer[int](1000)
	data := make([]int, 100)
	for i := range data {
		data[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.WriteMany(data)
	}
}

func BenchmarkBusyPollRingBuffer_Read(b *testing.B) {
	buf := NewBusyPollRingBuffer[int](b.N)
	reader := buf.Reader()

	// Pre-fill buffer
	for i := 0; i < b.N; i++ {
		buf.Write(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader.Read()
	}
}

func BenchmarkBusyPollRingBuffer_WriteRead(b *testing.B) {
	buf := NewBusyPollRingBuffer[int](b.N)
	reader := buf.Reader()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Write(i)
		reader.Read()
	}
}

func BenchmarkBusyPollRingBuffer_BeginEndWrite(b *testing.B) {
	buf := NewBusyPollRingBuffer[int](1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ptr := buf.BeginWrite()
		*ptr = i
		buf.EndWrite()
	}
}
