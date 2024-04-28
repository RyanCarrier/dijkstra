package dijkstra

import (
	"cmp"
	"math"
	"slices"
)

// BestPath contains the solution of the most optimal path
type BestPath[T any] struct {
	Distance uint64
	Path     []T
}
type BestPaths[T any] struct {
	Distance uint64
	Paths    [][]T
}

func (bps BestPaths[T]) SmallestPath() BestPath[T] {
	if len(bps.Paths) == 0 {
		return BestPath[T]{}
	}
	smallest := slices.MinFunc(bps.Paths, func(a, b []T) int {
		return cmp.Compare(len(a), len(b))
	})
	return BestPath[T]{
		Distance: bps.Distance,
		Path:     smallest,
	}
}

type currentDistance struct {
	id       int
	distance uint64
}

const (
	listShortAuto = iota
	listLongAuto
	listShortPQ
	listLongPQ
	listShortLL
	listLongLL
)

// Shortest calculates the shortest path from src to dest
func (g Graph) Shortest(src, dest int) (BestPath[int], error) {
	return g.evaluate(src, dest, true, listShortAuto)
}

// Longest calculates the longest path from src to dest
func (g Graph) Longest(src, dest int) (BestPath[int], error) {
	return g.evaluate(src, dest, false, listLongAuto)
}

// ShortestAll calculates all of the longest paths from src to dest
func (g Graph) ShortestAll(src, dest int) (BestPaths[int], error) {
	return g.evaluateAll(src, dest, true, listShortAuto)
}

// LongestAll calculates all of the longest paths from src to dest
func (g Graph) LongestAll(src, dest int) (BestPaths[int], error) {
	return g.evaluateAll(src, dest, false, listLongAuto)
}

func (g Graph) getList(i int) dijkstraList {
	switch i {
	case listShortAuto:
		//LL seems to be faster for less than 100 verticies
		if len(g.vertexArcs) < 100 {
			return g.getList(listShortLL)
		} else {
			return g.getList(listShortPQ)
		}
	case listLongAuto:
		//LL seems to be faster for less than 100 verticies
		if len(g.vertexArcs) < 100 {
			return g.getList(listLongLL)
		} else {
			return g.getList(listLongPQ)
		}
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

func (g Graph) evaluate(src, dest int, shortest bool, listOption int) (BestPath[int], error) {
	if err := g.vertexValid(src); err != nil {
		return BestPath[int]{}, err
	}
	if err := g.vertexValid(dest); err != nil {
		return BestPath[int]{}, err
	}
	var current currentDistance
	visitedDest := false
	var newDefault uint64
	var best uint64
	var better = func(a, b uint64) bool {
		return a >= b
	}
	var shouldSkip = func(currentDistance, best, storedDistance uint64) bool {
		return currentDistance < storedDistance
	}
	if shortest {
		newDefault = uint64(math.MaxUint64) - 2
		best = uint64(math.MaxUint64)
		better = func(a, b uint64) bool {
			return a < b
		}
		shouldSkip = func(currentDistance, best, storedDistance uint64) bool {
			return currentDistance >= best || currentDistance > storedDistance
		}
	}
	var visiting = g.getList(listOption)
	distances := make([]uint64, len(g.vertexArcs))
	bestVerticie := make([]int, len(g.vertexArcs))
	for i := range bestVerticie {
		bestVerticie[i] = -1
		distances[i] = newDefault
	}
	distances[src] = 0
	visiting.PushOrdered(currentDistance{src, distances[src]})
	for visiting.Len() > 0 {
		current = visiting.PopOrdered()
		if shouldSkip(current.distance, best, distances[current.id]) {
			continue
		}
		for to, dist := range g.vertexArcs[current.id] {
			if better(current.distance+dist, distances[to]) {
				if bestVerticie[current.id] == to && to != dest {
					return BestPath[int]{}, newErrLoop(current.id, to)
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
		return BestPath[int]{}, newErrNoPath(src, dest)
	}
	var path []int
	for c := dest; c != src; c = bestVerticie[c] {
		path = append(path, c)
	}
	path = append(path, src)
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return BestPath[int]{distances[dest], path}, nil
}

func (g *Graph) evaluateAll(src, dest int, shortest bool, listOption int) (BestPaths[int], error) {
	if err := g.vertexValid(src); err != nil {
		return BestPaths[int]{}, err
	}
	if err := g.vertexValid(dest); err != nil {
		return BestPaths[int]{}, err
	}
	var current currentDistance
	visitedDest := false
	var newDefault uint64
	var best uint64
	var better = func(a, b uint64) bool {
		return a >= b
	}
	var shouldSkip = func(currentDistance, best, storedDistance uint64) bool {
		return currentDistance < storedDistance
	}
	if shortest {
		newDefault = uint64(math.MaxUint64) - 2
		best = uint64(math.MaxUint64)
		better = func(a, b uint64) bool {
			return a <= b
		}
		shouldSkip = func(currentDistance, best, storedDistance uint64) bool {
			return currentDistance > best || currentDistance > storedDistance
		}
	}
	var visiting = g.getList(listOption)
	distances := make([]uint64, len(g.vertexArcs))
	bestVerticies := make([][]int, len(g.vertexArcs))
	for i := range bestVerticies {
		bestVerticies[i] = []int{-1}
		distances[i] = newDefault
	}
	distances[src] = 0
	visiting.PushOrdered(currentDistance{src, distances[src]})
	for visiting.Len() > 0 {
		current = visiting.PopOrdered()
		if shouldSkip(current.distance, best, distances[current.id]) {
			continue
		}
		for to, dist := range g.vertexArcs[current.id] {
			if better(current.distance+dist, distances[to]) {
				if bestVerticies[current.id][0] == to && to != dest {
					return BestPaths[int]{}, newErrLoop(current.id, to)
				}
				if (current.distance + dist) == distances[to] {
					bestVerticies[to] = append(bestVerticies[to], current.id)
				} else {
					distances[to] = current.distance + dist
					bestVerticies[to][0] = current.id
					if to == dest {
						visitedDest = true
						best = distances[to]
						continue // Do not push if dest
					}
					visiting.PushOrdered(currentDistance{to, distances[to]})
				}
			}
		}
	}
	if !visitedDest {
		return BestPaths[int]{}, newErrNoPath(src, dest)
	}
	return BestPaths[int]{
		Distance: distances[dest],
		Paths:    bestPaths(bestVerticies, src, dest),
	}, nil
}

func bestPaths(bestVerticies [][]int, src, dest int) [][]int {
	paths := visitPath(bestVerticies, src, dest, dest)
	best := [][]int{}
	//reverse order of paths
	for indexPaths := range paths {
		for i, j := 0, len(paths[indexPaths])-1; i < j; i, j = i+1, j-1 {
			paths[indexPaths][i], paths[indexPaths][j] = paths[indexPaths][j], paths[indexPaths][i]
		}
		best = append(best, paths[indexPaths])
	}
	return best
}

// visitPath is a recursive function that will visit all the bestVerticies of a Vertex
func visitPath(bestVerticies [][]int, src, dest, currentVerticie int) [][]int {
	if currentVerticie == src {
		return [][]int{
			{currentVerticie},
		}
	}
	paths := [][]int{}
	for _, vertex := range bestVerticies[currentVerticie] {
		subPaths := visitPath(bestVerticies, src, dest, vertex)
		for i := range subPaths {
			//just note that we are walking in reverse order
			subPathTilHere := append([]int{currentVerticie}, subPaths[i]...)
			paths = append(paths, subPathTilHere)
		}
	}
	return paths
}
