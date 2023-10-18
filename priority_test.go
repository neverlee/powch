package powch

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Example_Prioirty() {
	// create a pub channel
	pri := NewPriority[int](10, func(a, b int) bool {
		return a < b
	})

	ich := pri.InChan()
	indata := []int{1, 2, 2, 4, 1, 3, 5, 7}
	for _, d := range indata {
		ich <- d
	}
	pri.Close(false)

	for d := range pri.OutChan() {
		fmt.Print(d, ",")
	}
	// Output:
	// 1,1,2,2,3,4,5,7,
}

func Test_Priority(t *testing.T) {
	pq := NewPriority[int](10, func(a, b int) bool {
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

func Test_Priority2(t *testing.T) {
	pq := NewPriority[int](4, func(a, b int) bool {
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
