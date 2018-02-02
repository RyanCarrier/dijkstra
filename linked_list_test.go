package dijkstra

import "testing"

//if the other shit is passing this is working lol

func TestLazyInit(t *testing.T) {
	ll := linkedList{root: NewVertex(-1), len: 0}
	if ll.root.next != nil || ll.root.prev != nil {
		t.Error("pre init should have root failing")
	}
	ll.init(true)
	if ll.root.next == nil || ll.root.prev == nil {
		t.Error("post init should have root passing")
	}
}

func TestEmptyList(t *testing.T) {
	ll := linkedList{root: NewVertex(-1), len: 0}
	if ll.front() != nil || ll.back() != nil {
		t.Error("empty list front()/back() should be nil")
	}
}

func TestAdding(t *testing.T) {
	verticies := getVertexs()
	short := false
	var popped *Vertex
	for i := 0; i < 2; i++ {
		ll := new(linkedList).init(short)
		checkNil(t, ll)
		ll.PushOrdered(verticies[0])
		checkNonNil(t, ll)
		ll.PushOrdered(verticies[1])
		checkNonNil(t, ll)
		if ll.PopOrdered() != verticies[i] {
			t.Error("Wrong pop")
		}
		checkNonNil(t, ll)
		ll.PushOrdered(verticies[2])
		checkNonNil(t, ll)
		if ll.PopOrdered() != verticies[i+1] {
			t.Error("Wrong pop")
		}
		popped = ll.PopOrdered()
		if (short && popped != verticies[0]) || (!short && popped != verticies[2]) {
			t.Error("Wrong pop")
		}
		checkNil(t, ll)
		short = !short
	}
}

func checkNil(t *testing.T, ll *linkedList) {
	if ll.front() != nil {
		t.Errorf("empty list front() should be nil \n[%+v]", ll.front())
	}
	if ll.back() != nil {
		t.Errorf("empty list back() should be nil \n[%+v]", ll.back())
	}
}

func checkNonNil(t *testing.T, ll *linkedList) {
	if ll.front() == nil {
		t.Error("non empty list should never (front()) give nil")
	}
	if ll.back() == nil {
		t.Error("non empty list should never (back()) give nil")
	}
}

func getVertexs() []*Vertex {
	v := make([]*Vertex, 0)
	for i := 0; i < 10; i++ {
		v = append(v, NewVertex(i))
		v[i].distance = int64(i + 1)
	}
	return v
}
