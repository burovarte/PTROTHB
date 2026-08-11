package patterns

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func MainErrGroup() {
	fmt.Println("--- errgroup ---")

	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		select {
		case <-time.After(200 * time.Millisecond):
			fmt.Println("task A: ok")
			return nil
		case <-ctx.Done():
			fmt.Println("task A: cancelled")
			return ctx.Err()
		}
	})

	g.Go(func() error {
		select {
		case <-time.After(500 * time.Millisecond):
			fmt.Println("task B: failed")
			return errors.New("task B failed")
		case <-ctx.Done():
			fmt.Println("task B: cancelled")
			return ctx.Err()
		}
	})

	g.Go(func() error {
		fmt.Println("task C: started")
		select {
		case <-time.After(3 * time.Second):
			fmt.Println("task C: finished (unexpected)")
			return nil
		case <-ctx.Done():
			fmt.Println("task C: cancelled by ctx")
			return ctx.Err()
		}
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("errgroup Wait: %v\n", err)
	}
}
