//go:build go1.23
// +build go1.23

package quantainer

import (
	"reflect"
	"testing"
)

func TestRingBufferValues(t *testing.T) {
	rb := NewRingBuffer[int](3)
	rb.AddLast(1)
	rb.AddLast(2)
	rb.AddLast(3)
	rb.AddLast(4)

	result := []int{}
	for v := range rb.Values() {
		result = append(result, v)
	}
	if !reflect.DeepEqual(result, []int{2, 3, 4}) {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestRingBufferAll(t *testing.T) {
	rb := NewRingBuffer[int](3)
	rb.AddLast(1)
	rb.AddLast(2)
	rb.AddLast(3)
	rb.AddLast(4)

	result := []int{}
	for i, v := range rb.All() {
		if v != i+2 {
			t.Fail()
		}
		result = append(result, v)
	}
	if !reflect.DeepEqual(result, []int{2, 3, 4}) {
		t.Errorf("unexpected result: %v", result)
	}
}
