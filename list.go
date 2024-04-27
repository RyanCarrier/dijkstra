package dijkstra

import "sort"

type dijkstraList interface {
	PushOrdered(currentDistance)
	PopOrdered() currentDistance
	Len() int
}

//THE FOLLOWING FILE IS partially BASED FROM GO AUTHORS EDITED MINORLY AND LAZILY TO SUIT MY NEEDS
//https://golang.org/src/container/list/list.go?m=text
//AVOID USING MINE AS A TEMPLATE AS I REMOVED MOST SAFETIES (that's why they are
// all private now)

// element is an element of a linked list.
type element struct {
	next, prev *element
	Value      currentDistance
}

// linkedList represents a doubly linked list.
// The zero value for linkedList is an empty list ready to use.
type linkedList struct {
	root  element // sentinel list element, only &root, root.prev, and root.next are used
	len   int     // current list length excluding (this) sentinel element
	short bool
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
func (l *linkedList) PushOrdered(v currentDistance) {
	l.pushOrdered(v)
}

// Init initializes or clears list l.
func (l *linkedList) PopOrdered() currentDistance {
	if l.short {
		return l.popFront()
	}
	return l.popBack()
}

// Init initializes or clears list l.
func (l *linkedList) Len() int {
	return l.len
}

// Init initializes or clears list l.
func (l *linkedList) init(short bool) *linkedList {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	l.short = short
	return l
}

// front returns the first element of list l or nil.
func (l *linkedList) front() *element {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// popFront pops the Vertex off the front of the list
func (l *linkedList) popFront() currentDistance {
	e := l.front()
	l.remove(e)
	return e.Value
}

// popFront pops the Vertex off the front of the list
func (l *linkedList) popBack() currentDistance {
	e := l.back()
	l.remove(e)
	return e.Value
}

// back returns the last element of list l or nil.
func (l *linkedList) back() *element {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// pushOrdered pushes the value into the linked list in the correct position
// (ascending)
func (l *linkedList) pushOrdered(v currentDistance) *element {
	if l.len == 0 {
		return l.pushFront(v)
	}
	back := l.back()
	if back.Value.distance < v.distance {
		return l.insertValue(v, l.root.prev)
	}
	current := l.front()
	//don't need to chack if current=back cause back already checked
	for current.Value.distance < v.distance {
		if current.Value.id == v.id {
			if l.short {
				return current
			} else {
				old := current
				current = current.next
				l.remove(old)
			}
			// return current
		} else {
			current = current.next
		}
	}
	return l.insertValue(v, current.prev)
}

// insert inserts e after at, increments l.len, and returns e.
func (l *linkedList) insert(e, at *element) *element {
	n := at.next
	at.next = e
	e.prev = at
	e.next = n
	n.prev = e
	l.len++
	return e
}

// insertValue is a convenience wrapper for insert(&element{Value: v}, at).
func (l *linkedList) insertValue(v currentDistance, at *element) *element {
	return l.insert(&element{Value: v}, at)
}

// remove removes e from its list, decrements l.len, and returns e.
func (l *linkedList) remove(e *element) *element {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	l.len--
	return e
}

// pushFront inserts a new element e with value v at the front of list l and returns e.
func (l *linkedList) pushFront(v currentDistance) *element {
	return l.insertValue(v, &l.root)
}

// PriorityQueueNewShort creates a new priority queue for short solving
func priorityQueueNewShort() dijkstraList {
	return &priorityQueueWrapper{new(priorityQueueShort)}
}

// PriorityQueueNewLong creates a new priority queue for long solving
func priorityQueueNewLong() dijkstraList {
	return &priorityQueueWrapper{new(priorityQueueLong)}
}

type priorityQueueShort struct{ priorityQueueBase }
type priorityQueueLong struct{ priorityQueueBase }
type priorityQueueInterface interface {
	sort.Interface
	Push(x currentDistance)
	Pop() currentDistance
}
type priorityQueueWrapper struct{ priorityQueueInterface }

func (pq priorityQueueShort) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq.priorityQueueBase[i].distance < pq.priorityQueueBase[j].distance
}

func (pq priorityQueueLong) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq.priorityQueueBase[i].distance > pq.priorityQueueBase[j].distance
}

type priorityQueueBase []currentDistance

func (pq priorityQueueBase) Len() int { return len(pq) }

func (pq priorityQueueBase) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueueBase) Push(v currentDistance) {
	*pq = append(*pq, v)
}

func (pq *priorityQueueWrapper) PushOrdered(v currentDistance) {
	pq.Push(v)
	pq.up(pq.Len() - 1)
}

func (pq *priorityQueueWrapper) PopOrdered() currentDistance {
	n := pq.Len() - 1
	pq.Swap(0, n)
	pq.down(0, n)
	return pq.Pop()
}

func (pq *priorityQueueBase) Pop() currentDistance {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func (pq *priorityQueueWrapper) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !pq.Less(j, i) {
			break
		}
		pq.Swap(i, j)
		j = i
	}
}

func (pq *priorityQueueWrapper) down(i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && pq.Less(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !pq.Less(j, i) {
			break
		}
		pq.Swap(i, j)
		i = j
	}
	return i > i0
}
