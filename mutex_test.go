package powch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Mutex_LockUnlock(t *testing.T) {
	m := NewMutex()
	assert.True(t, m.TryLock())
	m.Unlock()

	m.Lock()
	assert.False(t, m.TryLock())
	m.Unlock()
}
