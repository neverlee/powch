package mutex

import (
	"sync"
)

type ExMutex struct {
	// disable atomic.Bool
	disable bool
	mu      sync.Mutex
}

func (m *ExMutex) Lock() {
	if !m.disable {
		m.mu.Lock()
	}
}

func (m *ExMutex) Unlock() {
	if !m.disable {
		m.mu.Unlock()
	}
}

func (m *ExMutex) Apply(fn func()) {
	m.Lock()
	defer m.Unlock()

	m.disable = true
	fn()
	m.disable = false
}
