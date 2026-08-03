package patterns

import (
	"fmt"
	"sync"
)

// источник: кладёт задачи в канал и закрывает его
func fanGenRepeat(nums ...int) chan int {
	ch := make(chan int)

	go func() {
		for _, num := range nums {
			ch <- num
		}
		close(ch)
	}()

	return ch
}

// воркер: читает из общего канала, пока тот не закрыт
func workerFanOutRepeat(jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		fmt.Println(job)
	}
}

// fan-out: N воркеров на одном канале
func fanOutRepeat() {
	jobs := fanGenRepeat(4, 5, 34, 53)

	const n = 2 // фиксированное число воркеров

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go workerFanOutRepeat(jobs, &wg)
	}

	wg.Wait()
}
