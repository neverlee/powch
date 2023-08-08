package heap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Heap(t *testing.T) {
	h := NewHeap[int](func(a, b int) bool {
		return a < b
	})

	for i := 10; i > 0; i-- {
		h.Push(i)
	}
	for i := 1; i <= 10; i++ {
		assert.Equal(t, i, h.Pop())
	}
}
