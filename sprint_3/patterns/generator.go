package patterns

func generator() <-chan int {
	out := make(chan int)

	go func() {
		for i := 0; i < 5; i++ {
			out <- i
		}

		close(out)
	}()

	return out
}

func generatorMain() {
	res := generator()

	println(res)
}
