package quantainer

import (
	"fmt"
	"reflect"
	"testing"
)

func ExampleSortedRingBuffer_AddLast() {
	l := NewSortedRingBuffer[int](3)
	l.AddLast(1)
	l.AddLast(2)
	l.AddLast(3)
	l.AddLast(4)
	for _, v := range l.ToSlice() {
		fmt.Println(v)
	}
	// Output:
	// 2
	// 3
	// 4
}

func ExampleSortedRingBuffer_SortedSlice() {
	l := NewSortedRingBuffer[int](3)
	l.AddLast(4)
	l.AddLast(3)
	l.AddLast(2)
	l.AddLast(1)
	for _, v := range l.SortedSlice() {
		fmt.Println(v)
	}
	// Output:
	// 1
	// 2
	// 3
}

// Test duplicates handling and TreeMap count sync when items are evicted
func TestSortedRingBuffer_DuplicatesAndEviction(t *testing.T) {
	l := NewSortedRingBuffer[int](5)
	l.AddLast(2)
	l.AddLast(1)
	l.AddLast(2)
	l.AddLast(3)
	l.AddLast(2) // buffer is full: [2,1,2,3,2]

	got := l.SortedSlice()
	want := []int{1, 2, 2, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("sorted before eviction: got %v want %v", got, want)
	}

	// Adding another element should evict the oldest (2) and keep counts correct
	l.AddLast(4) // buffer becomes logically [1,2,3,2,4]
	got = l.SortedSlice()
	want = []int{1, 2, 2, 3, 4}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("sorted after eviction: got %v want %v", got, want)
	}
}

// Test First/Last pointers across wrap/eviction
func TestSortedRingBuffer_FirstLast(t *testing.T) {
	l := NewSortedRingBuffer[int](3)
	l.AddLast(10)
	l.AddLast(20)
	if f := l.First(); f == nil || *f != 10 {
		t.Fatalf("First want 10 got %v", f)
	}
	if lst := l.Last(); lst == nil || *lst != 20 {
		t.Fatalf("Last want 20 got %v", lst)
	}

	l.AddLast(30)
	if f := l.First(); f == nil || *f != 10 {
		t.Fatalf("First want 10 got %v", f)
	}
	if lst := l.Last(); lst == nil || *lst != 30 {
		t.Fatalf("Last want 30 got %v", lst)
	}

	// Evict 10
	l.AddLast(40)
	if f := l.First(); f == nil || *f != 20 {
		t.Fatalf("First after eviction want 20 got %v", f)
	}
	if lst := l.Last(); lst == nil || *lst != 40 {
		t.Fatalf("Last after eviction want 40 got %v", lst)
	}
}

// PopFirst should return nil on empty buffer
func TestSortedRingBuffer_PopFirst_Empty(t *testing.T) {
	l := NewSortedRingBuffer[int](2)
	if v := l.PopFirst(); v != nil {
		t.Fatalf("PopFirst on empty want nil got %v", v)
	}
}

// Clear should reset both underlying ring buffer and the sorted view
func TestSortedRingBuffer_Clear(t *testing.T) {
	l := NewSortedRingBuffer[int](3)
	l.AddLast(3)
	l.AddLast(1)
	l.AddLast(2)
	l.Clear()
	if l.Count() != 0 {
		t.Fatalf("Count after Clear want 0 got %d", l.Count())
	}
	if s := l.ToSlice(); len(s) != 0 {
		t.Fatalf("ToSlice after Clear want empty got %v", s)
	}
	if s := l.SortedSlice(); len(s) != 0 {
		t.Fatalf("SortedSlice after Clear want empty got %v", s)
	}
	if l.First() != nil || l.Last() != nil {
		t.Fatalf("First/Last after Clear should be nil")
	}
	if l.Full() {
		t.Fatalf("Full() should be false after Clear")
	}
}

// Full and MaxSize behavior
func TestSortedRingBuffer_FullAndMaxSize(t *testing.T) {
	l := NewSortedRingBuffer[int](2)
	if l.MaxSize() != 2 {
		t.Fatalf("MaxSize want 2 got %d", l.MaxSize())
	}
	if l.Full() {
		t.Fatalf("Full() should be false initially")
	}
	l.AddLast(1)
	l.AddLast(2)
	if !l.Full() || l.Count() != 2 {
		t.Fatalf("Full() true and Count==2 expected; got full=%v count=%d", l.Full(), l.Count())
	}
	l.AddLast(3) // evict 1
	if !l.Full() || l.Count() != 2 {
		t.Fatalf("Full() should remain true and Count==2 after overwrite; got full=%v count=%d", l.Full(), l.Count())
	}
	if got, want := l.ToSlice(), []int{2, 3}; !reflect.DeepEqual(got, want) {
		t.Fatalf("ToSlice after overwrite got %v want %v", got, want)
	}
	if got, want := l.SortedSlice(), []int{2, 3}; !reflect.DeepEqual(got, want) {
		t.Fatalf("SortedSlice after overwrite got %v want %v", got, want)
	}
}

// For size==0, underlying ring buffer panics on AddLast; SortedRingBuffer should propagate
func TestSortedRingBuffer_SizeZero_AddPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic on AddLast with size 0")
		}
	}()
	l := NewSortedRingBuffer[int](0)
	l.AddLast(1)
}
