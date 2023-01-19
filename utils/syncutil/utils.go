package syncutil

import "sync"

type RLocked[T any] struct {
	lock  sync.RWMutex
	value T
}

func (rl *RLocked[T]) Store(value T) {
	rl.lock.Lock()
	defer rl.lock.Unlock()
	rl.value = value
}

func (rl *RLocked[T]) Load() T {
	rl.lock.RLock()
	defer rl.lock.RUnlock()
	return rl.value
}
