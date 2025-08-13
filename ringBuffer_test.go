package quantainer

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

// ExampleRingBuffer_At demonstrates accessing elements by logical index
// from oldest (0) to newest (count-1), including after the buffer wraps.
func ExampleRingBuffer_At() {
	rb := NewRingBuffer[int](3)
	rb.AddLast(1)
	rb.AddLast(2)
	rb.AddLast(3)
	fmt.Println(*rb.At(0), *rb.At(1), *rb.At(2), *rb.At(-1), *rb.At(-2))

	// Adding another element overwrites the oldest (1)
	rb.AddLast(4)
	fmt.Println(*rb.At(0), *rb.At(1), *rb.At(2), *rb.At(-1), *rb.At(-2))

	// Output:
	// 1 2 3 3 2
	// 2 3 4 4 3
}

func TestRingBuffer_At(t *testing.T) {
	rb := NewRingBuffer[int](3)
	rb.AddLast(1)

	v := rb.At(-1)
	if v == nil || *v != 1 {
		t.Fatalf("At(-1) want 1 got %v", *v)
	}
	v = rb.At(0)
	if v == nil || *v != 1 {
		t.Fatalf("At(0) want 1 got %v", *v)
	}

	rb.AddLast(2)
	rb.AddLast(3)
	rb.AddLast(4)

	v = rb.At(-1)
	if v == nil || *v != 4 {
		t.Fatalf("At(-1) want 4 got %v", *v)
	}
	v = rb.At(0)
	if v == nil || *v != 2 {
		t.Fatalf("At(0) want 1 got %v", *v)
	}
}

// helper to compare slices in tests
func equalSlice[T comparable](a, b []T) bool {
	return reflect.DeepEqual(a, b)
}

func TestRingBuffer_Basic(t *testing.T) {
	rb := NewRingBuffer[int](3)
	if rb.Count() != 0 {
		t.Fatalf("expected empty count=0, got %d", rb.Count())
	}

	rb.AddLast(1)
	rb.AddLast(2)
	if rb.Count() != 2 {
		t.Fatalf("count want 2 got %d", rb.Count())
	}
	if got := rb.ToSlice(); !equalSlice(got, []int{1, 2}) {
		t.Fatalf("ToSlice want [1 2] got %v", got)
	}

	first := rb.First()
	if first == nil || *first != 1 {
		t.Fatalf("First want 1 got %v", first)
	}
}

func TestRingBuffer_WrapAndOverwrite(t *testing.T) {
	rb := NewRingBuffer[int](3)
	rb.AddLast(1)
	rb.AddLast(2)
	rb.AddLast(3)
	if got := rb.ToSlice(); !equalSlice(got, []int{1, 2, 3}) {
		t.Fatalf("ToSlice before wrap want [1 2 3] got %v", got)
	}
	// overwrite oldest (1)
	rb.AddLast(4)
	if got := rb.ToSlice(); !equalSlice(got, []int{2, 3, 4}) {
		t.Fatalf("ToSlice after wrap want [2 3 4] got %v", got)
	}
	if rb.Count() != 3 {
		t.Fatalf("count want 3 got %d", rb.Count())
	}
	if first := rb.First(); first == nil || *first != 2 {
		t.Fatalf("First want 2 got %v", first)
	}
}

func TestRingBuffer_ToSliceAndSort(t *testing.T) {
	rb := NewRingBuffer[int](5)
	rb.AddLast(5)
	rb.AddLast(1)
	rb.AddLast(3)
	rb.AddLast(4)
	rb.AddLast(2)
	// validate order
	if got := rb.ToSlice(); !equalSlice(got, []int{5, 1, 3, 4, 2}) {
		t.Fatalf("ToSlice want [5 1 3 4 2] got %v", got)
	}
	// sort into provided slice
	buf := make([]int, rb.Count())
	out := rb.ToSliceAndSort(buf)
	if &buf[0] != &out[0] { // ensure reuse when provided
		t.Fatalf("ToSliceAndSort should reuse provided slice")
	}
	if !equalSlice(out, []int{1, 2, 3, 4, 5}) {
		t.Fatalf("sorted want [1 2 3 4 5] got %v", out)
	}
}

func TestRingBuffer_Clear(t *testing.T) {
	rb := NewRingBuffer[int](2)
	rb.AddLast(9)
	rb.AddLast(8)
	rb.Clear()
	if rb.Count() != 0 {
		t.Fatalf("count after Clear want 0 got %d", rb.Count())
	}
	// First on empty returns zero-value
	if first := rb.First(); first != nil {
		t.Fatalf("First on empty want nil got %v", first)
	}
}

// Tests below validate dynamic resizing behavior via NewRingBufferConfigurable.
// If you never change maxSize at runtime, you can skip these.

func TestRingBuffer_Grow_NoWrap(t *testing.T) {
	size := 3
	rb := NewRingBufferConfigurable[int](func() int { return size })
	rb.AddLast(1)
	rb.AddLast(2)
	size = 5 // grow before wrap
	rb.AddLast(3)
	rb.AddLast(4)
	rb.AddLast(5)
	if got := rb.ToSlice(); !equalSlice(got, []int{1, 2, 3, 4, 5}) {
		t.Fatalf("grow no wrap want [1 2 3 4 5] got %v", got)
	}
}

func TestRingBuffer_Grow_WithWrap(t *testing.T) {
	size := 3
	rb := NewRingBufferConfigurable[int](func() int { return size })
	rb.AddLast(1)
	rb.AddLast(2)
	rb.AddLast(3)
	rb.AddLast(4) // wraps, buffer now [2,3,4]
	size = 6      // grow after wrap
	rb.AddLast(5)
	// Expected logical sequence is [2,3,4,5]
	if got := rb.ToSlice(); !equalSlice(got, []int{2, 3, 4, 5}) {
		t.Fatalf("grow with wrap want [2 3 4 5] got %v", got)
	}
}

func TestRingBuffer_Shrink_FromFull(t *testing.T) {
	size := 5
	rb := NewRingBufferConfigurable[int](func() int { return size })
	rb.AddLast(1)
	rb.AddLast(2)
	rb.AddLast(3)
	rb.AddLast(4)
	rb.AddLast(5) // full [1..5]
	size = 3      // shrink to 3
	// After shrink, the logical window should contain the most recent 3 when next elements come in.
	rb.AddLast(6)
	// Expect tail behavior to keep last 3 elements in order
	if got := rb.ToSlice(); !equalSlice(got, []int{4, 5, 6}) {
		t.Fatalf("shrink from full want [4 5 6] got %v", got)
	}
}

func TestRingBuffer_Shrink_WhenSparse(t *testing.T) {
	size := 5
	rb := NewRingBufferConfigurable[int](func() int { return size })
	rb.AddLast(1)
	rb.AddLast(2) // count=2
	size = 3      // shrink, still enough capacity
	if got := rb.ToSlice(); !equalSlice(got, []int{1, 2}) {
		t.Fatalf("shrink sparse keep order want [1 2] got %v", got)
	}
	rb.AddLast(3)
	if got := rb.ToSlice(); !equalSlice(got, []int{1, 2, 3}) {
		t.Fatalf("shrink sparse then add want [1 2 3] got %v", got)
	}
}

// Optional: zero-size behavior. If size==0 is not a supported use-case, panic is acceptable.
func TestRingBuffer_SizeZero_AddPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic on AddLast with size 0")
		}
	}()
	rb := NewRingBuffer[int](0)
	rb.AddLast(1)
}

func TestRingBuffer_Grow_WithWrap_TailLESizeDiff(t *testing.T) {
	size := 5
	rb := NewRingBufferConfigurable[int](func() int { return size })
	// Fill and wrap once so tail=1
	rb.AddLast(1)
	rb.AddLast(2)
	rb.AddLast(3)
	rb.AddLast(4)
	rb.AddLast(5)
	rb.AddLast(6) // now logically [2,3,4,5,6]

	// Grow by 2 so sizeDiff=2 and current tail==1 (<= sizeDiff)
	size = 7
	rb.AddLast(7)

	got := rb.ToSlice()
	want := []int{2, 3, 4, 5, 6, 7}
	if len(got) != len(want) {
		t.Fatalf("len got %d want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("grow wrap tail<=diff: idx %d got %v want %v; slice=%v", i, got[i], want[i], got)
		}
	}
}

func TestRingBuffer_Grow_WithWrap_TailEqSizeDiff(t *testing.T) {
	size := 5
	rb := NewRingBufferConfigurable[int](func() int { return size })
	// Make tail==2
	rb.AddLast(1)
	rb.AddLast(2)
	rb.AddLast(3)
	rb.AddLast(4)
	rb.AddLast(5)
	rb.AddLast(6) // tail=1
	rb.AddLast(7) // tail=2

	// Grow by 2 so sizeDiff=2 and tail==sizeDiff
	size = 7
	rb.AddLast(8)

	got := rb.ToSlice()
	want := []int{3, 4, 5, 6, 7, 8}
	if len(got) != len(want) {
		t.Fatalf("len got %d want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("grow wrap tail==diff: idx %d got %v want %v; slice=%v", i, got[i], want[i], got)
		}
	}
}

func TestRingBuffer_Grow_WithWrap_TailGtSizeDiff(t *testing.T) {
	size := 5
	rb := NewRingBufferConfigurable[int](func() int { return size })
	// Make tail==3 (> sizeDiff=2)
	rb.AddLast(1)
	rb.AddLast(2)
	rb.AddLast(3)
	rb.AddLast(4)
	rb.AddLast(5)
	rb.AddLast(6) // tail=1
	rb.AddLast(7) // tail=2
	rb.AddLast(8) // tail=3

	size = 7 // sizeDiff=2
	rb.AddLast(9)

	got := rb.ToSlice()
	want := []int{4, 5, 6, 7, 8, 9}
	if len(got) != len(want) {
		t.Fatalf("len got %d want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("grow wrap tail>diff: idx %d got %v want %v; slice=%v", i, got[i], want[i], got)
		}
	}
}

func TestRingBuffer_LongWraps(t *testing.T) {
	rb := NewRingBuffer[int](3)
	for i := 1; i <= 10; i++ {
		rb.AddLast(i)
	}
	got := rb.ToSlice()
	want := []int{8, 9, 10}
	if len(got) != 3 || got[0] != want[0] || got[1] != want[1] || got[2] != want[2] {
		t.Fatalf("long wraps want %v got %v", want, got)
	}
	if !rb.Full() {
		t.Fatalf("expected Full() after many wraps")
	}
}

func TestRingBuffer_FullFlag(t *testing.T) {
	rb := NewRingBuffer[int](3)
	if rb.Full() {
		t.Fatalf("Full() should be false initially")
	}
	rb.AddLast(1)
	rb.AddLast(2)
	rb.AddLast(3)
	if !rb.Full() {
		t.Fatalf("Full() should be true when at capacity")
	}
	rb.AddLast(4)
	if !rb.Full() || rb.Count() != 3 {
		t.Fatalf("Full() should remain true and count stay at max; got full=%v count=%d", rb.Full(), rb.Count())
	}
}

func TestRingBuffer_FloatNaN_Sort(t *testing.T) {
	rb := NewRingBuffer[float64](5)
	rb.AddLast(5)
	rb.AddLast(math.NaN())
	rb.AddLast(3)
	rb.AddLast(1)
	rb.AddLast(2)

	out := rb.ToSliceAndSortNaN(nil)
	if len(out) != 5 {
		t.Fatalf("len want 5 got %d", len(out))
	}
	// Validate: one NaN present; non-NaNs are sorted ascending
	nanCount := 0
	nums := make([]float64, 0, 5)
	for _, v := range out {
		if math.IsNaN(v) {
			nanCount++
		} else {
			nums = append(nums, v)
		}
	}
	if nanCount != 1 {
		t.Fatalf("expected exactly 1 NaN, got %d (out=%v)", nanCount, out)
	}
	for i := 1; i < len(nums); i++ {
		if nums[i-1] > nums[i] {
			t.Fatalf("non-NaN segment not sorted: %v", nums)
		}
	}
}

func TestRingBuffer_Shrink_KeepLatest(t *testing.T) {
	size := 5
	rb := NewRingBufferConfigurable[int](func() int { return size })
	rb.AddLast(1)
	rb.AddLast(2)
	rb.AddLast(3)
	rb.AddLast(4)
	size = 3
	rb.AddLast(5)
	got := rb.ToSlice()
	want := []int{3, 4, 5}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("shrink then add: idx %d got %v want %v; slice=%v", i, got[i], want[i], got)
		}
	}
	rb.AddLast(6)
	got = rb.ToSlice()
	want = []int{4, 5, 6}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("shrink then add2: idx %d got %v want %v; slice=%v", i, got[i], want[i], got)
		}
	}
}
