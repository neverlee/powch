package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleSyncWMap() {
	swmap := NewSyncWMap[string, int]()

	swmap.Set("10", 10)

	swmap.Do(func() {
		swmap.Set("9", 9)
		swmap.Delete("10", "9")
		swmap.Get("10")
	})
}

func TestSyncWMap(t *testing.T) {
	swmap := NewSyncWMap[string, int]()

	{
		v, ok := swmap.Get("10")
		assert.Equal(t, 0, v)
		assert.False(t, ok)
	}
	swmap.Set("10", 10)
	{
		v, ok := swmap.Get("10")
		assert.Equal(t, 10, v)
		assert.True(t, ok)
	}

	swmap.Do(func() {
		swmap.Set("9", 9)
		swmap.Delete("10")
	})

	{
		v, ok := swmap.Get("9")
		assert.Equal(t, 9, v)
		assert.True(t, ok)
	}
}
