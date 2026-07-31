package patterns

import (
	"fmt"
	"sync"
)

type task int

func worker(jobs chan task, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		fmt.Println(job * job)
	}

}

func workerPool() {
	jobs := make(chan task)

	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go worker(jobs, &wg)
	}

	for i := task(0); i < 67; i++ {
		jobs <- i
	}

	close(jobs)

	wg.Wait()
}
