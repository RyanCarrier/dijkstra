package dijkstra

import (
	"math"
)

const (
	listShortAuto = iota
	listLongAuto
	listShortPQ
	listLongPQ
	listShortLL
	listLongLL
)

// BestPath contains the solution of the most optimal path
type BestPath struct {
	Distance int64
	Path     []int
}

// Shortest calculates the shortest path from src to dest
func (g *Graph) Shortest(src, dest int) (BestPath, error) {
	return g.evaluate(src, dest, true, listShortAuto)
}

// Longest calculates the longest path from src to dest
func (g *Graph) Longest(src, dest int) (BestPath, error) {
	return g.evaluate(src, dest, false, listLongAuto)
}

func (g *Graph) forceList(i int) dijkstraList {
	switch i {
	case listShortAuto:
		//LL seems to be faster for less than 100 verticies
		if len(g.Verticies) < 100 {
			return g.forceList(listShortLL)
		}
		return g.forceList(listShortPQ)
	case listLongAuto:
		//LL seems to be faster for less than 100 verticies
		if len(g.Verticies) < 100 {
			return g.forceList(listLongLL)
		}
		return g.forceList(listLongPQ)
	case listShortPQ:
		return priorityQueueNewShort()
	case listLongPQ:
		return priorityQueueNewLong()
	case listShortLL:
		return linkedListNewShort()
	case listLongLL:
		return linkedListNewLong()
	default:
		panic(i)
	}
}

func (g *Graph) bestPath(src, dest int) BestPath {
	var path []int
	for c := g.Verticies[dest]; c.ID != src; c = g.Verticies[c.bestVerticies[0]] {
		path = append(path, c.ID)
	}
	path = append(path, src)
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return BestPath{g.Verticies[dest].distance, path}
}

func (g *Graph) evaluate(src, dest int, shortest bool, listOption int) (BestPath, error) {
	return g.evaluateSafe(src, dest, shortest, listOption)
}

// ShortestSafe calculates the shortest path from src to dest with thread safety
//
//	(slightly slower than regular Shortest)
func (g *Graph) ShortestSafe(src, dest int) (BestPath, error) {
	return g.evaluateSafe(src, dest, true, listShortAuto)
}

// LongestSafe calculates the longest path from src to dest with thread safety
//
//	(slightly slower than regular Longest)
func (g *Graph) LongestSafe(src, dest int) (BestPath, error) {
	return g.evaluateSafe(src, dest, false, listLongAuto)
}

func (g Graph) vertexValid(v int) error {
	if v >= len(g.Verticies) {
		return newErrVertexNotFound(v)
	}
	return nil
}

type currentDistance struct {
	id       int
	distance int64
}

func (g *Graph) evaluateSafe(src, dest int, shortest bool, listOption int) (BestPath, error) {
	var current currentDistance
	visitedDest := false
	var newDefault int64
	var best int64
	var better = func(a, b int64) bool {
		return a >= b
	}
	var shouldSkip = func(currentDistance, best, storedDistance int64) bool {
		return currentDistance < storedDistance
	}
	if shortest {
		newDefault = math.MaxInt64
		best = math.MaxInt64
		better = func(a, b int64) bool {
			return a < b
		}
		shouldSkip = func(currentDistance, best, storedDistance int64) bool {
			return currentDistance >= best || currentDistance > storedDistance
		}
	}
	var visiting = g.forceList(listOption)
	distances := make([]int64, len(g.Verticies))
	bestVerticie := make([]int, len(g.Verticies))
	for i := range bestVerticie {
		bestVerticie[i] = -1
		distances[i] = newDefault
	}
	distances[src] = 0
	initial := currentDistance{src, distances[src]}
	visiting.PushOrdered(initial)
	for visiting.Len() > 0 {
		current = visiting.PopOrdered()
		if shouldSkip(current.distance, best, distances[current.id]) {
			continue
		}
		for to, dist := range g.Verticies[current.id].arcs {
			if better(current.distance+dist, distances[to]) {
				if bestVerticie[current.id] == to && to != dest {
					return BestPath{}, newErrLoop(current.id, to)
				}
				distances[to] = current.distance + dist
				bestVerticie[to] = current.id
				if to == dest {
					visitedDest = true
					best = distances[to]
					continue // Do not push if dest
				}
				visiting.PushOrdered(currentDistance{to, distances[to]})
			}
		}
	}
	if !visitedDest {
		return BestPath{}, ErrNoPath
	}
	var path []int
	for c := dest; c != src; c = bestVerticie[c] {
		path = append(path, c)
	}
	path = append(path, src)
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return BestPath{distances[dest], path}, nil
}
