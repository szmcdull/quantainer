package quantainer

import (
	"fmt"
	"testing"
)

func TestPopFirst(t *testing.T) {
	l := NewList[int]()
	l.AddFirst(1)
	l.PopFirst()
	if l.First() != nil {
		t.Fail()
	}
	if l.Last() != nil || l.Count() != 0 {
		t.Fail()
	}

	l = NewList[int]()
	l.AddLast(1)
	l.AddLast(2)
	n := l.PopFirst()
	if n.Value != 1 {
		t.Errorf("Expected 1, got %v", n.Value)
	}

	n = l.First()
	if n.Value != 2 {
		t.Fail()
	}
	if n.prev != nil || n.next != nil || l.Count() != 1 {
		t.Errorf(`prev or next is not nil`)
	}
}

func TestPopLast(t *testing.T) {
	l := NewList[int]()
	l.AddFirst(1)
	l.PopLast()
	if l.First() != nil {
		t.Fail()
	}
	if l.Last() != nil || l.Count() != 0 {
		t.Fail()
	}

	l = NewList[int]()
	l.AddLast(1)
	l.AddLast(2)
	n := l.PopLast()
	if n.Value != 2 {
		t.Errorf("Expected 1, got %v", n.Value)
	}

	n = l.First()
	if n.Value != 1 {
		t.Fail()
	}
	if n.prev != nil || n.next != nil || l.Count() != 1 {
		t.Errorf(`prev or next is not nil`)
	}
}

func ExampleList_AddFirst() {
	l := NewList[int]()
	l.AddLast(1)
	l.AddFirst(3)
	l.AddLast(2)
	l.AddFirst(4)
	for n := l.First(); n != nil; n = n.Next() {
		fmt.Println(n.Value)
	}
	// Output:
	// 4
	// 3
	// 1
	// 2
}

func ExampleList_AddLast() {
	l := NewList[int]()
	l.AddFirst(3)
	l.AddLast(1)
	l.AddFirst(4)
	l.AddLast(2)
	for n := l.First(); n != nil; n = n.Next() {
		fmt.Println(n.Value)
	}
	// Output:
	// 4
	// 3
	// 1
	// 2
}

func ExampleList_AddAfter() {
	l := NewList[int]()
	l.AddFirst(3)
	l.AddLast(1)
	l.AddFirst(4)
	n := l.AddLast(2)
	l.AddAfter(n, 5)
	for n := l.First(); n != nil; n = n.Next() {
		fmt.Println(n.Value)
	}
	// Output:
	// 4
	// 3
	// 1
	// 2
	// 5
}

func ExampleList_AddBefore() {
	l := NewList[int]()
	l.AddFirst(3)
	l.AddLast(1)
	n := l.AddFirst(4)
	l.AddLast(2)
	l.AddBefore(n, 5)
	for n := l.First(); n != nil; n = n.Next() {
		fmt.Println(n.Value)
	}
	// Output:
	// 5
	// 4
	// 3
	// 1
	// 2
}

func ExampleList_Remove() {
	l := NewList[int]()
	l.AddFirst(3)
	l.AddLast(1)
	n := l.AddFirst(4)
	l.AddLast(2)
	l.Remove(n)
	for n := l.First(); n != nil; n = n.Next() {
		fmt.Println(n.Value)
	}
	fmt.Println(``)
	fmt.Println(l.First().Value)
	fmt.Println(l.Last().Value)
	// Output:
	// 3
	// 1
	// 2
	//
	// 3
	// 2
}

func ExampleList_Remove_second() {
	l := NewList[int]()
	n := l.AddLast(1)
	l.AddLast(2)
	l.AddLast(3)
	l.AddLast(4)
	l.Remove(n)
	for n := l.First(); n != nil; n = n.Next() {
		fmt.Println(n.Value)
	}
	fmt.Println(``)
	fmt.Println(l.First().Value)
	fmt.Println(l.Last().Value)
	// Output:
	// 2
	// 3
	// 4
	//
	// 2
	// 4
}

func ExampleList_Remove_third() {
	l := NewList[int]()
	l.AddLast(1)
	l.AddLast(2)
	l.AddLast(3)
	n := l.AddLast(4)
	l.Remove(n)
	for n := l.First(); n != nil; n = n.Next() {
		fmt.Println(n.Value)
	}
	fmt.Println(``)
	fmt.Println(l.First().Value)
	fmt.Println(l.Last().Value)
	// Output:
	// 1
	// 2
	// 3
	//
	// 1
	// 3
}

func ExampleList_PopFirstWhen() {
	l := NewList[int]()
	l.AddLast(1)
	l.AddLast(2)
	l.AddLast(3)
	l.AddLast(4)
	l.PopFirstWhen(func(v *int) bool {
		return *v < 3
	})
	for n := l.First(); n != nil; n = n.Next() {
		fmt.Println(n.Value)
	}
	// Output:
	// 3
	// 4
}

func TestPopFirstWhen(t *testing.T) {
	l := NewList[int]()
	l.AddLast(1)
	l.AddLast(2)
	l.AddLast(3)
	l.AddLast(4)
	l.PopFirstWhen(func(v *int) bool {
		return *v < 3
	})
	for n := l.First(); n != nil; n = n.Next() {
		fmt.Println(n.Value)
	}
	if l.Count() != 2 || l.First().prev != nil {
		t.Fail()
	}
}

func ExampleList_PopLastWhen() {
	l := NewList[int]()
	l.AddLast(1)
	l.AddLast(2)
	l.AddLast(3)
	l.AddLast(4)
	l.PopLastWhen(func(v *int) bool {
		return *v > 2
	})
	for n := l.First(); n != nil; n = n.Next() {
		fmt.Println(n.Value)
	}
	// Output:
	// 1
	// 2
}

func TestPopLastWhen(t *testing.T) {
	l := NewList[int]()
	l.AddLast(1)
	l.AddLast(2)
	l.AddLast(3)
	l.AddLast(4)
	l.PopLastWhen(func(v *int) bool {
		return *v > 2
	})
	for n := l.First(); n != nil; n = n.Next() {
		fmt.Println(n.Value)
	}
	if l.Count() != 2 || l.First().prev != nil {
		t.Fail()
	}
}

func TestPopFirstWhen2(t *testing.T) {
	l := NewList[int]()
	l.PopFirstWhen(func(v *int) bool { return true })
}

func TestPopFirstWhen3(t *testing.T) {
	l := NewList[int]()
	l.AddLast(1)
	l.AddLast(2)
	l.AddLast(3)
	l.AddLast(4)
	l.PopFirstWhen(func(v *int) bool { return true })
	if l.Count() != 0 || l._Count() != 0 {
		t.Fail()
	}
	if l.back != nil || l.front != nil {
		t.Fail()
	}
}

func TestPopLastWhen3(t *testing.T) {
	l := NewList[int]()
	l.AddLast(1)
	l.AddLast(2)
	l.AddLast(3)
	l.AddLast(4)
	l.PopLastWhen(func(v *int) bool { return true })
	if l.Count() != 0 || l._Count() != 0 {
		t.Fail()
	}
	if l.back != nil || l.front != nil {
		t.Fail()
	}
}

func TestAt(t *testing.T) {
	l := NewList[int]()
	l.AddLast(1)
	l.AddLast(2)
	l.AddLast(3)
	l.AddLast(4)

	// Test valid indices
	n := l.at(0)
	if n.Value != 1 {
		t.Errorf("Expected 1, got %v", n.Value)
	}

	n = l.at(1)
	if n.Value != 2 {
		t.Errorf("Expected 2, got %v", n.Value)
	}

	n = l.at(2)
	if n.Value != 3 {
		t.Errorf("Expected 3, got %v", n.Value)
	}

	n = l.at(3)
	if n.Value != 4 {
		t.Errorf("Expected 4, got %v", n.Value)
	}

	// Test invalid indices
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for index out of range")
		}
	}()
	l.at(4) // This should panic
}

func TestFromBack(t *testing.T) {
	l := NewList[int]()
	l.AddLast(1)
	l.AddLast(2)
	l.AddLast(3)
	l.AddLast(4)

	// Test valid indices
	n := l.fromBack(0)
	if n.Value != 4 {
		t.Errorf("Expected 4, got %v", n.Value)
	}

	n = l.fromBack(1)
	if n.Value != 3 {
		t.Errorf("Expected 3, got %v", n.Value)
	}

	n = l.fromBack(2)
	if n.Value != 2 {
		t.Errorf("Expected 2, got %v", n.Value)
	}

	n = l.fromBack(3)
	if n.Value != 1 {
		t.Errorf("Expected 1, got %v", n.Value)
	}

	// Test invalid indices
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for index out of range")
		}
	}()
	l.fromBack(4) // This should panic
}

func TestTrim(t *testing.T) {
	l := NewList[int]()
	l.AddLast(0)
	l.AddLast(1)
	l.AddLast(2)
	l.AddLast(3)
	l.AddLast(4)

	// Test trimming from the start
	l.Trim(0, 3)
	if l.Count() != 3 || l.First().Value != 0 || l.Last().Value != 2 {
		t.Errorf("Expected list [0, 1, 2], got %v", l.ToSlice())
	}

	// Reset list
	l = NewList[int]()
	l.AddLast(0)
	l.AddLast(1)
	l.AddLast(2)
	l.AddLast(3)
	l.AddLast(4)

	// Test trimming from the middle
	l.Trim(1, 3)
	if l.Count() != 2 || l.First().Value != 1 || l.Last().Value != 2 {
		t.Errorf("Expected list [1, 2], got %v", l.ToSlice())
	}

	// Reset list
	l = NewList[int]()
	l.AddLast(0)
	l.AddLast(1)
	l.AddLast(2)
	l.AddLast(3)
	l.AddLast(4)

	// Test trimming from the end
	l.Trim(2, 5)
	if l.Count() != 3 || l.First().Value != 2 || l.Last().Value != 4 {
		t.Errorf("Expected list [2, 3, 4], got %v", l.ToSlice())
	}

	// Reset list
	l = NewList[int]()
	l.AddLast(0)
	l.AddLast(1)
	l.AddLast(2)
	l.AddLast(3)
	l.AddLast(4)

	// Test trimming with negative indices
	l.Trim(-3, -1)
	if l.Count() != 2 || l.First().Value != 2 || l.Last().Value != 3 {
		t.Errorf("Expected list [2, 3], got %v", l.ToSlice())
	}

	// Reset list
	l = NewList[int]()
	l.AddLast(0)
	l.AddLast(1)
	l.AddLast(2)
	l.AddLast(3)
	l.AddLast(4)

	// Test trimming the entire list
	l.Trim(0, 0)
	if l.Count() != 0 || l.First() != nil || l.Last() != nil {
		t.Errorf("Expected empty list, got %v", l.ToSlice())
	}

	// Test invalid indices
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for invalid indices")
		}
	}()
	l.Trim(5, 6) // This should panic
}
