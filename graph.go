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
	Verticies []Vertex
	Visiting  *linkedList
	mapping   map[string]int
	usingMap  bool
}

//GetMapped gets the key assosciated with the mapped int
func (g *Graph) GetMapped(a int) (string, error) {
	for k, v := range g.mapping {
		if v == a {
			return k, nil
		}
	}
	return "", errors.New(fmt.Sprint(a, " not found in mapping"))
}

//GetMapping gets the index associated with the specified key
func (g *Graph) GetMapping(a string) (int, error) {
	if b, ok := g.mapping[a]; ok {
		return b, nil
	}
	return -1, errors.New(fmt.Sprint(a, " not found in mapping"))
}

//AddVerticies adds the listed verticies to the graph
func (g *Graph) AddVerticies(v ...Vertex) {
	g.Verticies = append(g.Verticies, v...)
}

func (g Graph) validate() error {
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
func (g *Graph) setDefaults(Distance int64, BestNode int) {
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
		for key := 0; key < len(v.Arcs); key++ {
			fmt.Fprint(f, " ", key, ",", v.Arcs[key])
		}
		fmt.Fprint(f, "\n")
	}
	return nil
}
