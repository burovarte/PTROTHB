package patterns

import (
	"context"
	"fmt"
	"time"
)

func generatorOrDone() <-chan int {
	out := make(chan int)

	go func() {
		for i := 0; i < 100; i++ {
			<-time.After(time.Second * 1)
			out <- i
		}

		close(out)
	}()

	return out
}

func orDone(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case out <- v:
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return

			}
		}

	}()

	return out

}

func MainOrDone() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)

	defer cancel()

	for v := range orDone(ctx, generatorOrDone()) {
		fmt.Printf("Значения из OrDone: %d\n", v)
	}

}
