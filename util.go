package dijkstra

import (
	"io/ioutil"
	"strconv"
	"strings"
)

//Import imports a graph from the specified file returns the Graph, a map for
// if the nodes are not integers and an error if needed.
func Import(filename string) (g Graph, m map[string]int, err error) {
	mapping := false
	var lowestIndex int
	var i int
	var arc int
	var dist int64
	got, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	m = map[string]int{}
	g.Visiting = NewList()

	input := strings.TrimSpace(string(got))
	for _, line := range strings.Split(input, "\n") {
		f := strings.Fields(line)
		//no need to check for size cause there must be something as the string is trimmed and split
		if mapping {
			if _, ok := m[f[0]]; !ok {
				m[f[0]] = lowestIndex
				i = lowestIndex
				lowestIndex++
			}
		} else {
			i, err = strconv.Atoi(f[0])
			if err != nil {
				mapping = true
				m[f[0]] = lowestIndex
				i = lowestIndex
				lowestIndex++
			}
		}
		if temp := len(g.Verticies); temp <= i { //Extend if we have to
			g.Verticies = append(g.Verticies, make([]Vertex, 1+i-len(g.Verticies))...)
			for ; temp < len(g.Verticies); temp++ {
				g.Verticies[temp].ID = temp
				g.Verticies[temp].Arcs = map[int]int64{}
			}
		}
		//g.Verticies[i].ID = i
		//g.Verticies[i].Arcs = map[int]int64{}
		if len(f) == 1 {
			//if there is no FROM here
			continue
		}
		for _, set := range f[1:] {
			got := strings.Split(set, ",")
			if len(got) != 2 {
				err = ErrWrongFormat
				return
			}
			dist, err = strconv.ParseInt(got[1], 10, 64)
			if err != nil {
				err = ErrWrongFormat
				return
			}
			if mapping {
				var ok bool
				arc, ok = m[got[0]]
				if !ok {
					arc = lowestIndex
					m[got[0]] = arc
					lowestIndex++
				}
			} else {
				arc, err = strconv.Atoi(got[0])
				if err != nil {
					err = ErrMixMapping
					return
				}
			}
			g.Verticies[i].Arcs[arc] = dist
		}
	}
	g.Visited = make([]bool, len(g.Verticies))
	err = g.validate()
	return
}
