package dijkstra

import "math"

// BestPaths contains the list of best solutions
type BestPaths []BestPath

// ShortestAll calculates all of the shortest paths from src to dest
func (g *Graph) ShortestAll(src, dest int) (BestPaths, error) {
	return g.evaluateAll(src, dest, true, listShortAuto)
}

// LongestAll calculates all the longest paths from src to dest
func (g *Graph) LongestAll(src, dest int) (BestPaths, error) {
	return g.evaluateAll(src, dest, false, listLongAuto)
}

func (g *Graph) evaluateAll(src, dest int, shortest bool, listOption int) (BestPaths, error) {
	return g.postSetupEvaluateAll(src, dest, shortest, listOption)
}

func (g *Graph) postSetupEvaluateAll(src, dest int, shortest bool, listOption int) (BestPaths, error) {
	// if err := g.vertexValid(src); err != nil {
	// 	return BestPaths{}, err
	// }
	// if err := g.vertexValid(dest); err != nil {
	// 	return BestPaths{}, err
	// }
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
			return a <= b
		}
		shouldSkip = func(currentDistance, best, storedDistance int64) bool {
			return currentDistance > best || currentDistance > storedDistance
		}
	}
	var visiting = g.forceList(listOption)
	distances := make([]int64, len(g.Verticies))
	bestVerticies := make([][]int, len(g.Verticies))
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
		for to, dist := range g.Verticies[current.id].arcs {
			if better(current.distance+dist, distances[to]) {
				if bestVerticies[current.id][0] == to && to != dest {
					return BestPaths{}, newErrLoop(current.id, to)
				}
				if (current.distance + dist) == distances[to] {
					bestVerticies[to] = append(bestVerticies[to], current.id)
				} else {
					distances[to] = current.distance + dist
					bestVerticies[to] = []int{current.id}
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
		return BestPaths{}, ErrNoPath
	}
	paths := bestPaths(bestVerticies, src, dest)
	result := make(BestPaths, len(paths))
	for i := range paths {
		result[i] = BestPath{Distance: distances[dest], Path: paths[i]}
	}
	return result, nil
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
