package dijkstra

import (
	"reflect"
	"testing"
)

func TestFailure(t *testing.T) {
	testSolution(t, BestPath{}, ErrNoPath, "testdata/I.txt", 0, 4)
}

func testSolution(t *testing.T, best BestPath, wanterr error, filename string, from, to int) {
	graph, _, _ := Import(filename)
	got, err := graph.Shortest(from, to)
	testErrors(t, wanterr, err)
	if got.Distance != best.Distance {
		t.Error("Distance incorrect\ngot: ", got.Distance, "\nwant: ", best.Distance)
	}
	if !reflect.DeepEqual(got.Path, best.Path) {
		t.Error("Path incorrect\ngot: ", got.Path, "\nwant: ", best.Path)
	}
}
