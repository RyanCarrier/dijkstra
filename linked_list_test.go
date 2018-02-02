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
		if ll.PopOrdered() != verticies[(((i+1)%2)*2)] {
			t.Error("Wrong pop")
		}
		checkNil(t, ll)
		short = !short
	}
}

func TestDig(t *testing.T) {
	verticies := getVertexs()
	short := true
	ll := new(linkedList).init(short)
	checkRootRoot(t, ll)
	checkNil(t, ll)
	ll.PushOrdered(verticies[0])
	if ll.root.next != verticies[0] || ll.root.prev != verticies[0] {
		t.Error("Root next or prev incorrect")
	}
	ll.PopOrdered()
	checkRootRoot(t, ll)
	ll.PushOrdered(verticies[0])
	checkRootNextPrev(t, ll, verticies[0], verticies[0])
	ll.PushOrdered(verticies[1])
	checkRootNextPrev(t, ll, verticies[0], verticies[1])
	ll.PushOrdered(verticies[2])
	checkRootNextPrev(t, ll, verticies[0], verticies[2])
	if ll.root.next.next != ll.root.prev.prev || ll.root.next.next != verticies[1] {
		t.Error("Root directions wrong")
	}
	ll.PopOrdered()
	checkRootNextPrev(t, ll, verticies[0], verticies[1])
	ll.PushOrdered(verticies[1])
	checkRootNextPrev(t, ll, verticies[0], verticies[1])
	if ll.len != 2 {
		ll.print()
		ll.reversePrint()
		t.Errorf("len should be stay on collision")

	}
}

func checkRootNextPrev(t *testing.T, ll *linkedList, n, p *Vertex) {
	if ll.root.next != n {
		t.Errorf("root.next is # should be #\n[%+v]\n[%+v]", ll.root.next, n)
	}
	if ll.root.prev != p {
		t.Errorf("root.prev is # should be #\n[%+v]\n[%+v]", ll.root.prev, p)
	}
}

func checkRootRoot(t *testing.T, ll *linkedList) {
	if ll.root.next != ll.root || ll.root.prev != ll.root {
		t.Error("Root next and prev should be root")
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
