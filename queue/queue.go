package queue

import "fmt"

type Queue struct {
	items []any
}

func (q *Queue) Push(el any) {
	q.items = append(q.items, el)
}

func (q *Queue) Pop() {
	if len(q.items) == 0 {
		fmt.Printf("Queue is empty\n")
		return
	}

	el := q.items[0]

	q.items = q.items[1:]

	fmt.Printf("Del fisrt el is %v\n", el)
}

func (q Queue) Peek() {
	if len(q.items) == 0 {
		fmt.Printf("Queue is empty\n")
		return
	}

	fmt.Printf("Its first el is %v\n", q.items[0])
}

func (q Queue) Size() {
	fmt.Printf("Size of queue is %v\n", len(q.items))
}
