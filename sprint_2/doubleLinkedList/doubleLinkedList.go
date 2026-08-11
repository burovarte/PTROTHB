package doublelinkedlist

import "fmt"

type Node struct {
	Value int
	Next  *Node
	Prev  *Node
}

type List struct {
	Head *Node
	Tail *Node
}

func (l *List) PushFront(el int) {
	n := &Node{Value: el}

	if l.Head == nil {
		l.Head = n
		l.Tail = n
		return
	}
	n.Next = l.Head
	l.Head.Prev = n
	l.Head = n
}

func (l *List) PushBack(el int) {
	n := &Node{Value: el}

	if l.Head == nil {
		l.Head = n
		l.Tail = n
		return
	}

	n.Prev = l.Tail
	l.Tail.Next = n
	l.Tail = n
}

func (l *List) PopFront() {
	if l.Head == nil {
		return
	}

	firstEl := l.Head.Value

	l.Head = l.Head.Next

	if l.Head != nil {
		l.Head.Prev = nil
	} else {
		l.Tail = nil
	}

	fmt.Printf("Head of linked list os %d\n", firstEl)
}

func (l *List) PopBack() {
	if l.Tail == nil {
		return
	}

	lastEl := l.Tail.Value

	l.Tail = l.Tail.Prev

	if l.Tail != nil {
		l.Tail.Next = nil
	} else {
		l.Head = nil
	}

	fmt.Printf("Tail of linked list os %d\n", lastEl)
}

func (l *List) PeekFront() {
	if l.Head == nil {
		return
	}

	fmt.Printf("Head of linked list os %d\n", l.Head.Value)
}

func (l *List) PeekBack() {
	if l.Tail == nil {
		return
	}

	fmt.Printf("Tail of linked list os %d\n", l.Tail.Value)
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
