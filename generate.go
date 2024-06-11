package dijkstra

import "math/rand"
import "math"

// Generate generates file with the amount of nodes specified
func Generate(nodes int) Graph {
	//reproducable random
	verticies := nodes
	seeded := rand.New(rand.NewSource(int64(verticies)))
	graph := Graph{}
	var from int
	for from = 0; from < verticies; from++ {
		graph.AddVertex(from)
	}
	for from = 0; from < verticies; from++ {
		// v := map[int]int64{}
		min := from - verticies/4
		if min < 2 {
			min = 2
		}
		max := from + verticies/4
		if max > verticies {
			max = verticies
		}
		for to := min; to < max; to++ {
			if to == from {
				continue
			}
			graph.AddArc(from, to, int64(2*verticies-to)+seeded.Int63n(int64(verticies)*int64(verticies-to+1)))
		}
	}
	return graph
}

// GenerateWorstCase generates a graph with the worst case scenario for Dijkstra's algorithm, has to hit every node on the way
func GenerateWorstCase(nodes int) Graph {
	graph := Graph{}
	verticies := nodes
	var err error
	var from int
	for from = 0; from < verticies; from++ {
		graph.AddVertex(from)
	}
	for from = 0; from < verticies; from++ {
		for to := from + 1; to < verticies-1; to++ {
			distance := int64(to - from)
			err = graph.AddArc(from, to, int64((verticies-from))*distance*distance)
			if err != nil {
				panic(err)
			}
		}
	}
	err = graph.AddArc(0, verticies-1, math.MaxInt64/2)
	if err != nil {
		panic(err)
	}
	err = graph.AddArc(verticies-1, 0, 0)
	if err != nil {
		panic(err)
	}
	return graph
}
