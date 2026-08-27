//go:build !solution

package keylock

import (
	"sync"
)

type KeyLock struct {
	mu      sync.Mutex
	locked  map[string]struct{}
	changed chan struct{}
}

func New() *KeyLock {
	return &KeyLock{
		locked:  make(map[string]struct{}),
		changed: make(chan struct{}),
	}
}

func (l *KeyLock) hasConflict(keys []string) bool {
	for _, key := range keys {
		if _, exists := l.locked[key]; exists {
			return true
		}
	}

	return false
}

func (l *KeyLock) LockKeys(keys []string, cancel <-chan struct{}) (canceled bool, unlock func()) {
	for {
		l.mu.Lock()

		if !l.hasConflict(keys) {
			for _, key := range keys {
				l.locked[key] = struct{}{}
			}

			l.mu.Unlock()

			unlock := func() {
				l.mu.Lock()

				for _, key := range keys {
					delete(l.locked, key)
				}

				close(l.changed)

				l.changed = make(chan struct{})

				l.mu.Unlock()
			}

			return false, unlock
		}

		changed := l.changed
		l.mu.Unlock()

		select {
		case <-changed:

		case <-cancel:
			return true, nil
		}
	}
}
