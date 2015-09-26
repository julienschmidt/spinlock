// Copyright 2015 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package spinlock

import (
	"runtime"
	"sync/atomic"
)

const (
	mutexUnlocked = 0
	mutexLocked   = 1
)

type Mutex struct {
	state int32
}

// Lock locks m.
// If the lock is already in use, the calling goroutine repetitively tries to
// acquire the the mutex until it is available (busy waiting).
func (m *Mutex) Lock() {
	for !atomic.CompareAndSwapInt32(&m.state, mutexUnlocked, mutexLocked) {
		runtime.Gosched()
	}
}

// TryLock tries to lock m.
// If the lock is already in use, the lock is not acquired and false is
// returned.
func (m *Mutex) TryLock() bool {
	return atomic.CompareAndSwapInt32(&m.state, mutexUnlocked, mutexLocked)
}

// Unlock unlocks m.
// It is a run-time error if m is not locked on entry to Unlock.
//
// A locked Mutex is not associated with a particular goroutine.
// It is allowed for one goroutine to lock a Mutex and then
// arrange for another goroutine to unlock it.
func (m *Mutex) Unlock() {
	state := atomic.AddInt32(&m.state, -mutexLocked)
	if state != mutexUnlocked {
		panic("spinlock: unlock of unlocked mutex")
	}
}
