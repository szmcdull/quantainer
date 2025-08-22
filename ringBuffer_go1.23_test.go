//go:build go1.23
// +build go1.23

package quantainer

import (
	"reflect"
	"testing"
)

func TestRingBufferAll(t *testing.T) {
	rb := NewRingBuffer[int](3)
	rb.AddLast(1)
	rb.AddLast(2)
	rb.AddLast(3)
	rb.AddLast(4)

	result := []int{}
	for v := range rb.All() {
		result = append(result, v)
	}
	if !reflect.DeepEqual(result, []int{2, 3, 4}) {
		t.Errorf("unexpected result: %v", result)
	}
}
