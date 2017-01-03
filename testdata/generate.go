package main

import (
	"flag"
	"log"
	"math/rand"
	"strconv"

	"github.com/RyanCarrier/dijkstra"
)

var nodes int
var filename string

func init() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatal("Pls specify amount of nodes")
	}
	var err error
	nodes, err = strconv.Atoi(flag.Args()[0])
	if err != nil {
		log.Fatal("Pls use numbers not words")
	}
	filename = flag.Args()[0] + ".txt"
}

func main() {
	graph := dijkstra.Graph{}
	var v dijkstra.Vertex
	var i int
	for i = 0; i < nodes; i++ {
		v = dijkstra.Vertex{ID: i}
		v.Arcs = map[int]int64{}
		for j := 0; j < nodes; j++ {
			if j == i {
				continue
			}
			v.Arcs[j] = max(int64(nodes-j)-rand.Int63n(int64(nodes)*int64(nodes-j)), 1)
		}
		graph.AddVerticies(v)
	}
	err := graph.ExportToFile(filename)
	if err != nil {
		log.Fatal(err)
	}
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
