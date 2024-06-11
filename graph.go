package dijkstra

// Graph contains all the graph details (verticies and arcs)
type Graph struct {
	//slice of all verticies available
	// maps need to be initialised so if vertexArcs[i] == nil then vertex i
	// does not exist
	vertexArcs []map[int]uint64
}

// NewGraph creates a new empty graph
func NewGraph() Graph {
	new := Graph{}
	return new
}

// AddNewVertex adds a new vertex at the next available index
func (g *Graph) AddNewEmptyVertex() (index int) {
	for i, v := range g.vertexArcs {
		if v == nil {
			g.vertexArcs[i] = map[int]uint64{}
			return i
		}
	}
	index = len(g.vertexArcs)
	//only error can be id colision
	_ = g.AddEmptyVertex(index)
	return index
}

// AddVertex adds a single vertex at the specified index
func (g *Graph) AddEmptyVertex(index int) error {
	if err := g.vertexAvailable(index); err != nil {
		return err
	}
	g.addVertex(index, map[int]uint64{})
	return nil
}

func (g Graph) vertexOK(index int) error {
	if index < 0 {
		return newErrVertexNegative(index)
	}
	return nil
}

func (g Graph) vertexExists(index int) error {
	if index >= len(g.vertexArcs) ||
		g.vertexArcs[index] == nil {
		return newErrVertexNotFound(index)
	}
	return nil
}
func (g Graph) vertexAvailable(index int) error {
	if err := g.vertexOK(index); err != nil {
		return err
	}
	if index < len(g.vertexArcs) && g.vertexArcs[index] != nil {
		return newErrVertexAlreadyExists(index)
	}
	return nil
}

func (g Graph) vertexValid(index int) error {
	if err := g.vertexOK(index); err != nil {
		return err
	}
	if err := g.vertexExists(index); err != nil {
		return err
	}
	return nil
}

// GetArc gets the arc distance from src to dest
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

// AddVertex adds or overwrites a vertex with specified arcs, if the destination of an arc
// does not exist, it will error. If arc destinations should be added, use
// AddVertexAndArcs instead
func (g *Graph) AddVertex(index int, vertexArcs map[int]uint64) error {
	if err := g.vertexOK(index); err != nil {
		return err
	}
	for to := range vertexArcs {
		if err := g.vertexOK(to); err != nil {
			return err
		}
		if err := g.vertexExists(to); err != nil {
			return newErrArcNotFound(index, to)
		}
	}
	g.addVertex(index, vertexArcs)
	return nil
}

// AddVertexAndArcs adds or overwrites a vertex with specified arcs, if the destination of an
// arc does not exist, it will be added.
func (g *Graph) AddVertexAndArcs(index int, vertexArcs map[int]uint64) error {
	var err error
	if err = g.vertexOK(index); err != nil {
		return err
	}
	for to := range vertexArcs {
		if err = g.vertexOK(to); err != nil {
			return err
		}
		if err = g.vertexExists(to); err != nil {
			if err = g.AddEmptyVertex(to); err != nil {
				return err
			}
		}
	}
	g.addVertex(index, vertexArcs)
	return nil
}
func (g *Graph) addVertex(index int, vertexArcs map[int]uint64) {
	if index >= len(g.vertexArcs) {
		if index == len(g.vertexArcs) {
			g.vertexArcs = append(g.vertexArcs, vertexArcs)
		} else {
			g.vertexArcs = append(g.vertexArcs, make([]map[int]uint64, index-len(g.vertexArcs)+1)...)
		}
	}
	g.vertexArcs[index] = vertexArcs
}

// AddArc adds an arc from src to dest with the specified distance
// will error if the arc already exists (remove then add)
func (g *Graph) AddArc(src, dest int, distance uint64) error {
	if src >= len(g.vertexArcs) || g.vertexArcs[src] == nil {
		return newErrVertexNotFound(src)
	}
	if src < 0 {
		return newErrVertexNegative(src)
	}
	if dest >= len(g.vertexArcs) || g.vertexArcs[dest] == nil {
		return newErrVertexNotFound(dest)
	}
	if dest < 0 {
		return newErrVertexNegative(dest)
	}
	g.vertexArcs[src][dest] = distance
	return nil
}

// RemoveArc removes the arc from src to dest
func (g *Graph) RemoveArc(src, dest int) error {
	if src >= len(g.vertexArcs) || g.vertexArcs[src] == nil {
		return newErrVertexNotFound(src)
	}
	if src < 0 {
		return newErrVertexNegative(src)
	}
	if dest >= len(g.vertexArcs) || g.vertexArcs[dest] == nil {
		return newErrVertexNotFound(dest)
	}
	if dest < 0 {
		return newErrVertexNegative(dest)
	}
	if _, ok := g.vertexArcs[src][dest]; !ok {
		return newErrArcNotFound(src, dest)
	}
	delete(g.vertexArcs[src], dest)
	return nil
}

// GetVertexArcs gets all the arcs from vertex index
func (g Graph) GetVertexArcs(index int) (map[int]uint64, error) {
	if err := g.vertexValid(index); err != nil {
		return nil, err
	}
	return g.vertexArcs[index], nil
}

// RemoveVertexAndArcs removes the vertex at index and all arcs pointing to it
func (g *Graph) RemoveVertexAndArcs(index int) error {
	if err := g.vertexValid(index); err != nil {
		return err
	}
	for i := range g.vertexArcs {
		delete(g.vertexArcs[i], index)
	}
	g.vertexArcs[index] = nil
	return nil
}

// RemoveVertex removes the vertex at index, fails if there are still
// arcs pointing to the vertex. To remove all arcs also, use RemoveVertexAndArcs
func (g *Graph) RemoveVertex(index int) error {
	if err := g.vertexValid(index); err != nil {
		return err
	}
	for i, v := range g.vertexArcs {
		for to := range v {
			if to == index {
				return newArcHanging(i, index)
			}
		}
	}
	g.vertexArcs[index] = nil
	return nil
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
