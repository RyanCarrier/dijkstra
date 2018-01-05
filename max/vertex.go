package max

//Set is a set of in and out tax's
type Set struct {
	In  int64
	Out int64
}

//Vertex is a single node in the network, contains it's ID, best distance (to
// itself from the src) and the weight to go to each other connected node (Vertex)
type Vertex struct {
	//ID of the Vertex
	ID int
	//Best distance to the Vertex
	best       int64
	bestVertex int
	Multiplier Set
	Addition   Set
	//A set of all weights to the nodes in the map
	arcs map[int]int64
}

//NewVertex creates a new vertex
func NewVertex(ID int, Multiplier, Addition Set) *Vertex {
	return &Vertex{ID: ID, arcs: map[int]int64{}, Multiplier: Multiplier, Addition: Addition}
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

//Evaluate evalutates the value in to if converted from this node (v) to the
// next node (to)
func (v *Vertex) Evaluate(to Vertex, conversion int64) int64 {
	return ((((((v.best * (multiplier - v.Multiplier.Out) / multiplier) - v.Addition.Out) * conversion) / multiplier) * to.Multiplier.In) / multiplier) - to.Addition.In
}
