package dijkstra

import (
	"errors"
	"fmt"
	"math"
)

//Shortest calculates the shortest path from src to dest
func (g *Graph) Shortest(src, dest int) (BestPath, error) {
	visitedDest := false
	g.Visiting = NewList()
	g.SetDefaults(int64(math.MaxInt64)-1, -1)
	g.Verticies[src].Distance = 0
	g.Visiting.PushFront(&g.Verticies[src])
	best := int64(math.MaxInt64)
	var current *Vertex
	for g.Visiting.Len() > 0 {
		current = g.Visiting.PopFront()
		if current.ID == dest {
			visitedDest = true
		}
		if current.Distance >= best {
			continue
		}
		for v, dist := range current.Arcs {
			if current.Distance+dist < g.Verticies[v].Distance {
				if g.Verticies[v].BestVertex == current.ID {
					//This seems familiar 8^)
					return BestPath{},
						errors.New(fmt.Sprint(ErrLoopDetected.Error(), "From node '", current.ID, "' to node '", "'"))
				}
				g.Verticies[v].Distance = current.Distance + dist
				g.Verticies[v].BestVertex = current.ID
				if v == dest {
					best = current.Distance + dist
				}
				g.Visiting.PushOrdered(&g.Verticies[v])
			}
		}
	}
	if !visitedDest {
		return BestPath{}, ErrNoPath
	}
	return g.bestPath(src, dest), nil
}

func (g *Graph) bestPath(src, dest int) BestPath {
	var path []int
	for c := g.Verticies[dest]; c.ID != src; c = g.Verticies[c.BestVertex] {
		path = append(path, c.ID)
	}
	path = append(path, src)
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return BestPath{g.Verticies[dest].Distance, path}
}

//Longest calculates the longest path from src to dest
func (g *Graph) Longest(src, dest int) (BestPath, error) {
	return BestPath{}, ErrNoPath
}

//BestPath contains the solution of the most optimal path
type BestPath struct {
	Distance int64
	Path     []int
}
