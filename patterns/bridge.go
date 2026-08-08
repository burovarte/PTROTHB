package patterns

import (
	"context"
	"fmt"
	"time"
)

func chanProducer(ctx context.Context) <-chan (<-chan int) {
	out := make(chan (<-chan int))

	go func() {
		defer close(out)

		batches := [][]int{{1, 2}, {10, 20}, {100}}

		for _, batch := range batches {
			select {
			case <-ctx.Done():
				return
			case <-time.After(500 * time.Millisecond):
			}

			inner := make(chan int, len(batch))

			for _, v := range batch {
				inner <- v
			}
			close(inner)

			select {
			case out <- inner:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func bridge(ctx context.Context, chanStream <-chan (<-chan int)) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for {
			select {
			case inner, ok := <-chanStream:
				if !ok {
					return
				}

				for el := range inner {
					select {
					case out <- el:
					case <-ctx.Done():
						return
					}
				}

			case <-ctx.Done():
				return

			}
		}

	}()

	return out
}

func MainBridge() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for v := range bridge(ctx, chanProducer(ctx)) {
		fmt.Println(v)
	}
}
