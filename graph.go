package dijkstra

import (
	"errors"
	"fmt"
	"math"
)

//Graph contains all the graph details
type Graph struct {
	//slice of visited Vertex (ID's)
	Visited []bool
	//slice of all verticies available
	Verticies []Vertex

	Visiting *List
}

func (g Graph) validate() error {
	if len(g.Verticies) != len(g.Visited) {
		return errors.New("Verticies and visited slice are not same length")
	}
	for _, v := range g.Verticies {
		for a := range v.Arcs {
			if a >= len(g.Verticies) || (g.Verticies[a].ID == 0 && a != 0) {
				fmt.Printf("%+v", g)
				return errors.New(fmt.Sprint("Graph validation error;", "Vertex ", a, " referenced in arcs by Vertex ", v.ID))
			}
		}
	}
	return nil
}

//SetMaxDistances sets all Vertex distances to the max of int64
func (g *Graph) SetMaxDistances() {
	g.setAllDistances(math.MaxInt64)
}

//SetMinDistances sets all Vertex distances to the min of int64
func (g *Graph) SetMinDistances() {
	g.setAllDistances(math.MinInt64)
}

func (g *Graph) setAllDistances(max int64) {
	for i := range g.Verticies {
		g.Verticies[i].Distance = max
	}
}
