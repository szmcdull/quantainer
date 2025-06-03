package quantainer

import "fmt"

func ExampleSortedFixedList_AddLast() {
	l := NewSortedFixedList[int](3)
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

func ExampleSortedFixedList_SortedSlice() {
	l := NewSortedFixedList[int](3)
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

func ExampleSortedFixedList_AddFirst() {
	l := NewSortedFixedList[int](3)
	l.AddFirst(4)
	l.AddFirst(3)
	l.AddFirst(2)
	l.AddFirst(1)
	for _, v := range l.SortedSlice() {
		fmt.Println(v)
	}
	// Output:
	// 1
	// 2
	// 3
}
