// +build !solution

package pubsub

import (
	"context"
	"errors"
	"sync"
)

var _ Subscription = (*MySubscription)(nil)

type MySubscription struct {
	cancel context.CancelFunc
}

func (s *MySubscription) Unsubscribe() {
	s.cancel()
}

var _ PubSub = (*MyPubSub)(nil)

type channel struct {
	current chan wtf
}

type wtf struct {
	msg    interface{}
	w      chan wtf
	closed bool
}

type MyPubSub struct {
	wg sync.WaitGroup
	mp map[string]*channel

	mtx           sync.Mutex
	mainCtx       context.Context
	cancelMainCtx context.CancelFunc

	closed bool
}

func NewPubSub() PubSub {
	mps := &MyPubSub{
		mp: make(map[string]*channel),
	}
	mps.mainCtx, mps.cancelMainCtx = context.WithCancel(context.Background())
	return mps
}

func newChannel() *channel {
	return &channel{
		current: make(chan wtf, 1),
	}
}

func handler(ctx context.Context, w wtf, cb MsgHandler, group *sync.WaitGroup) {
	nw := w
f:
	for {
		select {
		case nw2 := <-nw.w:
			nw.w <- nw2
			if nw2.closed {
				break f
			}
			nw = nw2
			cb(nw2.msg)
		case <-ctx.Done():
			break f
		}
	}
	group.Done()
}

func (p *MyPubSub) Subscribe(subj string, cb MsgHandler) (Subscription, error) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.closed {
		return nil, ClosedError
	}

	ch, contain := p.mp[subj]

	if !contain {
		ch = newChannel()
		p.mp[subj] = ch
	}

	ctx, cancel := context.WithCancel(context.Background())

	p.wg.Add(1)

	go handler(ctx, wtf{w: ch.current}, cb, &p.wg)

	ms := &MySubscription{
		cancel: cancel,
	}

	return ms, nil
}

func (p *MyPubSub) Publish(subj string, msg interface{}) error {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.closed {
		return ClosedError
	}

	ch, contain := p.mp[subj]

	if !contain {
		return nil
	}

	if msg == nil {
		return nil
	}

	w := wtf{
		msg: msg,
		w:   make(chan wtf, 1),
	}

	ch.current <- w
	//close(ch.current)
	ch.current = w.w

	return nil
}

var ClosedError = errors.New("pubsub already closed")

func (p *MyPubSub) Close(ctx context.Context) error {
	p.mtx.Lock()
	p.closed = true

	for _, v := range p.mp {
		w := wtf{
			w:      make(chan wtf, 1),
			closed: true,
		}

		v.current <- w
		//close(v.current)
		v.current = w.w
	}

	p.mtx.Unlock()

	c1 := make(chan int)

	go func() {
		p.wg.Wait()
		close(c1)
	}()

	select {
	case <-c1:
	case <-ctx.Done():
		p.cancelMainCtx()
		return ctx.Err()
	}

	return nil
}
