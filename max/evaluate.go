package max

import (
	"errors"
	"math"
)

func (g *Graph) setup(src int, begin int64) {
	//Get a new list regardless
	g.visiting = newLinkedList()
	//Reset state
	g.visitedDest = false
	//Reset the best current value (worst so it gets overwritten)
	// and set the defaults *almost* as bad
	// set all best verticies to -1 (unused)

	g.setDefaults(int64(math.MinInt64)+2, -1)
	g.best = int64(math.MinInt64)

	//Set the distance of initial vertex 0
	g.Verticies[src].best = begin
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
	return BestPath{g.Verticies[dest].best, path}
}

func (g *Graph) finally(src, dest int) (BestPath, error) {
	if !g.visitedDest {
		return BestPath{}, errors.New("No path found")
	}
	return g.bestPath(src, dest), nil
}
