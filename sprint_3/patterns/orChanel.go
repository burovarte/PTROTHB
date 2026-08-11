package patterns

import (
	"fmt"
	"time"
)

func after(d time.Duration) <-chan struct{} {
	out := make(chan struct{})

	go func() {
		<-time.After(d)

		close(out)
	}()

	return out
}

func or(channels ...<-chan struct{}) <-chan struct{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	default:
		out := make(chan struct{})
		go func() {
			select {
			case <-channels[0]:

			case <-or(channels[1:]...):
			}
			close(out)
		}()
		return out
	}
}

func MainOrChanel() {

	a := after(2 * time.Second)
	b := after(5 * time.Second)
	c := after(8 * time.Second)

	start := time.Now()
	<-or(a, b, c)
	fmt.Printf("OrChannel сработало через %v\n", time.Since(start))
}
