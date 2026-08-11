package patterns

import (
	"fmt"
	"sync"
)

func fanGen(nums ...int) chan int {
	ch := make(chan int)

	go func() {

		for _, num := range nums {
			ch <- num

		}

		close(ch)
	}()
	return ch
}

func workerFanOut(jobs chan int, wg *sync.WaitGroup) {

	for job := range jobs {

		fmt.Println(job)

	}

	wg.Done()
}

func fanOut() {

	a := fanGen(3, 4, 57676, 53442)

	var wg sync.WaitGroup

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go workerFanOut(a, &wg)

	}

	wg.Wait()

}
