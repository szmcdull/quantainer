package quantainer

import (
	"fmt"
	"testing"
)

func ExampleFixedList_AddLast() {
	l := NewFixedList[int](3)
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

func ExampleFixedList_AddFirst() {
	l := NewFixedList[int](3)
	l.AddFirst(1)
	l.AddFirst(2)
	l.AddFirst(3)
	l.AddFirst(4)
	for _, v := range l.ToSlice() {
		fmt.Println(v)
	}
	// Output:
	// 4
	// 3
	// 2
}

// Filled is Full without the size-0 degeneracy: empty and size-0 are not filled.
func TestFixedList_Filled(t *testing.T) {
	zero := NewFixedList[int](0)
	if !zero.Full() {
		t.Fatalf("size 0 Full want true")
	}
	if zero.Filled() {
		t.Fatalf("size 0 Filled want false")
	}

	l := NewFixedList[int](2)
	if l.Filled() {
		t.Fatalf("empty Filled want false")
	}
	l.AddLast(1)
	if l.Filled() {
		t.Fatalf("partial Filled want false")
	}
	l.AddLast(2)
	if !l.Filled() {
		t.Fatalf("at capacity Filled want true")
	}
}
