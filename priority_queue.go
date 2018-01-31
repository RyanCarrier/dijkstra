package dijkstra

import "container/heap"

type dijkstraList interface {
	PushOrdered(interface{})
	PopOrdered() interface{}
	Len() int
}

//PriorityQueueNewShort creates a new priority queue for short solving
func priorityQueueNewShort() dijkstraList {
	l := &priorityQueueWrapper{new(priorityQueueShort)}
	heap.Init(l)
	return l
}

//PriorityQueueNewLong creates a new priority queue for long solving
func priorityQueueNewLong() dijkstraList {
	l := &priorityQueueWrapper{new(priorityQueueLong)}
	heap.Init(l)
	return l
}

type priorityQueueShort struct{ priorityQueueBase }
type priorityQueueLong struct{ priorityQueueBase }
type priorityQueueInterface interface {
	Push(interface{})
	Pop() interface{}
	Less(i, j int) bool
	Swap(i, j int)
	Len() int
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

func (pq *priorityQueueBase) Push(v interface{}) {
	/*	id := v.(*Vertex).ID
		for _, v := range *pq {
			if v.value.ID == id {
				return
			}
		}*/
	*pq = append(*pq, &Item{v.(*Vertex)})
}

func (pq *priorityQueueWrapper) PushOrdered(v interface{}) {
	heap.Push(pq, v)
}

func (pq *priorityQueueWrapper) PopOrdered() interface{} {
	return heap.Pop(pq)
}

func (pq *priorityQueueBase) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item.value
}
