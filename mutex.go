package powch

import "time"

var Lock = Unit{}

// Mutex this mutex is more powerful than sync.Mutex. it supports use with select, try lock and lock with timeout
type Mutex struct {
	ch chan Unit
}

// NewMutex returns a Mutex
func NewMutex() Mutex {
	return Mutex{
		ch: make(chan Unit, 1),
	}
}

// Lock just like the Lock of sync.Mutex
func (m Mutex) Lock() Mutex {
	m.ch <- Lock
	return m
}

// Unlock just like the Unlock of sync.Mutex
func (m Mutex) Unlock() {
	<-m.ch
}

// TryLock nonblock lock
func (m Mutex) TryLock() bool {
	select {
	case m.ch <- Lock:
		return true
	default:
		return false
	}
}

// LockWithTimeout lock with timeout
func (m Mutex) LockWithTimeout(t time.Duration) bool {
	select {
	case m.ch <- Lock:
		return true
	case <-time.After(t):
		return false
	}
}

// LockChan it can support select with this
func (m Mutex) LockChan() chan<- Unit {
	return m.ch
}
