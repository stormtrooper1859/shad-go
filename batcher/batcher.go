// +build !solution

package batcher

import (
	"gitlab.com/slon/shad-go/batcher/slow"
	"sync/atomic"
	"unsafe"
)

type entry struct {
	value interface{}
	ready chan struct{}
}

type Batcher struct {
	slowValue *slow.Value
	entry     unsafe.Pointer
}

var nullEntry *entry = &entry{}

func NewBatcher(v *slow.Value) *Batcher {
	b := Batcher{
		slowValue: v,
		entry:     unsafe.Pointer(nullEntry),
	}
	return &b
}

func (b *Batcher) Load() interface{} {
	var currentEntry *entry
	for {
		currentEntry = (*entry)(atomic.LoadPointer(&b.entry))

		if currentEntry != nullEntry {
			break
		}

		currentEntry = &entry{
			ready: make(chan struct{}),
		}

		if !atomic.CompareAndSwapPointer(&b.entry, unsafe.Pointer(nullEntry), unsafe.Pointer(currentEntry)) {
			continue
		}

		currentEntry.value = b.slowValue.Load()

		atomic.StorePointer(&b.entry, unsafe.Pointer(nullEntry))

		close(currentEntry.ready)

		break
	}

	<-currentEntry.ready

	return currentEntry.value
}
