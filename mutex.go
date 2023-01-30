package powch

import "time"

var Lock = Unit{}

type Mutex struct {
	ch chan Unit
}

func NewMutex() Mutex {
	return Mutex{
		ch: make(chan Unit, 1),
	}
}


func (m Mutex) Lock() Mutex {
	m.ch <- Lock
    return m
}

func (m Mutex) Unlock() {
	<-m.ch
}

func (m Mutex) TryLock() bool {
    select {
    case m.ch <- Lock:
        return true
    default:
        return false
    }
}

func (m Mutex) LockWithTimeout(t time.Duration) bool {
    select {
    case m.ch <- Lock:
        return true
    case <-time.After(t):
        return false
    }
}

func (m Mutex) LockChan() chan<- Unit {
	return m.ch
}

