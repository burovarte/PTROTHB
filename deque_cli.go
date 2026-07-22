package main

import (
	"fmt"
	"github/burovarte/PTROTHB/deque"
)

func handleDeque(d *deque.Deque, args []string) {
	if len(args) == 0 {
		fmt.Println("Need correct command")
		return
	}

	switch args[0] {
	case "pushFront":
		if len(args) < 2 {
			fmt.Println("need: pushFront <value>")
			return
		}

		d.PushFront(args[1])
	case "pushBack":
		d.PushBack(args[1])
	case "popFront":
		d.PopFront()
	case "popBack":
		d.PopBack()
	case "front":
		d.Front()
	case "back":
		d.Back()
	case "size":
		d.Size()
	}
}
