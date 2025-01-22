package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWMap(t *testing.T) {
	{
		wm := NewWMap[string, int](10)
		assert.Equal(t, 0, wm.Len())
	}

	{
		wm := ToWMap(map[int]int{1: 11, 2: 22, 4: 44})
		assert.Equal(t, 3, wm.Len())

		{
			v, ok := wm.Get(1)
			assert.True(t, ok)
			assert.Equal(t, 11, v)
		}

		{
			v, ok := wm.Get(3)
			assert.False(t, ok)
			assert.Equal(t, 0, v)
		}

		wm.Set(3, 33)
		{
			v, ok := wm.Get(3)
			assert.True(t, ok)
			assert.Equal(t, 33, v)
		}

		wm.Delete(3)
		{
			v, ok := wm.Get(3)
			assert.False(t, ok)
			assert.Equal(t, 0, v)
		}

		{
			do := 2
			wm.Range(func(k, v int) bool {
				if do > 0 {
					do--
					return true
				}
				return false
			})
			assert.Equal(t, 0, do)
		}

	}
}
