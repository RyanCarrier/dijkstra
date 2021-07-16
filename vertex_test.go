package dijkstra

import (
	"reflect"
	"testing"
)

func TestNewVertex(t *testing.T) {
	v := NewVertex(10)
	if v.ID != 10 {
		t.Error("NewVertex ID not set")
	}
	if v.arcs == nil {
		t.Error("NewVertex arcs map not initialised")
	}
}

func TestArc(t *testing.T) {
	v := Vertex{}
	dist, ok := v.GetArc(100)
	if ok || dist != 0 {
		t.Error("GetArc with no arcs failed (should 0,nil)")
	}
	v.AddArc(100, 1)
	if !reflect.DeepEqual(v.arcs, map[int]int64{100: 1}) {
		t.Error("AddArc failed to add arc")
	}
	v.AddArc(100, 2)
	if !reflect.DeepEqual(v.arcs, map[int]int64{100: 2}) {
		t.Error("AddArc failed to overwrite arc")
	}
	v.AddArc(101, 1)
	if !reflect.DeepEqual(v.arcs, map[int]int64{100: 2, 101: 1}) {
		t.Error("AddArc failed to add second arc")
	}
	dist, ok = v.GetArc(100)
	if !ok || dist != 2 {
		t.Error("GetArc failed")
	}
	v.RemoveArc(100)
	if !reflect.DeepEqual(v.arcs, map[int]int64{101: 1}) {
		t.Error("RemoveArc failed to remove arc")
	}
	_, ok = v.GetArc(100)
	if ok {
		t.Error("GetArc failed")
	}

}
