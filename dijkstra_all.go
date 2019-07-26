package dijkstra

//ShortestAll calculates all of the shortest paths from src to dest
func (g *Graph) ShortestAll(src, dest int) (BestPaths, error) {
	return g.evaluateAll(src, dest, true)
}

//LongestAll calculates all the longest paths from src to dest
func (g *Graph) LongestAll(src, dest int) (BestPaths, error) {
	return g.evaluateAll(src, dest, false)
}

func (g *Graph) evaluateAll(src, dest int, shortest bool) (BestPaths, error) {
	//Setup graph
	g.setup(shortest, src, -1)
	return g.postSetupEvaluateAll(src, dest, shortest)
}

func (g *Graph) postSetupEvaluateAll(src, dest int, shortest bool) (BestPaths, error) {
	var current *Vertex
	oldCurrent := -1
	for g.visiting.Len() > 0 {
		//Visit the current lowest distanced Vertex
		current = g.visiting.PopOrdered()
		if oldCurrent == current.ID {
			continue
		}
		oldCurrent = current.ID
		//If the current distance is already worse than the best try another Vertex
		if shortest && current.distance > g.best {
			continue
		}
		for v, dist := range current.arcs {
			//If the arc has better access, than the current best, update the Vertex being touched
			if (shortest && current.distance+dist < g.Verticies[v].distance) ||
				(!shortest && current.distance+dist > g.Verticies[v].distance) ||
				(current.distance+dist == g.Verticies[v].distance && !g.Verticies[v].containsBest(current.ID)) {
				//if g.Verticies[v].bestVertex == current.ID && g.Verticies[v].ID != dest {
				if current.containsBest(v) && g.Verticies[v].ID != dest {
					//also only do this if we aren't checkout out the best distance again
					//This seems familiar 8^)
					return BestPaths{}, newErrLoop(current.ID, v)
				}
				if current.distance+dist == g.Verticies[v].distance {
					//At this point we know it's not in the list due to initial check
					g.Verticies[v].bestVerticies = append(g.Verticies[v].bestVerticies, current.ID)
				} else {
					g.Verticies[v].distance = current.distance + dist
					g.Verticies[v].bestVerticies = []int{current.ID}
				}
				if v == dest {
					g.visitedDest = true
					g.best = current.distance + dist
					continue
					//If this is the destination update best, so we can stop looking at
					// useless Verticies
				}
				//Push this updated Vertex into the list to be evaluated, pushes in
				// sorted form
				g.visiting.PushOrdered(&g.Verticies[v])
			}
		}
	}
	if !g.visitedDest {
		return BestPaths{}, ErrNoPath
	}
	return g.bestPaths(src, dest), nil
}

func (g *Graph) bestPaths(src, dest int) BestPaths {
	paths := g.visitPath(src, dest, dest, [][]int{})
	best := BestPaths{}
	for indexPaths := range paths {
		for i, j := 0, len(paths[indexPaths])-1; i < j; i, j = i+1, j-1 {
			paths[indexPaths][i], paths[indexPaths][j] = paths[indexPaths][j], paths[indexPaths][i]
		}
		best = append(best, BestPath{g.Verticies[dest].distance, paths[indexPaths]})
	}

	return best
}

func (g *Graph) visitPath(src, dest, currentNode int, paths [][]int) [][]int {
	if currentNode == src {
		paths[len(paths)-1] = append(paths[len(paths)-1], currentNode)
		return paths
	}
	for _, vertex := range g.Verticies[currentNode].bestVerticies {
		if currentNode == dest {
			paths = append(paths, []int{dest})
		} else {
			paths[len(paths)-1] = append(paths[len(paths)-1], currentNode)
		}
		paths = g.visitPath(src, dest, vertex, paths)
	}
	return paths
}
