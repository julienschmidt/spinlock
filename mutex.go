package spinlock

import (
	"runtime"
	"sync/atomic"
)

const (
	mutexUnlocked = 0
	mutexLocked   = 1
)

type Mutex int32

// Lock locks m.
// If the lock is already in use, the calling goroutine repetitively tries to
// acquire the the mutex until it is available (busy waiting).
func (m *Mutex) Lock() {
	for !atomic.CompareAndSwapInt32(&m, mutexUnlocked, mutexLocked) {
	}
}

// LockYield locks m.
// Behaves like Lock, but yields the processor, allowing other goroutines to
// run, each time it fails to acquire the mutex.
func (m *Mutex) LockYield() {
	for !atomic.CompareAndSwapInt32(&m, mutexUnlocked, mutexLocked) {
		runtime.Gosched()
	}
}

func (m *Mutex) Unlock() {
	state := atomic.AddInt32(&m, -mutexLocked)
	if (state + mutexLocked) != mutexUnlocked {
		panic("spinlock: unlock of unlocked mutex")
	}
}
