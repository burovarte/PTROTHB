package patterns

import (
	"fmt"
)

func gen(nums ...int) chan int {
	out := make(chan int)

	go func() {
		for _, num := range nums {
			out <- num
		}

		close(out)
	}()

	return out
}

func square(a chan int) chan int {
	out := make(chan int)

	go func() {
		for v := range a {
			out <- v * v

		}

		close(out)
	}()

	return out
}

func somePrint(a chan int) chan int {
	out := make(chan int)

	go func() {
		for v := range a {

			fmt.Println("Я ещё кручусь в последней горутине")

			out <- v

		}
		close(out)
	}()

	return out

}

func Pipeline() {

	for v := range somePrint(square(gen(1, 2, 34, 67))) {
		fmt.Println(v)
	}

}
