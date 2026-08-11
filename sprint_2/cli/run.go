package cli

import (
	"bufio"
	"os"
	"strings"

	"github/burovarte/PTROTHB/sprint_2/deque"
	doublelinkedlist "github/burovarte/PTROTHB/sprint_2/doubleLinkedList"
	"github/burovarte/PTROTHB/sprint_2/heap"
	linkedlist "github/burovarte/PTROTHB/sprint_2/linkedList"
	"github/burovarte/PTROTHB/sprint_2/queue"
	"github/burovarte/PTROTHB/sprint_2/stack"
)

// Run читает команды из stdin и прокидывает их в нужный handler.
func Run() {
	scanner := bufio.NewScanner(os.Stdin)

	d := &deque.Deque{}
	dll := &doublelinkedlist.List{}
	h := &heap.Heap{}
	ll := &linkedlist.List{}
	q := &queue.Queue{}
	s := &stack.Stack{}

	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "deque":
			HandleDeque(d, parts[1:])
		case "doublelinkedlist":
			HandleDoubleLinkedList(dll, parts[1:])
		case "heap":
			HandleHeap(h, parts[1:])
		case "linkedlist":
			HandleLinkedList(ll, parts[1:])
		case "queue":
			HandleQueue(q, parts[1:])
		case "stack":
			HandleStack(s, parts[1:])
		}
	}
}
