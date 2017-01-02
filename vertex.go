package dijkstra

//Vertex is a single node in the network, contains it's ID, best distance (to
// itself from the src) and the weight to go to each other connected node (Vertex)
type Vertex struct {
	//ID of the Vertex
	ID int
	//Best Distance to the Vertex
	Distance   int64
	BestVertex int
	//A set of all weights to the nodes in the map
	Arcs map[int]int64
}
