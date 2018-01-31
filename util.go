package dijkstra

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

//Import imports a graph from the specified file returns the Graph, a map for
// if the nodes are not integers and an error if needed.
func Import(filename string) (g Graph, err error) {
	g.usingMap = false
	var lowestIndex int
	var i int
	var arc int
	var dist int64
	var ok bool
	got, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	g.mapping = map[string]int{}

	input := strings.TrimSpace(string(got))
	for _, line := range strings.Split(input, "\n") {
		f := strings.Fields(strings.TrimSpace(line))
		if len(f) == 0 || (len(f) == 1 && f[0] == "") {
			continue
		}
		//no need to check for size cause there must be something as the string is trimmed and split
		if g.usingMap {
			if i, ok = g.mapping[f[0]]; !ok {
				g.mapping[f[0]] = lowestIndex
				i = lowestIndex
				lowestIndex++
			}
		} else {
			i, err = strconv.Atoi(f[0])
			if err != nil {
				g.usingMap = true
				g.mapping[f[0]] = lowestIndex
				i = lowestIndex
				lowestIndex++
			}
		}
		if temp := len(g.Verticies); temp <= i { //Extend if we have to
			g.Verticies = append(g.Verticies, make([]Vertex, 1+i-len(g.Verticies))...)
			for ; temp < len(g.Verticies); temp++ {
				g.Verticies[temp].ID = temp
				g.Verticies[temp].arcs = map[int]int64{}
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
			dist, err = strconv.ParseInt(got[1], 10, 64)
			if err != nil {
				err = ErrWrongFormat
				return
			}
			if g.usingMap {
				arc, ok = g.mapping[got[0]]
				if !ok {
					arc = lowestIndex
					g.mapping[got[0]] = arc
					lowestIndex++
				}
			} else {
				arc, err = strconv.Atoi(got[0])
				if err != nil {
					err = ErrMixMapping
					return
				}
			}
			g.Verticies[i].arcs[arc] = dist
		}
	}
	err = g.validate()
	return
}

//ExportToFile exports the verticies to file currently does not take into account
// mappings (from string to int)
func (g Graph) ExportToFile(filename string) error {
	var i string
	var err error
	os.MkdirAll(filename, 0777)
	if _, err = os.Stat(filename); err == nil {
		os.Remove(filename)
	}
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, v := range g.Verticies {
		if g.usingMap {
			if i, err = g.GetMapped(v.ID); err != nil {
				return errors.New("Mapping fail when exporting; " + err.Error())
			}
			fmt.Fprint(f, i)
		} else {
			fmt.Fprint(f, v.ID)
		}
		for key, val := range v.arcs {
			if g.usingMap {
				if i, err = g.GetMapped(key); err != nil {
					return errors.New("Mapping fail when exporting; " + err.Error())
				}
				fmt.Fprint(f, " ", i, ",", val)
			} else {
				fmt.Fprint(f, " ", key, ",", val)
			}
		}
		fmt.Fprint(f, "\n")
	}
	return nil
}
