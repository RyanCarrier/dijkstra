package dijkstra

import (
	"fmt"
	"reflect"
	"testing"
)

func TestWrongFormat(t *testing.T) {
	testWrongFormat(t, "testdata/D.txt")
	testWrongFormat(t, "testdata/E.txt")
}

func testWrongFormat(t *testing.T, filename string) {
	_, _, err := Import(filename)
	testErrors(t, ErrWrongFormat, err)
}

func TestCorrectFormat(t *testing.T) {
	test(t, getAGraph(), map[string]int{}, nil, "testdata/A.txt")
}

func TestCorrectFormatNegatives(t *testing.T) {
	test(t, getCGraph(), map[string]int{}, nil, "testdata/C.txt")
}
func TestMixingIntString(t *testing.T) {
	test(t, Graph{}, map[string]int{}, ErrMixMapping, "testdata/H.txt")
}
func TestImportCorrectMap(t *testing.T) {
	wantgraph, wantmap := getGGraph()
	test(t, wantgraph, wantmap, nil, "testdata/G.txt")
}

func test(t *testing.T, wantgraph Graph, wantmap map[string]int, wanterr error, filename string) {
	graph, gmap, err := Import(filename)
	testErrors(t, wanterr, err)
	if !reflect.DeepEqual(gmap, wantmap) {
		t.Error("maps are different",
			"\ngot:\n", fmt.Sprintf("%+v", gmap),
			"\nwant:\n", fmt.Sprintf("%+v", wantmap))
	}
	assertGraphsEqual(t, graph, wantgraph)
}

func testErrors(t *testing.T, wanterr, err error) {
	if wanterr == nil {
		assertErrNil(t, err)
	}
	if err == nil {
		t.Error("err should be nil; ", err.Error())
	}
	if err.Error() != wanterr.Error() {
		t.Error("want:", wanterr.Error(),
			"\ngot:", err.Error())
	}
}
func assertErrNil(t *testing.T, err error) {
	if err != nil {
		t.Error("Error should be nil;" + err.Error())
	}
}

func assertGraphsEqual(t *testing.T, a, b Graph) {
	if len(a.Verticies) != len(b.Verticies) || len(a.Visited) != len(b.Visited) || len(a.Visited) != len(b.Visited) {
		t.Error("Error in graph sizes a size:", len(a.Verticies), "\tb size:", len(b.Verticies))
	}
	for i := range a.Visited {
		if a.Visited[i] != b.Visited[i] {
			t.Error("Index ", i, " not the same visited, a:", a.Visited, "\tb:", b.Visited)
		}
		if !reflect.DeepEqual(a.Verticies[i], b.Verticies[i]) {
			t.Error("Index ", i, " not the same vertex, ",
				"a:\n", fmt.Sprintf("%+v", a.Verticies[i]),
				"\nb:\n", fmt.Sprintf("%+v", b.Verticies[i]))
		}
	}
}

func getAGraph() Graph {
	return Graph{
		make([]bool, 5),
		[]Vertex{
			Vertex{0, 0, map[int]int64{
				1: 4,
				2: 2},
			},
			Vertex{1, 0, map[int]int64{
				3: 2,
				2: 3,
				4: 3},
			},
			Vertex{2, 0, map[int]int64{
				1: 1,
				3: 4,
				4: 5},
			},
			Vertex{3, 0, map[int]int64{}},
			Vertex{4, 0, map[int]int64{
				3: 1},
			},
		},
	}
}

func getBGraph() Graph {
	return Graph{
		make([]bool, 6),
		[]Vertex{
			Vertex{0, 0, map[int]int64{
				1: 4,
				2: 2},
			},
			Vertex{1, 0, map[int]int64{
				3: 2,
				2: 3,
				4: 3},
			},
			Vertex{2, 0, map[int]int64{
				1: 1,
				3: 4,
				4: 5},
			},
			Vertex{3, 0, map[int]int64{
				5: 10}},
			Vertex{4, 0, map[int]int64{
				3: 1},
			},
			Vertex{5, 0, map[int]int64{
				3: 10},
			},
		},
	}
}

func getCGraph() Graph {
	return Graph{
		make([]bool, 6),
		[]Vertex{
			Vertex{0, 0, map[int]int64{
				1: -4,
				2: 2},
			},
			Vertex{1, 0, map[int]int64{
				3: 2,
				2: -3,
				4: 3},
			},
			Vertex{2, 0, map[int]int64{
				1: 1,
				3: 4,
				4: 5},
			},
			Vertex{3, 0, map[int]int64{
				5: -10}},
			Vertex{4, 0, map[int]int64{
				3: 1},
			},
			Vertex{5, 0, map[int]int64{
				3: -10},
			},
		},
	}
}

func getGGraph() (Graph, map[string]int) {
	return Graph{
			make([]bool, 3),
			[]Vertex{
				Vertex{0, 0, map[int]int64{
					1: 2},
				},
				Vertex{0, 0, map[int]int64{
					3: 5},
				},
				Vertex{0, 0, map[int]int64{
					0: 1,
					1: 1},
				},
			},
		}, map[string]int{
			"A": 0,
			"B": 1,
			"C": 2,
		}
}
