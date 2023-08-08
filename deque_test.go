package powch

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Deque(t *testing.T) {
	dq := NewDeque[int](4)
	assert.Equal(t, 0, dq.lst.Len())
	dq.BinChan() <- 1
	dq.BinChan() <- 2
	dq.BinChan() <- 3
	dq.BinChan() <- 4

	atDefault := false
	select {
	case dq.BinChan() <- 1:
	case dq.FinChan() <- 1:
	default:
		atDefault = true
	}
	assert.True(t, atDefault)

	assert.Equal(t, 4, <-dq.BoutChan())
	assert.Equal(t, 3, <-dq.BoutChan())
	assert.Equal(t, 1, <-dq.FoutChan())
	assert.Equal(t, 2, <-dq.FoutChan())

	atDefault = false
	select {
	case <-dq.BoutChan():
	case <-dq.FoutChan():
	default:
		atDefault = true
	}
	assert.True(t, atDefault)

	dq.FinChan() <- 1
	dq.FinChan() <- 2
	dq.FinChan() <- 3
	dq.FinChan() <- 4
	assert.Equal(t, 4, <-dq.FoutChan())
	assert.Equal(t, 3, <-dq.FoutChan())
	assert.Equal(t, 1, <-dq.BoutChan())
	assert.Equal(t, 2, <-dq.BoutChan())

	assert.False(t, dq.closed)
	dq.Close(false)
	runtime.Gosched()
	assert.True(t, dq.closed)
	assert.True(t, dq.loopEnd)

	assert.Panics(t, func() {
		dq.BinChan() <- 1
	})
	assert.Panics(t, func() {
		dq.FinChan() <- 1
	})
	{
		v, ok := <-dq.FoutChan()
		assert.Equal(t, 0, v)
		assert.False(t, ok)
	}
	{
		v, ok := <-dq.BoutChan()
		assert.Equal(t, 0, v)
		assert.False(t, ok)
	}
}
