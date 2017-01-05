package dijkstra

import "testing"

func TestMapInitialised(t *testing.T) {
	g := NewGraph()
	g.GetMapping("test")
	//if the program didn't crash then the map was initialised bahahah
}

func TestGetVertex(t *testing.T) {
	g := NewGraph()
	g.AddVertex(99)
	if v, _ := g.GetVertex(99); v.ID != 99 {
		t.Error("Getting vertex failed (99)")
	}
	if _, err := g.GetVertex(100); err == nil {
		t.Error("Vertex should not be found (100)")
	}
}

func TestValidateCorrect(t *testing.T) {
	if newGraph().validate() != nil {
		t.Error(newGraph().validate().Error(), " should be nil")
	}
}

func TestValidateIncorrect(t *testing.T) {
	if newBadGraph().validate() == nil {
		t.Error("graph should not have validated")
	}
}

func newGraph() Graph {
	return getBGraph()
}

func newBadGraph() Graph {
	g := getBGraph()
	g.Verticies[0].AddArc(9999, 1)
	return g
}
