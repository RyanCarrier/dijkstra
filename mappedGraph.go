package dijkstra

import "errors"

type MappedGraph[T comparable] struct {
	graph   Graph
	mapping map[T]int
}

// NewMappedGraph creates a new empty mapped graph
func NewMappedGraph[T comparable]() MappedGraph[T] {
	return MappedGraph[T]{
		graph:   NewGraph(),
		mapping: make(map[T]int),
	}
}

// AddEmptyVertex adds a single empty vertex
func (mg *MappedGraph[T]) AddEmptyVertex(item T) error {
	_, err := mg.addMap(item)
	return err
}

// RemoveArc removes the arc between two vertices
func (mg *MappedGraph[T]) RemoveArc(src, dest T) error {
	srcIndex, destIndex, err := mg.getMap2(src, dest)
	if err != nil {
		return err
	}
	return mg.graph.RemoveArc(srcIndex, destIndex)
}

// AddArc adds an arc between two vertices
func (mg *MappedGraph[T]) AddArc(src, dest T, distance uint64) error {
	srcIndex, destIndex, err := mg.getMap2(src, dest)
	if err != nil {
		return err
	}
	return mg.graph.AddArc(srcIndex, destIndex, distance)
}

// GetArc gets the distance from src to dest
func (mg MappedGraph[T]) GetArc(src, dest T) (uint64, error) {
	var err error
	var dist uint64
	srcIndex, destIndex, err := mg.getMap2(src, dest)
	if err != nil {
		return 0, err
	}
	dist, err = mg.graph.GetArc(srcIndex, destIndex)
	if errors.Is(err, ErrArcNotFound) {
		return dist, newErrArcNotFound(src, dest)
	}
	return dist, err
}

// AddVertex adds a vertex with specified arcs, if the destination of an arc
// does not exist, it will error. If arc destinations should be added, use
// AddVertexAndArcs instead
func (mg *MappedGraph[T]) AddVertex(item T, vertexArcs map[T]uint64) error {
	index, err := mg.getMap(item)
	if err != nil {
		return err
	}
	mappedVertexArcs, err := mg.getMappedMap(vertexArcs)
	if err != nil {
		return err
	}
	err = mg.graph.AddVertex(index, mappedVertexArcs)
	if errors.Is(err, ErrArcNotFound) {
		//this should NEVER happen
		return err
	}
	return err
}

// AddVertexAndArcs adds a vertex with specified arcs, if the destination of an
// arc does not exist, it will be added.
func (mg *MappedGraph[T]) AddVertexAndArcs(item T, vertexArcs map[T]uint64) error {
	index, err := mg.getMap(item)
	if err != nil {
		return err
	}
	mappedVertexArcs, err := mg.getMappedMapCreate(vertexArcs)
	if err != nil {
		return err
	}
	return mg.graph.AddVertexAndArcs(index, mappedVertexArcs)
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
func (mg MappedGraph[T]) getInverseMap(index int) (T, error) {
	for k, v := range mg.mapping {
		if v == index {
			return k, nil
		}
	}
	var t T
	return t, newErrMapNotFound(index)
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
func (mg *MappedGraph[T]) addVertex(item T, vertexArcs map[T]uint64) error {
	index, ok := mg.mapping[item]
	if ok {
		return newErrMappedVertexAlreadyExists(item)
	}
	realVertex := make(map[int]uint64)
	for toID, v := range vertexArcs {
		toid, err := mg.getMap(toID)
		if err != nil {
			return err
		}
		realVertex[toid] = v
	}
	return mg.graph.AddVertex(index, realVertex)
}

// GetVertex gets the reference of the specified vertex. An error is thrown if
// there is no vertex with that index/ID.
func (mg *MappedGraph[T]) GetVertexArcs(item T) (map[T]uint64, error) {
	var arcs map[int]uint64
	id, err := mg.getMap(item)
	if err != nil {
		return nil, err
	}
	arcs, err = mg.graph.GetVertexArcs(id)
	if err != nil {
		return nil, err
	}
	return mg.getInverseMappedMap(arcs)
}

func (mg *MappedGraph[T]) RemoveVertexAndArcs(item T) error {
	index, err := mg.getMap(item)
	if err != nil {
		return err
	}
	return mg.graph.RemoveVertexAndArcs(index)
}

func (mg *MappedGraph[T]) RemoveVertex(item T) error {
	index, err := mg.getMap(item)
	if err != nil {
		return err
	}
	return mg.graph.RemoveVertex(index)
}

func (mg MappedGraph[T]) getInverseMappedMap(arcs map[int]uint64) (map[T]uint64, error) {
	mappedArcs := make(map[T]uint64)
	for k, v := range arcs {
		name, err := mg.getInverseMap(k)
		if err != nil {
			return mappedArcs, err
		}
		mappedArcs[name] = v
	}
	return mappedArcs, nil
}

func (mg MappedGraph[T]) getMappedMap(arcs map[T]uint64) (map[int]uint64, error) {
	mappedArcs := make(map[int]uint64)
	for k, v := range arcs {
		id, err := mg.getMap(k)
		if err != nil {
			return mappedArcs, err
		}
		mappedArcs[id] = v
	}
	return mappedArcs, nil
}

func (mg MappedGraph[T]) getMappedMapCreate(arcs map[T]uint64) (map[int]uint64, error) {
	mappedArcs := make(map[int]uint64)
	for k, v := range arcs {
		if _, ok := mg.mapping[k]; !ok {
			_, err := mg.addMap(k)
			if err != nil {
				return mappedArcs, err
			}
		}
		id, err := mg.getMap(k)
		if err != nil {
			return mappedArcs, err
		}
		mappedArcs[id] = v
	}
	return mappedArcs, nil
}

func (mg MappedGraph[T]) validate() error {
	for k, v := range mg.mapping {
		if _, err := mg.graph.GetVertexArcs(v); err != nil {
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

func (mg MappedGraph[T]) evaluate(src, dest T, shortest bool) (BestPath[T], error) {
	var bp BestPath[T]
	var bpOriginal BestPath[int]
	srcId, destId, err := mg.getMap2(src, dest)
	if err != nil {
		return bp, err
	}
	if shortest {
		bpOriginal, err = mg.graph.Shortest(srcId, destId)
	} else {
		bpOriginal, err = mg.graph.Longest(srcId, destId)
	}
	if err != nil {
		return bp, err
	}
	return mg.toMappedBestPath(bpOriginal)
}

func (mg MappedGraph[T]) Shortest(src, dest T) (BestPath[T], error) {
	return mg.evaluate(src, dest, true)
}

func (mg MappedGraph[T]) Longest(src, dest T) (BestPath[T], error) {
	return mg.evaluate(src, dest, false)
}

func (mg MappedGraph[T]) evaluateAll(src, dest T, shortest bool) (BestPaths[T], error) {
	var bp BestPaths[T]
	var bpOriginal BestPaths[int]
	srcId, destId, err := mg.getMap2(src, dest)
	if err != nil {
		return bp, err
	}
	if shortest {
		bpOriginal, err = mg.graph.ShortestAll(srcId, destId)
	} else {
		bpOriginal, err = mg.graph.LongestAll(srcId, destId)
	}
	if err != nil {
		return bp, err
	}
	return mg.toMappedBestPaths(bpOriginal)
}

func (mg MappedGraph[T]) ShortestAll(src, dest T) (BestPaths[T], error) {
	return mg.evaluateAll(src, dest, true)
}

func (mg MappedGraph[T]) LongestAll(src, dest T) (BestPaths[T], error) {
	return mg.evaluateAll(src, dest, false)
}
