package powch

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PriorityQueue(t *testing.T) {
	pq := NewPriorityQueue[int](10, func(a, b int) bool {
		return a > b
	})

	for i := 0; i < 10; i++ {
		pq.InChan() <- i
	}
	for i := 0; i < 10; i++ {
		v, ok := <-pq.OutChan()
		assert.True(t, ok)
		assert.Equal(t, 9-i, v)
	}
	pq.Close(true)
}

func Test_PriorityQueue2(t *testing.T) {
	pq := NewPriorityQueue[int](4, func(a, b int) bool {
		return a > b
	})
	assert.Equal(t, 0, pq.h.Len())
	pq.InChan() <- 1
	pq.InChan() <- 2
	pq.InChan() <- 3
	pq.InChan() <- 4

	atDefault := false
	select {
	case pq.InChan() <- 1:
	case pq.InChan() <- 1:
	default:
		atDefault = true
	}
	assert.True(t, atDefault)

	assert.Equal(t, 4, <-pq.OutChan())
	assert.Equal(t, 3, <-pq.OutChan())
	assert.Equal(t, 2, <-pq.OutChan())
	assert.Equal(t, 1, <-pq.OutChan())

	assert.False(t, pq.closed)
	pq.Close(false)
	runtime.Gosched()
	assert.True(t, pq.closed)
	assert.True(t, pq.loopEnd)

	assert.Panics(t, func() {
		pq.InChan() <- 1
	})
	assert.Panics(t, func() {
		pq.InChan() <- 1
	})
	{
		v, ok := <-pq.OutChan()
		assert.Equal(t, 0, v)
		assert.False(t, ok)
	}
}
