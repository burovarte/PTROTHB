package main

import "github/burovarte/PTROTHB/queue"

func handleQueue(queue *queue.Queue, args []string) {

	switch args[0] {
	case "push":
		queue.Push(args[1])
	case "pop":
		queue.Pop()
	case "size":
		queue.Size()
	case "peek":
		queue.Peek()
	}

}
