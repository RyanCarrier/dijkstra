package dijkstra

//Vertex is a single node in the network, contains it's ID, best distance (to
// itself from the src) and the weight to go to each other connected node (Vertex)
type Vertex struct {
	//ID of the Vertex
	ID int
	//Best distance to the Vertex
	distance   int64
	bestVertex int
	//A set of all weights to the nodes in the map
	arcs map[int]int64
}

//AddVerticies adds the listed verticies to the graph
func (g *Graph) AddVerticies(v ...Vertex) {
	g.Verticies = append(g.Verticies, v...)
}

//AddArc adds an arc to the vertex, it's up to the user to make sure this is used
// correctly, firstly ensuring to use before adding to graph, or to use referenced
// of the Vertex instead of a copy. Secondly, to ensure the destination is a valid
// Vertex in the graph.
func (v *Vertex) AddArc(Destination int, Distance int64) {
	v.arcs[Destination] = Distance
}
