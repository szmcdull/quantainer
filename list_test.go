package quantainer

import (
	"fmt"
)

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
