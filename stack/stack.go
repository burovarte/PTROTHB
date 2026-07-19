package stack

import "fmt"

type Stack struct {
	items []int
}

func (s *Stack) Push(el int) {
	s.items = append(s.items,el)

	 fmt.Printf("Add el: %d\n", el)
}

func (s *Stack) Pop(){
	if len(s.items) == 0 {
		fmt.Println("Stack is empty")
		return
	}


	el := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]

	fmt.Printf("You delete %d\n", el)
}

func (s Stack) Peek(){
	if len(s.items) == 0 {
		fmt.Println("Stack is empty")
		return
	}

	fmt.Printf("Element on peek: %d\n", s.items[len(s.items)-1] )
}

func (s Stack) Size(){
	fmt.Printf("Size is : %d\n", len(s.items))
}