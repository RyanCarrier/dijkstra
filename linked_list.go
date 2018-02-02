package dijkstra

import (
	"fmt"
)

//THE FOLLOWING FILE IS BASED FROM GO AUTHORS EDITED MINORLY AND LAZILY TO SUIT MY NEEDS
//https://golang.org/src/container/list/list.go?m=text
//AVOID USING MINE AS A TEMPLATE AS I REMOVED MOST SAFETIES (that's why they are
// all private now)

// linkedList represents a doubly linked list.
// The zero value for linkedList is an empty list ready to use.
type linkedList struct {
	root  *Vertex // sentinel list element, only &root, root.prev, and root.next are used
	len   int     // current list length excluding (this) sentinel element
	short bool
}

func (l linkedList) print() {
	current := l.root
	fmt.Print("Next = ROOT")
	current = current.next
	for ; current != nil && current != l.root; current = current.next {
		fmt.Printf("->%d", current.ID)
	}
	if current == nil {
		fmt.Print("->NIL")
	}
	fmt.Println("->ROOT")
	l.reversePrint()
}

func (l linkedList) reversePrint() {
	current := l.root
	fmt.Print("Prev = ROOT")
	current = current.prev
	for ; current != nil && current != l.root; current = current.prev {
		fmt.Printf("->%d", current.ID)
	}
	if current == nil {
		fmt.Print("->NIL")
	}
	fmt.Println("->ROOT")
}

// Init initializes or clears list l.
func linkedListNewShort() dijkstraList {
	return dijkstraList(new(linkedList).init(true))
}

// Init initializes or clears list l.
func linkedListNewLong() dijkstraList {
	return dijkstraList(new(linkedList).init(false))
}

// Init initializes or clears list l.
func (l *linkedList) PushOrdered(v *Vertex) {
	l.pushOrdered(v)
}

// Init initializes or clears list l.
func (l *linkedList) PopOrdered() *Vertex {
	if l.short {
		return l.popBack()
	}
	return l.popFront()
}

// Init initializes or clears list l.
func (l *linkedList) Len() int {
	return l.len
}

// Init initializes or clears list l.
func (l *linkedList) init(short bool) *linkedList {
	l.root = NewVertex(-1)
	l.root.next = l.root
	l.root.prev = l.root
	l.len = 0
	l.short = short
	return l
}

// front returns the first element of list l or nil.
func (l *linkedList) front() *Vertex {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// back returns the last element of list l or nil.
func (l *linkedList) back() *Vertex {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

//popFront pops the Vertex off the front of the list
func (l *linkedList) popFront() *Vertex {
	return l.remove(l.front())
}

//popFront pops the Vertex off the back of the list
func (l *linkedList) popBack() *Vertex {
	return l.remove(l.back())
}

func (l *linkedList) fixV(v *Vertex) *Vertex {
	return l.pushOrdered(l.remove(v))
}

//pushOrdered pushes the value into the linked list in the correct position
// (ascending)
func (l *linkedList) pushOrdered(v *Vertex) *Vertex {
	if v.inList {
		return l.fixV(v)
	}
	if l.len == 0 {
		return l.pushFront(v)
	}
	if l.back().compare(v) < 0 {
		return l.insert(v, l.root.prev)
	}
	current := l.front()
	for ; current.compare(v) < 0; current = current.next {
	}
	return l.insert(v, current.prev)
}

// insert inserts e after at, increments l.len, and returns e.
func (l *linkedList) insert(e, at *Vertex) *Vertex {
	n := at.next
	at.next = e
	e.prev = at

	e.next = n
	n.prev = e
	l.len++
	e.inList = true
	return e
}

// remove removes e from its list, decrements l.len, and returns e.
func (l *linkedList) remove(e *Vertex) *Vertex {
	e.prev.next = e.next
	e.next.prev = e.prev
	l.len--
	e.inList = false
	return e
}

// pushFront inserts a new element e with value v at the front of list l and returns e.
func (l *linkedList) pushFront(v *Vertex) *Vertex {
	return l.insert(v, l.root)
}
