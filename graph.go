package dijkstra

import (
	"errors"
	"fmt"
	"os"
)

//Graph contains all the graph details
type Graph struct {
	best        int64
	visitedDest bool
	//slice of all verticies available
	Verticies       []Vertex
	Visiting        *linkedList
	mapping         map[string]int
	usingMap        bool
	highestMapIndex int
}

func (g *Graph) GetVertex(i int) (*Vertex, error) {
	if i >= len(g.Verticies) {
		return nil, errors.New("Vertex not found")
	}
	return &g.Verticies[i], nil
}

func (g Graph) validate() error {
	for _, v := range g.Verticies {
		for a := range v.arcs {
			if a >= len(g.Verticies) || (g.Verticies[a].ID == 0 && a != 0) {
				fmt.Printf("%+v", g)
				return errors.New(fmt.Sprint("Graph validation error;", "Vertex ", a, " referenced in arcs by Vertex ", v.ID))
			}
		}
	}
	return nil
}

//SetDefaults sets the distance and best node to that specified
func (g *Graph) setDefaults(Distance int64, BestNode int) {
	for i := range g.Verticies {
		g.Verticies[i].bestVertex = BestNode
		g.Verticies[i].distance = Distance
	}
}

//ExportToFile exports the verticies to file
func (g Graph) ExportToFile(filename string) error {
	//TODO ADD MAP STUFF
	if _, err := os.Stat(filename); err == nil {
		os.Remove(filename)
	}
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, v := range g.Verticies {
		fmt.Fprint(f, v.ID)
		for key := 0; key < len(v.arcs); key++ {
			fmt.Fprint(f, " ", key, ",", v.arcs[key])
		}
		fmt.Fprint(f, "\n")
	}
	return nil
}
