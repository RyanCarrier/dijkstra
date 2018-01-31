package dijkstra

import "math/rand"

//Generate generates file with the amount of nodes specified
func Generate(nodes int) Graph {
	//	fmt.Println("Generating file "+filename+" with nodes ", nodes)
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
