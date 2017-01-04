package dijkstra

//THE FOLLOWING FILE IS BASED FROM GO AUTHORS EDITED MINORLY AND LAZILY TO SUIT MY NEEDS
//https://golang.org/src/container/list/list.go?m=text

// Element is an element of a linked list.
type Element struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *Element

	// The list to which this element belongs.
	list *LinkedList

	// The value stored with this element.
	Value *Vertex
}

// Next returns the next list element or nil.
func (e *Element) Next() *Element {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Prev returns the previous list element or nil.
func (e *Element) Prev() *Element {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// LinkedList represents a doubly linked list.
// The zero value for LinkedList is an empty list ready to use.
type LinkedList struct {
	root Element // sentinel list element, only &root, root.prev, and root.next are used
	len  int     // current list length excluding (this) sentinel element
}

// Init initializes or clears list l.
func (l *LinkedList) Init() *LinkedList {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// NewLinkedList returns an initialized list.
func NewLinkedList() *LinkedList { return new(LinkedList).Init() }

// Len returns the number of elements of list l.
// The complexity is O(1).
func (l *LinkedList) Len() int { return l.len }

// Front returns the first element of list l or nil.
func (l *LinkedList) Front() *Element {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

//PopFront pops the Vertex off the front of the list
func (l *LinkedList) PopFront() *Vertex {
	e := l.Front()
	if e.list == l {
		// if e.list == l, l must have been initialized when e was inserted
		// in l or l == nil (e is a zero Element) and l.remove will crash
		l.remove(e)
	}
	return e.Value
}

// Back returns the last element of list l or nil.
func (l *LinkedList) Back() *Element {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// lazyInit lazily initializes a zero LinkedList value.
func (l *LinkedList) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

//PushOrdered pushes the value into the linked list in the correct position
// (ascending)
func (l *LinkedList) PushOrdered(v *Vertex) *Element {
	l.lazyInit()
	if l.Len() == 0 {
		return l.PushFront(v)
	}
	back := l.Back()
	if back.Value.Distance < v.Distance {
		return l.insertValue(v, l.root.prev)
	}
	current := l.Front()
	for current.Value.Distance < v.Distance { //don't need to chack if current=back cause back already checked
		current = current.Next()
	}
	return l.insertValue(v, current.prev)
}

// insert inserts e after at, increments l.len, and returns e.
func (l *LinkedList) insert(e, at *Element) *Element {
	n := at.next
	at.next = e
	e.prev = at
	e.next = n
	n.prev = e
	e.list = l
	l.len++
	return e
}

// insertValue is a convenience wrapper for insert(&Element{Value: v}, at).
func (l *LinkedList) insertValue(v *Vertex, at *Element) *Element {
	return l.insert(&Element{Value: v}, at)
}

// remove removes e from its list, decrements l.len, and returns e.
func (l *LinkedList) remove(e *Element) *Element {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
	l.len--
	return e
}

// PushFront inserts a new element e with value v at the front of list l and returns e.
func (l *LinkedList) PushFront(v *Vertex) *Element {
	l.lazyInit()
	return l.insertValue(v, &l.root)
}
