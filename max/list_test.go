//From github.com/RyanCarrier/dijkstra
package max

import "testing"

//if the other shit is passing this is working lol

func TestLazyInit(t *testing.T) {
	ll := linkedList{root: element{}, len: 0}
	if ll.root.next != nil || ll.root.prev != nil {
		t.Error("pre init should have root failing")
	}
	ll.lazyinit()
	if ll.root.next == nil || ll.root.prev == nil {
		t.Error("post init should have root passing")
	}
}

func TestEmptyList(t *testing.T) {
	ll := linkedList{root: element{}, len: 0}
	if ll.front() != nil || ll.back() != nil {
		t.Error("empty list front()/back() should be nil")
	}
}
