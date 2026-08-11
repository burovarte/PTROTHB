package patterns

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func simulateRequest(ctx context.Context, limiter *rate.Limiter, id int) {

	err := limiter.Wait(ctx)

	if err != nil {
		return
	}

	fmt.Printf("%s request #%d\n", time.Now().Format("15:04:05.000"), id)

}

func MainRateLimiting() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var wg sync.WaitGroup

	limiter := rate.NewLimiter(rate.Every(500*time.Millisecond), 2)

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			simulateRequest(ctx, limiter, i)

			wg.Done()
		}()
	}

	wg.Wait()
}
