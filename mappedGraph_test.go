package dijkstra

import (
	"fmt"
	"testing"
)

func TestMapping(t *testing.T) {
	g := Graph{}
	if _, err := g.GetMapped(0); err == nil || err.Error() != "Map is not being used/initialised" {
		t.Error("No init map should return correct err\n", err)
	}
	if _, err := g.GetMapping("A"); err == nil || err.Error() != "Map is not being used/initialised" {
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
	if g.AddMappedVertex("B") != g.AddMappedVertex("B") {
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
