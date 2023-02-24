package powch

import "time"

type RMutex struct {
	ch chan Unit
}

func NewRMutex(max int) RMutex {
    if max < 1 {
        max = 1
    }
	return RMutex{
		ch: make(chan Unit, max),
	}
}


func (m RMutex) Lock() RMutex {
	m.ch <- Lock
    return m
}

func (m RMutex) Unlock() {
	<-m.ch
}

func (m RMutex) TryLock() bool {
    select {
    case m.ch <- Lock:
        return true
    default:
        return false
    }
}

func (m RMutex) LockWithTimeout(t time.Duration) bool {
    select {
    case m.ch <- Lock:
        return true
    case <-time.After(t):
        return false
    }
}

func (m RMutex) LockChan() chan<- Unit {
	return m.ch
}

