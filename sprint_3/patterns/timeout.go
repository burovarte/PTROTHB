package patterns

import (
	"context"
	"fmt"
	"time"
)

func taskWithTimeout() string {
	res := make(chan string)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	go func() {
		select {
		case <-time.After(7 * time.Second):
			select {
			case res <- "успех: работа выполнена за 7 сек":
			case <-ctx.Done():
			}
		case <-ctx.Done():
			return
		}
	}()

	select {
	case msg := <-res:
		return msg
	case <-ctx.Done():
		return fmt.Sprintf("timeout: лимит 3 сек, err=%v", ctx.Err())
	}
}

func taskWithDeadline() string {
	res := make(chan string)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	go func() {
		select {
		case <-time.After(7 * time.Second):
			select {
			case res <- "успех: работа выполнена за 7 сек":
			case <-ctx.Done():
			}
		case <-ctx.Done():
			return
		}
	}()

	select {
	case msg := <-res:
		return msg
	case <-ctx.Done():
		if d, ok := ctx.Deadline(); ok {
			return fmt.Sprintf("deadline: лимит до %s, err=%v", d.Format("15:04:05"), ctx.Err())
		}
		return fmt.Sprintf("deadline: err=%v", ctx.Err())
	}
}

func MainTimeout() {
	fmt.Println("--- WithTimeout ---")
	fmt.Println(taskWithTimeout())

	fmt.Println("--- WithDeadline ---")
	fmt.Println(taskWithDeadline())
}
