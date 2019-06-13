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
	solutions := BestPaths{}
	for g.visiting.Len() > 0 {
		//Visit the current lowest distanced Vertex
		current = g.visiting.PopOrdered()
		if oldCurrent == current.ID {
			continue
		}
		oldCurrent = current.ID
		//If the current distance is already worse than the best try another Vertex
		if shortest && current.distance >= g.best {
			continue
		}
		for v, dist := range current.arcs {
			//If the arc has better access, than the current best, update the Vertex being touched
			if (shortest && current.distance+dist < g.Verticies[v].distance) ||
				(!shortest && current.distance+dist > g.Verticies[v].distance) {
				//if g.Verticies[v].bestVertex == current.ID && g.Verticies[v].ID != dest {
				if current.bestVertex == v && g.Verticies[v].ID != dest {
					//also only do this if we aren't checkout out the best distance again
					//This seems familiar 8^)
					return BestPaths{}, newErrLoop(current.ID, v)
				}
				g.Verticies[v].distance = current.distance + dist
				g.Verticies[v].bestVertex = current.ID
				if v == dest {
					g.visitedDest = true
					if (shortest && current.distance+dist <= g.best) || (!shortest && current.distance+dist >= g.best) {
						if current.distance+dist != g.best {
							g.best = current.distance + dist
							solutions = BestPaths{}
						}
						solutions = append(solutions, g.bestPath(src, dest))
					}
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
	return solutions, nil
}
