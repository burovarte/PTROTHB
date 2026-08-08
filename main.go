package main

import (
	"bufio"
	"fmt"
	"github/burovarte/PTROTHB/deque"
	doublelinkedlist "github/burovarte/PTROTHB/doubleLinkedList"
	"github/burovarte/PTROTHB/heap"
	linkedlist "github/burovarte/PTROTHB/linkedList"
	"github/burovarte/PTROTHB/patterns"
	"github/burovarte/PTROTHB/queue"
	"github/burovarte/PTROTHB/stack"
	"os"
	"strings"
)

func main() {
	fmt.Println("Hello, World!")

	patterns.Pipeline()

	patterns.MainTimeout()

	patterns.MainOrDone()

	patterns.MainErrGroup()

	patterns.MainRateLimiting()

	patterns.MainOrChanel()

	patterns.MainBridge()

	scanner := bufio.NewScanner(os.Stdin)

	deque := &deque.Deque{}
	doublelinkedlist := &doublelinkedlist.List{}
	heap := &heap.Heap{}
	linkedlist := &linkedlist.List{}
	queue := &queue.Queue{}
	stack := &stack.Stack{}

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Fields(line)

		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "deque":
			handleDeque(deque, parts[1:])
		case "doublelinkedlist":
			handleDoubleLinkedList(doublelinkedlist, parts[1:])
		case "heap":
			handleHeap(heap, parts[1:])
		case "linkedlist":
			handleLinkedList(linkedlist, parts[1:])
		case "queue":
			handleQueue(queue, parts[1:])
		case "stack":
			handleStack(stack, parts[1:])

		}
	}
}
