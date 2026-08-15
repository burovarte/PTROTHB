//go:build !solution

package ratelimit

import (
	"context"
	"errors"
	"time"
)

type request struct {
	ctx    context.Context
	result chan error
}

// Limiter is precise rate limiter with context support.
type Limiter struct {
	requests chan request
	stop     chan struct{}
	done     chan struct{}
}

var ErrStopped = errors.New("limiter stopped")

// NewLimiter returns limiter that throttles rate of successful Acquire() calls
// to maxSize events at any given interval.
func NewLimiter(maxCount int, interval time.Duration) *Limiter {
	limiter := &Limiter{
		requests: make(chan request),
		stop:     make(chan struct{}),
		done:     make(chan struct{}),
	}

	go limiter.run(maxCount, interval)

	return limiter
}

func (l *Limiter) Acquire(ctx context.Context) error {
	result := make(chan error, 1)
	req := request{
		ctx:    ctx,
		result: result,
	}

	select {
	case l.requests <- req:
	case <-ctx.Done():
		return ctx.Err()
	case <-l.done:
		return ErrStopped
	}

	select {
	case err := <-result:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-l.done:
		return ErrStopped
	}
}

func (l *Limiter) Stop() {
	select {
	case l.stop <- struct{}{}:
		<-l.done
	case <-l.done:
	}
}

func (l *Limiter) run(maxCount int, interval time.Duration) {
	defer close(l.done)

	var waiting []request
	var acquiredAt []time.Time

	for {
		now := time.Now()
		if interval > 0 {
			firstActive := 0
			for firstActive < len(acquiredAt) && !acquiredAt[firstActive].Add(interval).After(now) {
				firstActive++
			}
			acquiredAt = acquiredAt[firstActive:]
		} else {
			acquiredAt = acquiredAt[:0]
		}

		for len(waiting) > 0 && (interval <= 0 || len(acquiredAt) < maxCount) {
			req := waiting[0]
			waiting[0] = request{}
			waiting = waiting[1:]

			if err := req.ctx.Err(); err != nil {
				req.result <- err
				continue
			}

			req.result <- nil
			if interval > 0 {
				acquiredAt = append(acquiredAt, time.Now())
			}
		}

		var timer *time.Timer
		var timerC <-chan time.Time
		if len(waiting) > 0 && interval > 0 && maxCount > 0 && len(acquiredAt) >= maxCount {
			delay := time.Until(acquiredAt[0].Add(interval))
			if delay < 0 {
				delay = 0
			}
			timer = time.NewTimer(delay)
			timerC = timer.C
		}

		select {
		case req := <-l.requests:
			waiting = append(waiting, req)
		case <-timerC:
		case <-l.stop:
			if timer != nil {
				timer.Stop()
			}
			return
		}

		if timer != nil && !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
	}
}
