package dijkstra

import (
	"math"
	"math/rand"
)

// Generate generates a graph with a moderate amount of connections with semi
// random weights
func Generate(verticies int) Graph {
	//reproducable random
	seeded := rand.New(rand.NewSource(int64(verticies)))
	graph := Graph{}
	var from int
	for from = 0; from < verticies; from++ {
		graph.AddEmptyVertex(from)
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
			graph.AddArc(from, to, uint64(2*verticies-to)+uint64(seeded.Int63n(int64(verticies)*int64(verticies-to+1))))
		}
	}
	return graph
}

func generateWorstCaseShortest(verticies int) (Graph, BestPath[int]) {
	graph := Graph{}
	var from int
	for from = 0; from < verticies; from++ {
		graph.AddEmptyVertex(from)
	}
	for from = 0; from < verticies; from++ {
		for to := from + 1; to < verticies-1; to++ {
			distance := uint64(to - from)
			graph.AddArc(from, to, uint64((verticies-from))*distance*distance)
		}
	}
	graph.AddArc(0, verticies-1, math.MaxUint64/2)
	graph.AddArc(verticies-1, 0, 0)
	shortestResult := BestPath[int]{math.MaxUint64 / 2, []int{0, verticies - 1}}
	return graph, shortestResult
}

func generateWorstCaseLongest(verticies int) (Graph, BestPath[int]) {
	graph := Graph{}
	var from int
	for from = 0; from < verticies; from++ {
		graph.AddEmptyVertex(from)
	}
	for from = 0; from < verticies; from++ {
		for to := from + 1; to < verticies-1; to++ {
			distance := uint64(verticies - (to - from))
			graph.AddArc(from, to, distance*distance)
		}
	}
	graph.AddArc(0, verticies-1, 0)
	longestResult := BestPath[int]{0, []int{0, verticies - 1}}
	return graph, longestResult
}
