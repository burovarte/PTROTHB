package patterns

import (
	"fmt"
	"sync"
)

func generatorFanIn(num int) <-chan int {
	out := make(chan int)

	go func() {
		for i := 0; i < 3; i++ {
			out <- i + num
		}

		close(out)
	}()

	return out
}

func fanIn(chanels ...<-chan int) <-chan int {
	out := make(chan int)

	var wg sync.WaitGroup

	for _, chanel := range chanels {
		wg.Add(1)

		go func(c <-chan int) {
			for el := range c {
				out <- el
			}

			wg.Done()
		}(chanel)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func mainFanIn() {
	a := fanIn(generatorFanIn(1), generatorFanIn(2), generatorFanIn(3))

	for v := range a {
    fmt.Println(v)
}

}
