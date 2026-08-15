//go:build !solution

package rwmutex

// A RWMutex is a reader/writer mutual exclusion lock.
// The lock can be held by an arbitrary number of readers or a single writer.
// The zero value for a RWMutex is an unlocked mutex.
//
// If a goroutine holds a RWMutex for reading and another goroutine might
// call Lock, no goroutine should expect to be able to acquire a read lock
// until the initial read lock is released. In particular, this prohibits
// recursive read locking. This is to ensure that the lock eventually becomes
// available; a blocked Lock call excludes new readers from acquiring the
// lock.
type RWMutex struct {
	countReader chan int
	lockWriter  chan struct{}
}

// New creates *RWMutex.
func New() *RWMutex {
	rwmutex := &RWMutex{
		countReader: make(chan int, 1),
		lockWriter:  make(chan struct{}, 1),
	}

	rwmutex.countReader <- 0
	rwmutex.lockWriter <- struct{}{}

	return rwmutex
}

// RLock locks rw for reading.
//
// It should not be used for recursive read locking; a blocked Lock
// call excludes new readers from acquiring the lock. See the
// documentation on the RWMutex type.
func (rw *RWMutex) RLock() {
	currentCountReader := <-rw.countReader

	if currentCountReader == 0 {
		<-rw.lockWriter
	}

	currentCountReader++

	rw.countReader <- currentCountReader
}

// RUnlock undoes a single RLock call;
// it does not affect other simultaneous readers.
// It is a run-time error if rw is not locked for reading
// on entry to RUnlock.
func (rw *RWMutex) RUnlock() {
	currentCountReader := <-rw.countReader

	if currentCountReader == 0 {
		rw.countReader <- currentCountReader
		panic("It is a run-time error if rw is not locked for reading")
	}

	currentCountReader--

	if currentCountReader == 0 {
		rw.lockWriter <- struct{}{}
	}

	rw.countReader <- currentCountReader
}

// Lock locks rw for writing.
// If the lock is already locked for reading or writing,
// Lock blocks until the lock is available.
func (rw *RWMutex) Lock() {
	<-rw.lockWriter
}

// Unlock unlocks rw for writing. It is a run-time error if rw is
// not locked for writing on entry to Unlock.
//
// As with Mutexes, a locked RWMutex is not associated with a particular
// goroutine. One goroutine may RLock (Lock) a RWMutex and then
// arrange for another goroutine to RUnlock (Unlock) it.
func (rw *RWMutex) Unlock() {

	select {
	case rw.lockWriter <- struct{}{}:

	default:
		panic("rwmutex: unlock of unlocked mutex")
	}
}
