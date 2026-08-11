package heap

import "fmt"

type Heap struct {
	items []int
}

func (h *Heap) Push(el int) {
	h.items = append(h.items, el)

	i := len(h.items) - 1

	for i > 0 && h.items[i] > h.items[(i-1)/2] {
		h.items[i], h.items[(i-1)/2] = h.items[(i-1)/2], h.items[i]

		i = (i - 1) / 2
	}
}

func (h *Heap) Pop() (int, bool) {
	if len(h.items) == 0 {
		fmt.Println("Heap is empty")
		return 0, false
	}
	first := h.items[0]

	last := len(h.items) - 1
	h.items[0] = h.items[last]
	h.items = h.items[:last]
	h.siftDown(0)
	fmt.Printf("Pop root: %d\n", first)
	return first, true
}

func (h *Heap) siftDown(i int) {
	for {
		left := 2*i + 1
		right := 2*i + 2
		largest := i
		if left < len(h.items) && h.items[left] > h.items[largest] {
			largest = left
		}
		if right < len(h.items) && h.items[right] > h.items[largest] {
			largest = right
		}
		if largest == i {
			return
		}
		h.items[i], h.items[largest] = h.items[largest], h.items[i]
		i = largest
	}
}

func (h *Heap) Peek() (int, bool) {
	if len(h.items) == 0 {
		fmt.Println("Heap is empty")
		return 0, false
	}
	fmt.Printf("Root: %d\n", h.items[0])
	return h.items[0], true
}

func (h *Heap) Size() int {
	fmt.Printf("Size of heap is %d\n", len(h.items))
	return len(h.items)
}
