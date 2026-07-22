package linkedlist

import "fmt"

type Node struct {
	Value int
	Next  *Node
}

type List struct {
	Head *Node
}

func (l *List) PushFront(el int) {
	n := &Node{Value: el, Next: l.Head}

	l.Head = n

}

func (l *List) PopFront() {
	if l.Head == nil {
		return
	}

	firstEl := l.Head.Value

	l.Head = l.Head.Next

	fmt.Printf("Head of linked list os %d", firstEl)
}

func (l *List) Peek() {
	if l.Head == nil {
		return
	}

	fmt.Printf("Head of linked list os %d\n", l.Head.Value)
}

func (l *List) PrintList() {
	if l.Head == nil {
		return
	}

	curEl := l.Head

	for curEl != nil {
		fmt.Printf("El of List is %d\n", curEl.Value)

		curEl = curEl.Next
	}
}

func (l *List) Size() {
	curEl := l.Head
	count := 0

	for curEl != nil {
		count++

		curEl = curEl.Next
	}

	fmt.Printf("Size of list is %d\n", count)
}
