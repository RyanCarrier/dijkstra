package dijkstra

//BestPath contains the solution of the most optimal path
type BestPath struct {
	Distance int64
	Path     []int
}

func (g *Graph) bestPath(src, dest int) BestPath {
	var path []int
	for c := g.Verticies[dest]; c.ID != src; c = g.Verticies[c.bestVertex] {
		path = append(path, c.ID)
	}
	path = append(path, src)
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return BestPath{g.Verticies[dest].distance, path}
}
