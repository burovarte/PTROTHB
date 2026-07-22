package main

import (
	"fmt"
	doublelinkedlist "github/burovarte/PTROTHB/doubleLinkedList"
	"strconv"
)

func handleDoubleLinkedList(doubleLinkedList *doublelinkedlist.List, args []string) {

	switch args[0] {
	case "pushFront":

		n, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("это не число")
			return
		}

		doubleLinkedList.PushFront(n)

	case "pushBack":
		n, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("это не число")
			return
		}

		doubleLinkedList.PushBack(n)

	case "popFront":

		doubleLinkedList.PopFront()

	case "popBack":
		doubleLinkedList.PopBack()

	case "peekFront":
		doubleLinkedList.PeekFront()

	case "peekBack":
		doubleLinkedList.PeekBack()

	case "printList":
		doubleLinkedList.PrintList()

	case "size":
		doubleLinkedList.Size()
	}

}
