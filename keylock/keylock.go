// +build !solution

package keylock

import (
	"sync"
)

type KeyLock struct {
	set  map[string]struct{}
	cond *sync.Cond
}

func New() *KeyLock {
	return &KeyLock{
		set:  make(map[string]struct{}),
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

func checkContain(set map[string]struct{}, keys []string) bool {
	for _, s := range keys {
		if _, c := set[s]; c {
			return true
		}
	}
	return false
}

func (l *KeyLock) LockKeys(keys []string, cancel <-chan struct{}) (canceled bool, unlock func()) {
	c := make(chan struct{}, 1)

	l.cond.L.Lock()
	for !canceled && checkContain(l.set, keys) {
		go func() {
			l.cond.L.Lock()
			if checkContain(l.set, keys) {
				l.cond.Wait()
			}
			l.cond.L.Unlock()
			c <- struct{}{}
		}()
		l.cond.L.Unlock()

		select {
		case <-cancel:
			canceled = true
		case <-c:
			canceled = false
		}
		l.cond.L.Lock()
	}

	if canceled {
		l.cond.Broadcast()
		l.cond.L.Unlock()
		return canceled, nil
	}

	for _, s := range keys {
		l.set[s] = struct{}{}
	}
	l.cond.L.Unlock()

	unlock = func() {
		l.cond.L.Lock()
		for _, s := range keys {
			delete(l.set, s)
		}

		l.cond.Broadcast()
		l.cond.L.Unlock()
	}

	return canceled, unlock
}
