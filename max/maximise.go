package dijkstraMax

import (
	"errors"
	"fmt"
)

//Maximise the path to the destination
func (g *Graph) Maximise(src, dest int, Initial int64) (BestPath, error) {
	shortest := false
	//setup with minimum values
	g.setup(src, Initial)
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
			continue
		}
		//If the current distance is already worse than the best try another Vertex
		if shortest && current.best >= g.best { //} || (!shortest && current.best <= g.best) {
			continue
		}
		for v, dist := range current.arcs {
			//If the arc has better access, than the current best, update the Vertex being touched
			if (shortest && current.best+dist < g.Verticies[v].best) ||
				(!shortest && current.best+dist > g.Verticies[v].best) {
				if g.Verticies[v].bestVertex == current.ID && g.Verticies[v].ID != dest {
					//also only do this if we aren't checkout out the best distance again
					//This seems familiar 8^)
					return BestPath{}, errors.New(fmt.Sprint("loop detected from ", current.ID, " to ", v))
				}
				g.Verticies[v].best = current.best + dist
				g.Verticies[v].bestVertex = current.ID
				if v == dest {
					//If this is the destination update best, so we can stop looking at
					// useless Verticies
					g.best = current.best + dist
				}
				//Push this updated Vertex into the list to be evaluated, pushes in
				// sorted form
				g.visiting.pushOrdered(&g.Verticies[v])
			}
		}
	}
	return g.finally(src, dest)
}
