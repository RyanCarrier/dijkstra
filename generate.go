package dijkstra

import "math/rand"

//Generate generates file with the amount of nodes specified
func Generate(nodes int) Graph {
	graph := Graph{}
	var i int
	for i = 0; i < nodes; i++ {
		v := NewVertex(i)
		for j := 0; j < nodes; j++ {
			if j == i {
				continue
			}
			v.AddArc(j, int64(2*nodes-j)+rand.Int63n(int64(nodes)*int64(nodes-j+1)))
		}
		graph.AddVerticies(*v)
	}
	return graph
}

//GenerateWorstCase generates a graph with the worst case scenario for Dijkstra's algorithm, has to hit every node on the way
func GenerateWorstCase(nodes int) Graph {
	graph := Graph{}
	var i int
	for i = 0; i < nodes; i++ {
		v := NewVertex(i)
		for j := 0; j < nodes; j++ {
			if j == i {
				continue
			}
			if j == i+1 {
				v.AddArc(j, int64(1))
				continue
			}
			v.AddArc(j, int64(2*nodes))
		}
		graph.AddVerticies(*v)
	}
	return graph
}
