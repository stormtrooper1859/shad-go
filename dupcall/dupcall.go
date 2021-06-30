// +build !solution

package dupcall

import (
	"context"
	"sync"
)

type Call struct {
	lock         sync.Mutex
	counter      int
	done         chan int
	cancellation context.CancelFunc

	result interface{}
	err    error
}

func (o *Call) Do(
	ctx context.Context,
	cb func(context.Context) (interface{}, error),
) (result interface{}, err error) {
	o.lock.Lock()
	if o.counter == 0 {
		o.done = make(chan int)
		innerCtx, cancel := context.WithCancel(context.Background())
		o.cancellation = cancel
		go func() {
			o.result, o.err = cb(innerCtx)
			close(o.done)
		}()
	}
	o.counter += 1
	o.lock.Unlock()

	select {
	case <-ctx.Done():
		o.lock.Lock()
		o.counter -= 1
		if o.counter == 0 {
			o.cancellation()
		}
		o.lock.Unlock()
		err = ctx.Err()
	case <-o.done:
		o.lock.Lock()
		o.counter -= 1
		result, err = o.result, o.err
		o.lock.Unlock()
	}

	return result, err
}
