//go:build !solution

package dupcall

import ("context"
 "sync")

type invocation struct {
	result interface{}
	err error

	done chan struct{}

	waiters int 

	cancel context.CancelFunc
}

type Call struct {
	mu sync.Mutex
	current *invocation
}

func (o *Call) Do(
	ctx context.Context,
	cb func(context.Context) (interface{}, error),
) (result interface{}, err error) {
	o.mu.Lock()

	var curr *invocation

	var ctxCurrent context.Context
	shouldStart := false

	if o.current == nil {
		var cancel context.CancelFunc
ctxCurrent, cancel = context.WithCancel(context.Background())

		curr = &invocation{
			done:    make(chan struct{}),
			waiters: 1,
			cancel:  cancel,
		}

		o.current = curr
		shouldStart = true

	} else {
		curr = o.current
		curr.waiters++
	}

	o.mu.Unlock()

	if shouldStart {
		go func() {
			result, err := cb(ctxCurrent)

			o.mu.Lock()
			curr.result = result
			curr.err = err
			close(curr.done)
			if o.current == curr {
	o.current = nil
}
			o.mu.Unlock()
		}()
	}

	select {
	case <-curr.done:
		return curr.result, curr.err
	case <-ctx.Done():
		o.mu.Lock()
		curr.waiters--
		
	if curr.waiters == 0 {
		if o.current == curr {
			o.current = nil
		}
		curr.cancel()
	}

		o.mu.Unlock()
		return nil, ctx.Err()
	}

	
}
