package main

import (
	"fmt"
	linkedlist "github/burovarte/PTROTHB/linkedList"
	"strconv"
)

func handleLinkedList(list *linkedlist.List, args []string) {
	switch args[0] {
	case "pushFront":

		n, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("это не число")
			return
		}

		list.PushFront(n)

	case "popFront":
		list.PopFront()

	case "peek":
		list.Peek()

	case "printList":
		list.PrintList()

	case "size":
		list.Size()
	}
}
