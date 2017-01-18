package max

import (
	"errors"
	"fmt"
)

//Maximise the path to the destination
func (g *Graph) Maximise(src, dest int, Initial int64) (BestPath, error) {
	//setup with minimum values
	g.setup(src, Initial)
	var current *Vertex
	for g.visiting.len > 0 {
		//Visit the current lowest distanced Vertex
		current = g.visiting.popFront()

		//If we have hit the destination set the flag, cheaper than checking it's
		// distance change at the end
		if current.ID == dest {
			g.visitedDest = true
			continue
		}
		//don't cut short because a 'shorter' immediate path might end up having greater final
		for v, dist := range current.arcs {
			//If the arc has better access, than the current best, update the Vertex being touched
			//TODO ADD CONVERSION INSTEAD OF DIST
			vNewD := current.Evaluate(g.Verticies[v], dist)
			if vNewD > g.Verticies[v].best {
				if g.Verticies[v].bestVertex == current.ID && g.Verticies[v].ID != dest {
					//also only do this if we aren't checkout out the best distance again
					//This seems familiar 8^)
					return BestPath{}, errors.New(fmt.Sprint("loop detected from ", current.ID, " to ", v))
				}
				g.Verticies[v].best = vNewD
				g.Verticies[v].bestVertex = current.ID
				if v == dest {
					//If this is the destination update best, so we can stop looking at
					// useless Verticies
					g.best = vNewD
				}
				//Push this updated Vertex into the list to be evaluated, pushes in
				// sorted form
				g.visiting.push(&g.Verticies[v])
			}
		}
	}
	return g.finally(src, dest)
}
