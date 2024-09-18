package powch

import "time"

// KIMutex a hash key mutex
type KIMutex struct {
	locks  []Mutex
	count  uint
	handle HashHandle
}

// New return a keymutex
// It require the number of mutexs(prime number is better)
func New(count uint) *KIMutex {
	var km KIMutex
	km.count = count
	km.handle = ELFHash
	km.locks = make([]Mutex, count, count)
	return &km
}

// Count the number of mutexs
func (km *KIMutex) Count() uint {
	return km.count
}

// LockID lock by idx
func (km *KIMutex) LockID(idx uint) {
	km.locks[idx%km.count].Lock()
}

// UnlockID unlock by idx
func (km *KIMutex) UnlockID(idx uint) {
	km.locks[idx%km.count].Unlock()
}

// TryLock nonblock lock
func (km *KIMutex) TryLockID(idx uint) bool {
	lock := km.locks[idx%km.count]
	select {
	case lock.ch <- Lock:
		return true
	default:
		return false
	}
}

// LockIDWithTimeout lock with timeout
func (km *KIMutex) LockIDWithTimeout(idx uint, t time.Duration) bool {
	lock := km.locks[idx%km.count]
	select {
	case lock.ch <- Lock:
		return true
	case <-time.After(t):
		return false
	}
}

// LockIDChan it can support select with this
func (km *KIMutex) LockIDChan(idx uint) chan<- Unit {
	lock := km.locks[idx%km.count]
	return lock.ch
}

// Lock the key
func (km *KIMutex) Lock(key string) {
	km.LockID(km.handle(key))
}

// Unlock the key
func (km *KIMutex) Unlock(key string) {
	km.UnlockID(km.handle(key))
}

// TryLock nonblock lock
func (km *KIMutex) TryLock(key string) bool {
	lock := km.locks[km.handle(key)%km.count]
	select {
	case lock.ch <- Lock:
		return true
	default:
		return false
	}
}

// LockWithTimeout lock with timeout
func (km *KIMutex) LockWithTimeout(key string, t time.Duration) bool {
	lock := km.locks[km.handle(key)%km.count]
	select {
	case lock.ch <- Lock:
		return true
	case <-time.After(t):
		return false
	}
}

// LockChan it can support select with this
func (km *KIMutex) LockChan(key string) chan<- Unit {
	lock := km.locks[km.handle(key)%km.count]
	return lock.ch
}

