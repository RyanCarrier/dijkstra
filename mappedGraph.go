package dijkstra

type MappedGraph[T comparable] struct {
	graph   Graph
	mapping map[T]int
}

// NewGraph creates a new empty graph
func NewMappedGraph[T comparable]() MappedGraph[T] {
	return MappedGraph[T]{
		graph:   NewGraph(),
		mapping: make(map[T]int),
	}
}

// AddVertex adds a single vertex
func (mg *MappedGraph[T]) AddEmptyVertex(item T) error {
	_, err := mg.addMap(item)
	return err
}

func (mg *MappedGraph[T]) RemoveArc(from, to T) error {
	src, dest, err := mg.getMap2(from, to)
	if err != nil {
		return err
	}
	return mg.graph.RemoveArc(src, dest)
}
func (mg *MappedGraph[T]) AddArc(from, to T, distance uint64) error {
	src, dest, err := mg.getMap2(from, to)
	if err != nil {
		return err
	}
	return mg.graph.AddArc(src, dest, distance)
}

func (mg MappedGraph[T]) GetArc(src, dest int) (uint64, error) {
	if src >= len(mg.graph.vertexArcs) {
		return 0, newErrVertexNotFound(src)
	}
	got, ok := mg.graph.vertexArcs[src][dest]
	if !ok {
		return 0, newErrArcNotFound(src, dest)
	}
	return got, nil
}
func (mg MappedGraph[T]) getMap2(a, b T) (int, int, error) {
	var aid, bid int
	var err error
	aid, err = mg.getMap(a)
	if err != nil {
		return aid, bid, err
	}
	bid, err = mg.getMap(b)
	return aid, bid, err
}
func (mg MappedGraph[T]) getMap(item T) (int, error) {
	id, ok := mg.mapping[item]
	if !ok {
		return id, newErrMappedVertexNotFound(item)
	}
	return id, nil
}
func (mg MappedGraph[T]) getInverseMap(id int) (T, error) {
	for k, v := range mg.mapping {
		if v == id {
			return k, nil
		}
	}
	var t T
	return t, newErrMapNotFound(id)
}
func (mg *MappedGraph[T]) addMap(item T) (int, error) {
	var id int
	var ok bool
	id, ok = mg.mapping[item]
	if ok {
		return id, newErrMappedVertexAlreadyExists(item)
	}
	id = mg.graph.AddNewEmptyVertex()
	mg.mapping[item] = id
	return id, nil
}

// AddVertex adds a single vertex
func (mg *MappedGraph[T]) AddVertex(item T, vertex map[T]uint64) error {
	id, ok := mg.mapping[item]
	if ok {
		return newErrMappedVertexAlreadyExists(item)
	}
	realVertex := make(map[int]uint64)
	for toID, v := range vertex {
		toid, err := mg.getMap(toID)
		if err != nil {
			return err
		}
		realVertex[toid] = v
	}
	return mg.graph.AddVertex(id, realVertex)
}

// GetVertex gets the reference of the specified vertex. An error is thrown if
// there is no vertex with that index/ID.
func (mg *MappedGraph[T]) GetVertex(item T) (Vertex, error) {
	id, err := mg.getMap(item)
	if err != nil {
		return nil, err
	}
	return mg.graph.GetVertex(id)
}

func (mg MappedGraph[T]) validate() error {
	for k, v := range mg.mapping {
		if _, err := mg.graph.GetVertex(v); err != nil {
			return newErrMappedVertexNotFound(k)
		}
	}
	return mg.graph.validate()
}
func (mg MappedGraph[T]) inverseMapPath(path []int) ([]T, error) {
	var err error
	var result []T = make([]T, len(path))
	for i, v := range path {
		result[i], err = mg.getInverseMap(v)
		if err != nil {
			return result, err
		}
	}
	return result, err
}
func (mg MappedGraph[T]) toMappedBestPath(bp BestPath[int]) (BestPath[T], error) {
	var err error
	var result BestPath[T]
	result.Distance = bp.Distance
	result.Path, err = mg.inverseMapPath(bp.Path)
	return result, err
}
func (mg MappedGraph[T]) toMappedBestPaths(bp BestPaths[int]) (BestPaths[T], error) {
	var err error
	var result BestPaths[T]
	result.Paths = make([][]T, len(bp.Paths))
	for i, v := range bp.Paths {
		result.Paths[i], err = mg.inverseMapPath(v)
		if err != nil {
			return result, err
		}
	}
	return result, err
}
func (mg MappedGraph[T]) Shortest(src, dest T) (BestPath[T], error) {
	var bp BestPath[T]
	var bpOriginal BestPath[int]
	srcId, destId, err := mg.getMap2(src, dest)
	if err != nil {
		return bp, err
	}
	bpOriginal, err = mg.graph.Shortest(srcId, destId)
	if err != nil {
		return bp, err
	}
	return mg.toMappedBestPath(bpOriginal)
}
func (mg MappedGraph[T]) ShortestAll(src, dest T) (BestPaths[T], error) {
	var bp BestPaths[T]
	var bpOriginal BestPaths[int]
	srcId, destId, err := mg.getMap2(src, dest)
	if err != nil {
		return bp, err
	}
	bpOriginal, err = mg.graph.ShortestAll(srcId, destId)
	if err != nil {
		return bp, err
	}
	return mg.toMappedBestPaths(bpOriginal)
}
func (mg MappedGraph[T]) Longest(src, dest T) (BestPath[T], error) {
	var bp BestPath[T]
	var bpOriginal BestPath[int]
	srcId, destId, err := mg.getMap2(src, dest)
	if err != nil {
		return bp, err
	}
	bpOriginal, err = mg.graph.Longest(srcId, destId)
	if err != nil {
		return bp, err
	}
	return mg.toMappedBestPath(bpOriginal)
}
func (mg MappedGraph[T]) LongestAll(src, dest T) (BestPaths[T], error) {
	var bp BestPaths[T]
	var bpOriginal BestPaths[int]
	srcId, destId, err := mg.getMap2(src, dest)
	if err != nil {
		return bp, err
	}
	bpOriginal, err = mg.graph.LongestAll(srcId, destId)
	if err != nil {
		return bp, err
	}
	return mg.toMappedBestPaths(bpOriginal)
}
