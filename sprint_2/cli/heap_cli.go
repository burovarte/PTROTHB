package cli

import (
	"fmt"
	"github/burovarte/PTROTHB/sprint_2/heap"
	"strconv"
)

func HandleHeap(heap *heap.Heap, args []string) {

	switch args[0] {
	case "push":
		n, err := strconv.Atoi(args[1])

		if err != nil {
			fmt.Println("Это не число")
			return
		}

		heap.Push(n)

	case "pop":
		heap.Pop()

	case "peek":
		heap.Peek()

	case "size":
		heap.Size()
	}
}
