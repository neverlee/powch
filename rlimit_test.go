package powch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Rlimit(t *testing.T) {
	rl := NewRlimit(2)
	rl.Enter()
	assert.True(t, rl.TryEnter())
	assert.False(t, rl.TryEnter())
	rl.Leave()
	rl.Leave()
}
