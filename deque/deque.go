package deque

import "fmt"

type Deque struct {
	items []any
}

func (d *Deque) PushFront(el any) {
	var newDeque []any
	newDeque = append(newDeque, el)
	newDeque = append(newDeque, d.items...)
	d.items = newDeque
}

func (d *Deque) PushBack(el any) {
	d.items = append(d.items, el)
}

func (d *Deque) PopFront() {
	if len(d.items) == 0 {
		return
	}

	delEl := d.items[0]

	d.items = d.items[1:]

	fmt.Printf("Delete el (%v) from front\n", delEl)
}

func (d *Deque) PopBack() {
	if len(d.items) == 0 {
		return
	}

	delEl := d.items[len(d.items)-1]

	d.items = d.items[:len(d.items)-1]

	fmt.Printf("Delete el (%v) from back\n", delEl)
}

func (d *Deque) Front() {
	fmt.Printf("Front element is %v", d.items[0])
}

func (d *Deque) Back() {
	fmt.Printf("Back element is %v", d.items[len(d.items)-1])
}

func (d *Deque) Size() {
	fmt.Printf("Size of deques is %v", len(d.items))
}
