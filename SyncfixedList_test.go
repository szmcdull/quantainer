package quantainer

import "fmt"

func ExampleSyncFixedList_AddLast() {
	l := NewSyncFixedList[int](3)
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

func ExampleSyncFixedList_AddFirst() {
	l := NewSyncFixedList[int](3)
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
