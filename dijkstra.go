package dijkstra

import (
	"math"
)

// BestPath contains the solution of the most optimal path
type BestPath struct {
	Distance int64
	Path     []int
}

// Shortest calculates the shortest path from src to dest
func (g *Graph) Shortest(src, dest int) (BestPath, error) {
	return g.evaluate(src, dest, true)
}

// Longest calculates the longest path from src to dest
func (g *Graph) Longest(src, dest int) (BestPath, error) {
	return g.evaluate(src, dest, false)
}

func (g *Graph) setup(shortest bool, src int, list int) {
	//-1 auto list
	//Get a new list regardless
	if list >= 0 {
		g.forceList(list)
	} else if shortest {
		g.forceList(-1)
	} else {
		g.forceList(-2)
	}
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
	g.visiting.PushOrdered(&g.Verticies[src])
}

func (g *Graph) forceList(i int) {
	//-2 long auto
	//-1 short auto
	//0 short pq
	//1 long pq
	//2 short ll
	//3 long ll
	switch i {
	case -2:
		if len(g.Verticies) < 800 {
			g.forceList(2)
		} else {
			g.forceList(0)
		}
	case -1:
		if len(g.Verticies) < 800 {
			g.forceList(3)
		} else {
			g.forceList(1)
		}
	case 0:
		g.visiting = priorityQueueNewShort()
	case 1:
		g.visiting = priorityQueueNewLong()
	case 2:
		g.visiting = linkedListNewShort()
	case 3:
		g.visiting = linkedListNewLong()
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

func (g *Graph) evaluate(src, dest int, shortest bool) (BestPath, error) {
	if g.running {
		return BestPath{}, ErrAlreadyCalculating
	}
	g.running = true
	defer func() { g.running = false }()
	//Setup graph
	g.setup(shortest, src, -1)
	return g.postSetupEvaluate(src, dest, shortest)
}

// ShortestSafe calculates the shortest path from src to dest with thread safety
//
//	(slightly slower than regular Shortest)
func (g *Graph) ShortestSafe(src, dest int) (BestPath, error) {
	return g.evaluate(src, dest, true)
}

// LongestSafe calculates the longest path from src to dest with thread safety
//
//	(slightly slower than regular Longest)
func (g *Graph) LongestSafe(src, dest int) (BestPath, error) {
	return g.evaluate(src, dest, false)
}

func (g *Graph) evaluateSafe(src, dest int, shortest bool) (BestPath, error) {
	var current *Vertex
	var best int64
	var newDefault int64
	var visiting dijkstraList
	visitedDest := false
	oldCurrent := -1
	verticies := make([]Vertex, len(g.Verticies))
	copy(verticies, g.Verticies)
	if shortest {
		if len(verticies) < 800 {
			visiting = linkedListNewLong()
		} else {
			visiting = priorityQueueNewLong()
		}
		newDefault = int64(math.MaxInt64) - 2
		best = int64(math.MaxInt64)
	} else {
		if len(verticies) < 800 {
			visiting = linkedListNewShort()
		} else {
			visiting = priorityQueueNewShort()
		}
		newDefault = int64(math.MinInt64) + 2
		best = int64(math.MinInt64)
	}
	for i := range verticies {
		verticies[i].bestVerticies = []int{-1}
		verticies[i].distance = newDefault
	}
	verticies[src].distance = 0
	visiting.PushOrdered(&verticies[src])

	for visiting.Len() > 0 {
		//Visit the current lowest distanced Vertex
		//TODO WTF
		current = visiting.PopOrdered()
		if oldCurrent == current.ID {
			continue
		}
		oldCurrent = current.ID
		//If the current distance is already worse than the best try another Vertex
		if shortest && current.distance >= best {
			continue
		}
		for v, dist := range current.arcs {
			//If the arc has better access, than the current best, update the Vertex being touched
			if (shortest && current.distance+dist < verticies[v].distance) ||
				(!shortest && current.distance+dist > verticies[v].distance) {
				if current.bestVerticies[0] == v && verticies[v].ID != dest {
					//also only do this if we aren't checkout out the best distance again
					//This seems familiar 8^)
					return BestPath{}, newErrLoop(current.ID, v)
				}
				verticies[v].distance = current.distance + dist
				verticies[v].bestVerticies[0] = current.ID
				if v == dest {
					//If this is the destination update best, so we can stop looking at
					// useless Verticies
					best = current.distance + dist
					visitedDest = true
					continue // Do not push if dest
				}
				//Push this updated Vertex into the list to be evaluated, pushes in
				// sorted form
				visiting.PushOrdered(&verticies[v])
			}
		}
	}
	if !visitedDest {
		return BestPath{}, ErrNoPath
	}
	var path []int
	for c := verticies[dest]; c.ID != src; c = verticies[c.bestVerticies[0]] {
		path = append(path, c.ID)
	}
	path = append(path, src)
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return BestPath{verticies[dest].distance, path}, nil
}
func (g *Graph) postSetupEvaluate(src, dest int, shortest bool) (BestPath, error) {
	var current *Vertex
	oldCurrent := -1
	for g.visiting.Len() > 0 {
		//Visit the current lowest distanced Vertex
		//TODO WTF
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
				if current.bestVerticies[0] == v && g.Verticies[v].ID != dest {
					//also only do this if we aren't checkout out the best distance again
					//This seems familiar 8^)
					return BestPath{}, newErrLoop(current.ID, v)
				}
				g.Verticies[v].distance = current.distance + dist
				g.Verticies[v].bestVerticies[0] = current.ID
				if v == dest {
					//If this is the destination update best, so we can stop looking at
					// useless Verticies
					g.best = current.distance + dist
					g.visitedDest = true
					continue // Do not push if dest
				}
				//Push this updated Vertex into the list to be evaluated, pushes in
				// sorted form
				g.visiting.PushOrdered(&g.Verticies[v])
			}
		}
	}
	return g.finally(src, dest)
}

func (g *Graph) finally(src, dest int) (BestPath, error) {
	if !g.visitedDest {
		return BestPath{}, ErrNoPath
	}
	return g.bestPath(src, dest), nil
}
