package dijkstra

import (
	"sort"
)

type dijkstraList interface {
	PushOrdered(*Vertex)
	PopOrdered() *Vertex
	Len() int
}

//PriorityQueueNewShort creates a new priority queue for short solving
func priorityQueueNewShort() dijkstraList {
	l := &priorityQueueWrapper{new(priorityQueueShort)}
	n := l.Len()
	for i := n/2 - 1; i >= 0; i-- {
		l.down(i, n)
	}
	return l
}

//PriorityQueueNewLong creates a new priority queue for long solving
func priorityQueueNewLong() dijkstraList {
	l := &priorityQueueWrapper{new(priorityQueueLong)}
	n := l.Len()
	for i := n/2 - 1; i >= 0; i-- {
		l.down(i, n)
	}
	return l
}

type priorityQueueShort struct{ priorityQueueBase }
type priorityQueueLong struct{ priorityQueueBase }
type priorityQueueInterface interface {
	sort.Interface
	Push(x *Vertex)
	Pop() *Vertex
}
type priorityQueueWrapper struct{ priorityQueueInterface }

func (pq priorityQueueShort) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq.priorityQueueBase[i].value.distance > pq.priorityQueueBase[j].value.distance
}

func (pq priorityQueueLong) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq.priorityQueueBase[i].value.distance < pq.priorityQueueBase[j].value.distance
}

// An Item is something we manage in a priority queue.
type Item struct {
	value *Vertex // The value of the item; arbitrary.
}
type priorityQueueBase []*Item

func (pq priorityQueueBase) Len() int { return len(pq) }

func (pq priorityQueueBase) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueueBase) Push(v *Vertex) {
	/*	id := v.(*Vertex).ID
		for _, v := range *pq {
			if v.value.ID == id {
				return
			}
		}*/
	*pq = append(*pq, &Item{v})
}

func (pq *priorityQueueWrapper) PushOrdered(v *Vertex) {
	pq.Push(v)
	pq.up(pq.Len() - 1)
}

func (pq *priorityQueueWrapper) PopOrdered() *Vertex {
	n := pq.Len() - 1
	pq.Swap(0, n)
	pq.down(0, n)
	return pq.Pop()
}

func (pq *priorityQueueBase) Pop() *Vertex {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item.value
}

////////////////////////HEAP CODE/////////////////// Copied from container/heap

// Fix re-establishes the heap ordering after the element at index i has changed its value.
// Changing the value of the element at index i and then calling Fix is equivalent to,
// but less expensive than, calling Remove(h, i) followed by a Push of the new value.
// The complexity is O(log(n)) where n = h.Len().
func (pq *priorityQueueWrapper) fix(i int) {
	if !pq.down(i, pq.Len()) {
		pq.up(i)
	}
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
