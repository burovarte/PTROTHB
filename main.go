package main

import (
	"bufio"
	"fmt"
	"github/burovarte/PTROTHB/stack"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Hello, World!")

	scanner := bufio.NewScanner(os.Stdin)

	s := &stack.Stack{}

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Fields(line)


		switch parts[0]{
		case "pop":
			s.Pop()
		case "peek":
			s.Peek()
		case "size":
			s.Size()
		case "push":
			if len(parts) == 2{
				n, err := strconv.Atoi(parts[1])
				if err != nil {
						fmt.Println("Error")
				}

				s.Push(n)
			} else {
				fmt.Println("need number")
			}
			case "exit":
				return
			default:
				fmt.Println("try push ..., peek, size, pop or exit")
		}
	}	
}