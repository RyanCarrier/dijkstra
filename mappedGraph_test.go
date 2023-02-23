package dijkstra

import (
	"fmt"
	"testing"
)

func TestMapping(t *testing.T) {
	g := Graph{}
	if _, err := g.GetMapped(0); err == nil || err.Error() != ErrNoMap.Error() {
		t.Error("No init map should return correct err\n", err)
	}
	if _, err := g.GetMapping("A"); err == nil || err.Error() != ErrNoMap.Error() {
		t.Error("No init map should return correct err\n", err)
	}
	g = *NewGraph()
	g.AddMappedVertex("A")
	if _, err := g.GetMapped(1); err == nil || err.Error() != fmt.Sprint(1, " not found in mapping") {
		t.Error("Empty map should return correct err\n", err)
	}
	if _, err := g.GetMapping("B"); err == nil || err.Error() != fmt.Sprint("B", " not found in mapping") {
		t.Error("Empty map should return correct err\n", err)
	}
	//One will create, the other will get the created
	a := g.AddMappedVertex("B")
	b := g.AddMappedVertex("B")
	if a != b {
		t.Error("Adding same map should return same index")
	}
	if got, err := g.GetMapped(1); got != "B" || err != nil {
		t.Error("GetMapped failed")
	}
	if got, err := g.GetMapping("B"); got != 1 || err != nil {
		t.Error("GetMapping failed")
	}
	if err := g.AddMappedArc("A", "C", 1); err != nil {
		t.Error("AddMappedArc was not successful when destination not created yet")
	}
	if err := g.AddMappedArc("C", "A", 1); err != nil {
		t.Error("AddMappedArc was not successful when source not created yet")
	}
	if err := g.AddMappedArc("A", "B", 1); err != nil {
		t.Error("AddMappedArc was failed when valid")
	}
	if err := g.AddArc(0, 99, 1); err == nil {
		t.Error("AddArc didn't fail when referencing non existant vertex (destination)")
	}
	if err := g.AddArc(99, 0, 1); err == nil {
		t.Error("AddArc didn't fail when referencing non existant vertex (source)")
	}
}

func TestRemoveArc(t *testing.T) {
	var err error
	g := newGraph()
	if err = g.RemoveArc(1, 10); err == nil {
		t.Error("RemoveArc should fail on verticies that don't exist")
	}
	g.AddVertex(1)
	if err = g.RemoveArc(1, 10); err == nil {
		t.Error("RemoveArc should fail on destination verticies that don't exist")
	}
	g.AddVertex(10)
	g.AddVertex(100)
	if err = g.RemoveArc(1, 10); err != nil {
		t.Error("RemoveArc should not fail on Verticies that exist")
	}
	v, _ := g.GetVertex(1)
	if _, ok := v.GetArc(10); ok {
		t.Error("Arc should not yet exist")
	}
	_ = g.AddArc(1, 10, 100)
	if _, ok := v.GetArc(10); !ok {
		t.Error("Arc should exist")
	}
	if err = g.RemoveArc(1, 10); err != nil {
		t.Error("RemoveArc should not fail when removing valid arcs")
	}
	if _, ok := v.GetArc(10); ok {
		t.Error("Arc not should exist")
	}
}
