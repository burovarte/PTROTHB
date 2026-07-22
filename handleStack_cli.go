package main

import (
	"fmt"
	"github/burovarte/PTROTHB/stack"
	"strconv"
)

func handleStack(stack *stack.Stack, args []string) {
	switch args[0] {
	case "push":

		n, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Это не число")
			return
		}
		stack.Push(n)
	case "pop":
		stack.Pop()
	case "peek":
		stack.Peek()
	}
}
