package dijkstra

import "math"

//Shortest calculates the shortest path from src to dest
func (g *Graph) Shortest(src, dest int) (BestPath, error) {
	return g.evaluate(src, dest, true)
}

//Longest calculates the longest path from src to dest
func (g *Graph) Longest(src, dest int) (BestPath, error) {
	return g.evaluate(src, dest, false)
}

func (g *Graph) finally(src, dest int) (BestPath, error) {
	if !g.visitedDest {
		return BestPath{}, ErrNoPath
	}
	return g.bestPath(src, dest), nil
}

func (g *Graph) setup(shortest bool, src int) {
	//Get a new list regardless
	g.visiting = newLinkedList()
	//Reset state
	g.visitedDest = false
	//Reset the best current value (worst so it gets overwritten)
	// and set the defaults *almost* as bad
	// set all best verticies to -1 (unused)
	if shortest {
		g.setDefaults(int64(math.MaxInt64)-2, -1)
		g.best = int64(math.MaxInt64)
	} else {
		g.setDefaults(int64(math.MinInt64)+2, -1)
		g.best = int64(math.MinInt64)
	}
	//Set the distance of initial vertex 0
	g.Verticies[src].distance = 0
	//Add the source vertex to the list
	g.visiting.pushFront(&g.Verticies[src])
}

func (g *Graph) bestPath(src, dest int) BestPath {
	var path []int
	for c := g.Verticies[dest]; c.ID != src; c = g.Verticies[c.bestVertex] {
		path = append(path, c.ID)
	}
	path = append(path, src)
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return BestPath{g.Verticies[dest].distance, path}
}

func (g *Graph) evaluate(src, dest int, shortest bool) (BestPath, error) {
	//Setup graph
	g.setup(shortest, src)
	var current *Vertex
	for g.visiting.len > 0 {
		//Visit the current lowest distanced Vertex
		if shortest {
			current = g.visiting.popFront()
		} else {
			current = g.visiting.popBack()
		}
		//If we have hit the destination set the flag, cheaper than checking it's
		// distance change at the end
		if current.ID == dest {
			g.visitedDest = true
		}
		//If the current distance is already worse than the best try another Vertex
		if shortest && current.distance >= g.best { //} || (!shortest && current.distance <= g.best) {
			continue
		}
		for v, dist := range current.arcs {
			//If the arc has better access, than the current best, update the Vertex being touched
			if (shortest && current.distance+dist < g.Verticies[v].distance) ||
				(!shortest && current.distance+dist > g.Verticies[v].distance) {
				if g.Verticies[v].bestVertex == current.ID {
					//This seems familiar 8^)
					return BestPath{}, newErrLoop(current.ID, v)
				}
				g.Verticies[v].distance = current.distance + dist
				g.Verticies[v].bestVertex = current.ID
				if v == dest {
					//If this is the destination update best, so we can stop looking at
					// useless Verticies
					g.best = current.distance + dist
				}
				//Push this updated Vertex into the list to be evaluated, pushes in
				// sorted form
				g.visiting.pushOrdered(&g.Verticies[v])
			}
		}
	}
	return g.finally(src, dest)
}

//BestPath contains the solution of the most optimal path
type BestPath struct {
	Distance int64
	Path     []int
}
