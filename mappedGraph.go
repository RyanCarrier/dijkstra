package dijkstra

import (
	"errors"
	"fmt"
)

//GetMapped gets the key assosciated with the mapped int
func (g *Graph) GetMapped(a int) (string, error) {
	if !g.usingMap {
		return "", errors.New("Map is not being used/initialised")
	}
	for k, v := range g.mapping {
		if v == a {
			return k, nil
		}
	}
	return "", errors.New(fmt.Sprint(a, " not found in mapping"))
}

//GetMapping gets the index associated with the specified key
func (g *Graph) GetMapping(a string) (int, error) {
	if !g.usingMap {
		return -1, errors.New("Map is not being used/initialised")
	}
	if b, ok := g.mapping[a]; ok {
		return b, nil
	}
	return -1, errors.New(fmt.Sprint(a, " not found in mapping"))
}

func (g *Graph) AddMappedVertex(ID string) int {
	if i, ok := g.mapping[ID]; ok {
		return i
	}
	i := g.highestMapIndex
	g.highestMapIndex++
	g.mapping[ID] = i
	g.Verticies[i] = Vertex{ID: i}
	return i
}

func (g *Graph) AddMappedArc(Source, Destination string, Distance int64) error {
	return g.AddArc(g.AddMappedVertex(Source), g.AddMappedVertex(Destination), Distance)
}

func (g *Graph) AddArc(Source, Destination int, Distance int64) error {
	if len(g.Verticies) <= Source || len(g.Verticies) <= Destination {
		return errors.New("Source/Destination not found")
	}
	g.Verticies[Source].AddArc(Destination, Distance)
	return nil
}
