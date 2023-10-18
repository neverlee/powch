package powch

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Dedup(t *testing.T) {
	dd := NewDedup[int, int](10, func(t int) int {
		return t % 8
	})

	for i := 0; i < 12; i++ {
		dd.InChan() <- i
	}
	res := []int{8, 9, 10, 11, 4, 5, 6, 7}
	dd.Close(false)
	i := 0
	for v := range dd.OutChan() {
		assert.Equal(t, res[i], v)
		i++
	}
}

func Test_Dedup2(t *testing.T) {
	dd := NewDedup[int](4, func(a int) int {
		return a % 10
	})
	assert.Equal(t, 0, dd.l.Len())
	dd.InChan() <- 1
	dd.InChan() <- 2
	dd.InChan() <- 3
	dd.InChan() <- 4

	atDefault := false
	select {
	case dd.InChan() <- 1:
	case dd.InChan() <- 1:
	default:
		atDefault = true
	}
	assert.True(t, atDefault)

	assert.Equal(t, 1, <-dd.OutChan())
	assert.Equal(t, 2, <-dd.OutChan())
	assert.Equal(t, 3, <-dd.OutChan())
	assert.Equal(t, 4, <-dd.OutChan())

	assert.False(t, dd.closed)
	dd.Close(false)
	runtime.Gosched()
	assert.True(t, dd.closed)
	assert.True(t, dd.loopEnd)

	assert.Panics(t, func() {
		dd.InChan() <- 1
	})
	assert.Panics(t, func() {
		dd.InChan() <- 1
	})
	{
		v, ok := <-dd.OutChan()
		assert.Equal(t, 0, v)
		assert.False(t, ok)
	}
}
