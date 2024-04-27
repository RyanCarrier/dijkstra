package dijkstra

import "testing"

func TestGetVertex(t *testing.T) {
	g := NewGraph()
	g.AddEmptyVertex(99)
	if _, err := g.GetVertex(99); err != nil {
		t.Error("Getting vertex failed (99)")
	}
	if _, err := g.GetVertex(100); err == nil {
		t.Error("Vertex should not be found (100)")
	}
}
func TestAddVertex(t *testing.T) {
	g := NewGraph()
	g.AddEmptyVertex(99)
	_, err := g.GetVertex(98)
	if err != nil {
		t.Error("should not have had ID set and err should be nil")
	}
	for i := 0; i <= 10; i++ {
		g.AddEmptyVertex(i)
	}
	v := g.AddNewEmptyVertex()
	if v != 11 {
		t.Error("Adding self assigned vertex fail")
	}
	g = NewGraph()
	for i := 0; i <= 10; i++ {
		g.AddNewEmptyVertex()
	}
	if v = g.AddNewEmptyVertex(); v != 11 {
		t.Error("Adding self assigned vertex fail when extending slice")
	}

}

func TestValidateCorrect(t *testing.T) {
	for i, test := range testGraphsCorrect {
		if test.graph.validate() != nil {
			t.Error("Graph ", i, "not valid;\n", test.graph.validate(), " should be nil")
		}
	}
}
