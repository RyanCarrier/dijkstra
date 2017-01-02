package dijkstra

import (
	"errors"
	"fmt"
)

//Graph contains all the graph details
type Graph struct {
	//slice of visited Vertex (ID's)
	Visited []bool
	//slice of all verticies available
	Verticies []Vertex
}

//Vertex is a single node in the network, contains it's ID, best distance (to
// itself from the src) and the weight to go to each other connected node (Vertex)
type Vertex struct {
	//ID of the Vertex
	ID int
	//Best Distance to the Vertex
	Distance int64
	//A set of all weights to the nodes in the map
	Arcs map[int]int64
}

func (g Graph) validate() error {
	if len(g.Verticies) != len(g.Visited) {
		return errors.New("Verticies and visited slice are not same length")
	}
	for _, v := range g.Verticies {
		for a := range v.Arcs {
			if a >= len(g.Verticies) || (g.Verticies[a].ID == 0 && a != 0) {
				return errors.New(fmt.Sprint("Vertex ", a, " referenced in arcs by Vertex ", v.ID))
			}
		}
	}
	return nil
}
