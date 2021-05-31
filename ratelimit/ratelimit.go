// +build !solution

package ratelimit

import (
	"context"
	"errors"
	"time"
)

type ct <-chan time.Time

// Limiter is precise rate limiter with context support.
type Limiter struct {
	interval time.Duration
	stopped  chan struct{}
	ch2      chan ct
}

var ErrStopped = errors.New("limiter stopped")

// NewLimiter returns limiter that throttles rate of successful Acquire() calls
// to maxSize events at any given interval.
func NewLimiter(maxCount int, interval time.Duration) *Limiter {
	limiter := Limiter{
		interval: interval,
		stopped:  make(chan struct{}),
		ch2:      make(chan ct, maxCount),
	}
	for i := 0; i < maxCount; i++ {
		c := make(chan time.Time, 1)
		c <- time.Time{}
		limiter.ch2 <- c
	}
	return &limiter
}

func (l *Limiter) Acquire(ctx context.Context) error {
	select {
	case <-l.stopped:
		return ErrStopped
	default:
	}

	select {
	case cc := <-l.ch2:
		select {
		case <-cc:
			r := time.After(l.interval)
			l.ch2 <- r
		case <-ctx.Done():
			return ctx.Err()
		case <-l.stopped:
			return ErrStopped
		}
	case <-ctx.Done():
		return ctx.Err()
	case <-l.stopped:
		return ErrStopped
	}

	return nil
}

func (l *Limiter) Stop() {
	close(l.stopped)
}
