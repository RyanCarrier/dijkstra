package dijkstra

import "testing"

func TestGetVertex(t *testing.T) {
	g := NewGraph()
	g.AddEmptyVertex(99)
	if _, err := g.GetVertexArcs(99); err != nil {
		t.Error("Getting vertex failed (99)")
	}
	if _, err := g.GetVertexArcs(100); err == nil {
		t.Error("Vertex should not be found (100)")
	}
}
func TestAddVertex(t *testing.T) {
	g := NewGraph()
	g.AddEmptyVertex(99)
	_, err := g.GetVertexArcs(98)
	if err == nil {
		t.Error("should not have had ID set and err should not be nil")
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
	err = g.AddVertex(11, map[int]uint64{5: 1})
	if err == nil {
		t.Error("adding vertex should have failed")
	}
	err = g.AddVertexAndArcs(11, map[int]uint64{5: 1})
	if err == nil {
		t.Error("adding vertex should have failed")
	}
	err = g.AddVertex(100, map[int]uint64{101: 1})
	if err == nil {
		t.Error("arc is out of range, should have errored")
	}
	err = g.AddVertexAndArcs(100, map[int]uint64{101: 1})
	if err != nil {
		t.Error("adding vertex should have succeeded")
	}
	if len(g.vertexArcs) != 102 {
		//as we should have created index 101 with arc
		t.Error("vertexArcs should have been extended")
	}
	if g.vertexArcs[101] == nil {
		t.Error("arc should have been added")
	}
}

func TestValidateCorrect(t *testing.T) {
	for i, test := range testGraphsCorrect {
		if test.graph.validate() != nil {
			t.Error("Graph ", i, "not valid;\n", test.graph.validate(), " should be nil")
		}
	}
}
