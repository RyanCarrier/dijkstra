package dijkstra

type Vertex map[int]uint64

// Graph contains all the graph details
type Graph struct {
	//slice of all verticies available
	vertexArcs []Vertex
}

// NewGraph creates a new empty graph
func NewGraph() Graph {
	new := Graph{}
	return new
}

// AddNewVertex adds a new vertex at the next available index
func (g *Graph) AddNewEmptyVertex() (id int) {
	for i, v := range g.vertexArcs {
		if v == nil {
			g.vertexArcs[i] = map[int]uint64{}
			return i
		}
	}
	id = len(g.vertexArcs)
	//only error can be id colision
	_ = g.AddEmptyVertex(id)
	return id
}

// AddVertex adds a single vertex
func (g *Graph) AddEmptyVertex(ID int) error {
	return g.AddVertex(ID, map[int]uint64{})
}

func (g Graph) vertexOK(ID int) error {
	if ID < 0 {
		return newErrVertexNegative(ID)
	}
	return nil
}

func (g Graph) vertexExists(ID int) error {
	if ID >= len(g.vertexArcs) ||
		g.vertexArcs[ID] == nil {
		return newErrVertexNotFound(ID)
	}
	return nil
}

func (g Graph) vertexValid(ID int) error {
	if err := g.vertexOK(ID); err != nil {
		return err
	}
	if err := g.vertexExists(ID); err != nil {
		return err
	}
	return nil
}

func (g Graph) GetArcs(src int) (map[int]uint64, error) {
	if err := g.vertexValid(src); err != nil {
		return nil, err
	}
	return g.vertexArcs[src], nil
}
func (g Graph) GetArc(src, dest int) (uint64, error) {
	err := g.vertexValid(src)
	if err != nil {
		return 0, err
	}
	err = g.vertexValid(dest)
	if err != nil {
		return 0, err
	}
	got, ok := g.vertexArcs[src][dest]
	if !ok {
		return 0, newErrArcNotFound(src, dest)
	}
	return got, nil
}

// AddVertex adds a single vertex
func (g *Graph) AddVertex(ID int, vertex Vertex) error {
	if err := g.vertexOK(ID); err != nil {
		return err
	}
	if ID < len(g.vertexArcs) && g.vertexArcs[ID] != nil {
		return newErrVertexAlreadyExists(ID)
	}
	for to := range vertex {
		if err := g.vertexOK(to); err != nil {
			return err
		}
		if err := g.vertexExists(to); err != nil {
			return newErrArcNotFound(ID, to)
		}
	}
	if ID >= len(g.vertexArcs) {
		if ID == len(g.vertexArcs) {
			g.vertexArcs = append(g.vertexArcs, vertex)
		} else {
			g.vertexArcs = append(g.vertexArcs, make([]Vertex, ID-len(g.vertexArcs)+1)...)
		}
	}
	g.vertexArcs[ID] = vertex
	return nil
}
func (g *Graph) AddArc(from, to int, distance uint64) error {
	if from >= len(g.vertexArcs) || g.vertexArcs[from] == nil {
		return newErrVertexNotFound(from)
	}
	if from < 0 {
		return newErrVertexNegative(from)
	}
	if to >= len(g.vertexArcs) || g.vertexArcs[to] == nil {
		return newErrVertexNotFound(to)
	}
	if to < 0 {
		return newErrVertexNegative(to)
	}
	g.vertexArcs[from][to] = distance
	return nil
}
func (g *Graph) RemoveArc(from, to int) error {
	if from >= len(g.vertexArcs) || g.vertexArcs[from] == nil {
		return newErrVertexNotFound(from)
	}
	if from < 0 {
		return newErrVertexNegative(from)
	}
	if to >= len(g.vertexArcs) || g.vertexArcs[to] == nil {
		return newErrVertexNotFound(to)
	}
	if to < 0 {
		return newErrVertexNegative(to)
	}
	if _, ok := g.vertexArcs[from][to]; !ok {
		return newErrArcNotFound(from, to)
	}
	delete(g.vertexArcs[from], to)
	return nil
}

// GetVertex gets the reference of the specified vertex. An error is thrown if
// there is no vertex with that index/ID.
func (g *Graph) GetVertex(ID int) (Vertex, error) {
	if ID >= len(g.vertexArcs) {
		return nil, newErrVertexNotFound(ID)
	}
	return g.vertexArcs[ID], nil
}

func (g Graph) validate() error {
	for vId, vertex := range g.vertexArcs {
		for k := range vertex {
			if k >= len(g.vertexArcs) || k < 0 {
				return newErrGraphNotValid(vId, k)
			}
		}
	}
	return nil
}
