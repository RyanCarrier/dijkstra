package dijkstra

import (
	"fmt"
	"strconv"
	"strings"
)

// Import imports a graph from the specified file returns the Graph, a map for
// if the nodes are not integers and an error if needed.
func Import(data string) (g Graph, err error) {
	var i int
	var arc int
	var dist uint64

	input := strings.TrimSpace(string(data))
	for _, line := range strings.Split(input, "\n") {
		f := strings.Fields(strings.TrimSpace(line))
		if len(f) == 0 || (len(f) == 1 && f[0] == "") {
			continue
		}
		//no need to check for size cause there must be something as the string is trimmed and split
		i, err = strconv.Atoi(f[0])
		if err != nil {
			return g, err
		}
		if temp := len(g.vertexArcs); temp <= i { //Extend if we have to
			g.vertexArcs = append(g.vertexArcs, make([]Vertex, 1+i-len(g.vertexArcs))...)
			for ; temp < len(g.vertexArcs); temp++ {
				g.vertexArcs[temp] = map[int]uint64{}
			}
		}
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
			dist, err = strconv.ParseUint(got[1], 10, 64)
			if err != nil {
				err = ErrWrongFormat
				return
			}
			arc, err = strconv.Atoi(got[0])
			if err != nil {
				return g, err
			}
			g.vertexArcs[i][arc] = dist
		}
	}
	return g, g.validate()
}
func ImportStringMapped(data string) (mg MappedGraph[string], err error) {
	var lowestIndex int
	var i int
	var arc int
	var dist uint64
	var ok bool
	mg.mapping = map[string]int{}

	input := strings.TrimSpace(string(data))
	for _, line := range strings.Split(input, "\n") {
		f := strings.Fields(strings.TrimSpace(line))
		if len(f) == 0 || (len(f) == 1 && f[0] == "") {
			continue
		}
		//no need to check for size cause there must be something as the string is trimmed and split
		if i, ok = mg.mapping[f[0]]; !ok {
			mg.mapping[f[0]] = lowestIndex
			i = lowestIndex
			lowestIndex++
		}

		if temp := len(mg.graph.vertexArcs); temp <= i { //Extend if we have to
			mg.graph.vertexArcs = append(mg.graph.vertexArcs, make([]Vertex, 1+i-len(mg.graph.vertexArcs))...)
			for ; temp < len(mg.graph.vertexArcs); temp++ {
				mg.graph.vertexArcs[temp] = map[int]uint64{}
			}
		}
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
			dist, err = strconv.ParseUint(got[1], 10, 64)
			if err != nil {
				err = ErrWrongFormat
				return
			}

			arc, ok = mg.mapping[got[0]]
			if !ok {
				arc = lowestIndex
				mg.mapping[got[0]] = arc
				lowestIndex++
			}

			mg.graph.vertexArcs[i][arc] = dist
		}
	}
	err = mg.validate()
	return
}

func (g Graph) ToString() (string, error) {
	var result = strings.Builder{}
	for id, v := range g.vertexArcs {
		result.WriteString(strconv.Itoa(id))
		for key, val := range v {
			result.WriteString(" " + strconv.Itoa(key) + "," + strconv.FormatUint(val, 10))
		}
		result.WriteRune('\n')
	}
	return result.String(), nil
}
func (mg MappedGraph[T]) ToString() (string, error) {
	var err error
	var result = strings.Builder{}
	for k, v := range mg.mapping {
		result.WriteString(fmt.Sprint(k))
		if err = mg.graph.vertexValid(v); err != nil {
			return "", err
		}
		for to, dist := range mg.graph.vertexArcs[v] {
			name, err := mg.getInverseMap(to)
			if err != nil {
				return "", err
			}
			result.WriteString(" " + fmt.Sprint(name) + "," + strconv.FormatUint(dist, 10))
		}
		result.WriteRune('\n')
	}
	return result.String(), nil
}
