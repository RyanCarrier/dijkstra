package dijkstra

import (
	"errors"
	"fmt"
	"os"
)

//Graph contains all the graph details
type Graph struct {
	//slice of visited Vertex (ID's)
	Visited []bool
	//slice of all verticies available
	Verticies []Vertex

	Visiting *List
}

//AddVerticies adds the listed verticies to the graph
func (g *Graph) AddVerticies(v ...Vertex) {
	g.Verticies = append(g.Verticies, v...)
	g.Visited = append(g.Visited, make([]bool, len(v))...)
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

//SetDefaults sets the distance and best node to that specified
func (g *Graph) SetDefaults(Distance int64, BestNode int) {
	for i := range g.Verticies {
		g.Verticies[i].BestVertex = BestNode
		g.Verticies[i].Distance = Distance
	}
}

//ExportToFile exports the verticies to file
func (g Graph) ExportToFile(filename string) error {
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
		for key, val := range v.Arcs {
			fmt.Fprint(f, " ", key, ",", val)
		}
		fmt.Fprint(f, "\n")
	}
	return nil
}
