package quantainer

import (
	"fmt"
	"time"
)

func ExampleFixedDurationSlice_Add() {
	l := NewFixedDurationSlice[int](time.Second * 2)
	l.Add(1)
	time.Sleep(time.Second)
	l.Add(2)
	time.Sleep(time.Second)
	l.Add(3)
	for _, v := range l.Values() {
		fmt.Println(v)
	}
	// Output:
	// 2
	// 3
}
