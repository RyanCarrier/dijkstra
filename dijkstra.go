package dijkstra

import (
	"errors"
	"fmt"
	"math"
)

//Shortest calculates the shortest path from src to dest
func (g *Graph) Shortest(src, dest int) (BestPath, error) {
	g.Visiting = NewList()
	g.Visiting.PushFront(&g.Verticies[src])
	best := int64(math.MaxInt64)
	var current *Vertex
	for g.Visiting.Len() > 0 {
		current = g.Visiting.PopFront()
		if !g.Visited[current.ID] || current.Distance > best {
			continue
		}
		g.Visited[current.ID] = true
		for v, dist := range current.Arcs {
			if current.Distance+dist < g.Verticies[v].Distance {
				if g.Verticies[v].BestVertex == current.ID {
					return BestPath{},
						errors.New(fmt.Sprint(ErrLoopDetected.Error(), "From node '", current.ID, "' to node '", "'"))
				}
				g.Verticies[v].Distance = current.Distance + dist
				g.Verticies[v].BestVertex = current.ID
				g.Visiting.PushOrdered(&g.Verticies[v])
			}
		}
	}
	if !g.Visited[dest] {
		return BestPath{}, ErrNoPath
	}
	return g.bestPath(src, dest), nil
}

func (g *Graph) bestPath(src, dest int) BestPath {
	var path []int
	for c := g.Verticies[dest]; c.ID != src; c = g.Verticies[c.BestVertex] {
		path = append(path, c.ID)
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return BestPath{g.Verticies[dest].Distance, path}
}

//Longest calculates the longest path from src to dest
func (g *Graph) Longest(src, dest int) (BestPath, error) {
	return BestPath{}, ErrNoPath
}

func (g *Graph) addVisit(v *Vertex, Ascending bool) {

}

//BestPath contains the solution of the most optimal path
type BestPath struct {
	Distance int64
	Path     []int
}
