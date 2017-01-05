package dijkstra

import "testing"

func TestMapInitialised(t *testing.T) {
	g := NewGraph()
	g.GetMapping("test")
	//if the program didn't crash then the map was initialised bahahah
}
