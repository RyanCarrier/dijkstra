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

//NewVertex creates a new vertex
func NewVertex(ID int) *Vertex {
	return &Vertex{ID: ID, arcs: map[int]int64{}}
}

//AddVerticies adds the listed verticies to the graph, overwrites any existing
// Vertex with the same ID.
func (g *Graph) AddVerticies(verticies ...Vertex) {
	for _, v := range verticies {
		if v.ID >= len(g.Verticies) {
			newV := make([]Vertex, v.ID+1-len(g.Verticies))
			g.Verticies = append(g.Verticies, newV...)
		}
		g.Verticies[v.ID] = v
	}
}

//AddArc adds an arc to the vertex, it's up to the user to make sure this is used
// correctly, firstly ensuring to use before adding to graph, or to use referenced
// of the Vertex instead of a copy. Secondly, to ensure the destination is a valid
// Vertex in the graph. Note that AddArc will overwrite any existing distance set
// if there is already an arc set to Destination.
func (v *Vertex) AddArc(Destination int, Distance int64) {
	if v.arcs == nil {
		v.arcs = map[int]int64{}
	}
	v.arcs[Destination] = Distance
}

/*
I decided you don't get that kind of privelage
#checkyourprivelage
//RemoveArc completely removes the arc to Destination (if it exists)
func (v *Vertex) RemoveArc(Destination int) {
	delete(v.arcs, Destination)
}*/

//GetArc gets the specified arc to Destination, bool is false if no arc found
func (v *Vertex) GetArc(Destination int) (distance int64, ok bool) {
	if v.arcs == nil {
		return 0, false
	}
	//idk why but doesn't work on one line?
	distance, ok = v.arcs[Destination]
	return
}
