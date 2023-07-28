package powch

import "time"

// Rlimit a goroutine limiter
type Rlimit struct {
	ch chan Unit
}

// returns a created Rlimit
func NewRlimit(max int) Rlimit {
	if max < 1 {
		max = 1
	}
	return Rlimit{
		ch: make(chan Unit, max),
	}
}

// Enter block until enter
func (m Rlimit) Enter() Rlimit {
	m.ch <- Lock
	return m
}

// Leave block until leave
func (m Rlimit) Leave() {
	<-m.ch
}

// TryEnter nonblock, returns true if success in entering otherwise false
func (m Rlimit) TryEnter() bool {
	select {
	case m.ch <- Lock:
		return true
	default:
		return false
	}
}

func (m Rlimit) EnterWithTimeout(t time.Duration) bool {
	select {
	case m.ch <- Lock:
		return true
	case <-time.After(t):
		return false
	}
}

func (m Rlimit) EnterChan() chan<- Unit {
	return m.ch
}
