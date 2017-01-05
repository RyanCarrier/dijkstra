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
	_, err := Import(filename)
	testErrors(t, ErrWrongFormat, err, filename)
}

func TestCorrectFormat(t *testing.T) {
	test(t, getAGraph(), map[string]int{}, nil, "testdata/A.txt")
}

func TestCorrectFormatNegatives(t *testing.T) {
	test(t, getCGraph(), map[string]int{}, nil, "testdata/C.txt")
}
func TestMixingIntString(t *testing.T) {
	filename := "testdata/H.txt"
	_, err := Import(filename)
	testErrors(t, ErrMixMapping, err, filename)
}
func TestImportCorrectMap(t *testing.T) {
	wantgraph, wantmap := getGGraph()
	test(t, wantgraph, wantmap, nil, "testdata/G.txt")
}

func test(t *testing.T, wantgraph Graph, wantmap map[string]int, wanterr error, filename string) {
	graph, err := Import(filename)
	gmap := graph.mapping
	testErrors(t, wanterr, err, filename)
	if !reflect.DeepEqual(gmap, wantmap) {
		t.Fatal("maps are different",
			"\ngot:\n", fmt.Sprintf("%+v", gmap),
			"\nwant:\n", fmt.Sprintf("%+v", wantmap))
	}
	assertGraphsEqual(t, graph, wantgraph)
}

//func assertMaps(t *testing.T, got, want map[string]int)

func testErrors(t *testing.T, wanterr, err error, filename string) {
	if wanterr == nil {
		assertErrNil(t, err, filename)
		return
	}
	if err == nil {
		t.Fatal("err should not be nil, want; ", wanterr.Error())
	}
	if err.Error() != wanterr.Error() {
		t.Fatal("want:", wanterr.Error(),
			"\ngot:", err.Error())
	}
}
func assertErrNil(t *testing.T, err error, filename string) {
	if err != nil {
		t.Fatal("Error should be nil;\n", filename, "\n"+err.Error())
	}
}

func assertGraphsEqual(t *testing.T, a, b Graph) {
	if len(a.Verticies) != len(b.Verticies) {
		t.Fatal("Error in graph sizes a size:", len(a.Verticies), "\tb size:", len(b.Verticies))
	}
}

func getAGraph() Graph {
	return Graph{
		0, false,
		[]Vertex{
			Vertex{0, 0, 0, map[int]int64{
				1: 4,
				2: 2},
			},
			Vertex{1, 0, 0, map[int]int64{
				3: 2,
				2: 3,
				4: 3},
			},
			Vertex{2, 0, 0, map[int]int64{
				1: 1,
				3: 4,
				4: 5},
			},
			Vertex{3, 0, 0, map[int]int64{}},
			Vertex{4, 0, 0, map[int]int64{
				3: 1},
			},
		},
		newLinkedList(),
		map[string]int{},
		false,
		0,
	}
}

func getBGraph() Graph {
	return Graph{
		0, false,
		[]Vertex{
			Vertex{0, 0, 0, map[int]int64{
				1: 4,
				2: 2},
			},
			Vertex{1, 0, 0, map[int]int64{
				3: 2,
				2: 3,
				4: 3},
			},
			Vertex{2, 0, 0, map[int]int64{
				1: 1,
				3: 4,
				4: 5},
			},
			Vertex{3, 0, 0, map[int]int64{
				5: 10}},
			Vertex{4, 0, 0, map[int]int64{
				3: 1},
			},
			Vertex{5, 0, 0, map[int]int64{
				3: 10},
			},
		},
		newLinkedList(),
		map[string]int{},
		false,
		0,
	}
}

func getBSol() BestPath {
	return BestPath{
		Distance: 15,
		Path:     []int{0, 2, 1, 3, 5},
	}
}

func getCGraph() Graph {
	return Graph{0, false,
		[]Vertex{
			Vertex{0, 0, 0, map[int]int64{
				1: -4,
				2: 2},
			},
			Vertex{1, 0, 0, map[int]int64{
				3: 2,
				2: -3,
				4: 3},
			},
			Vertex{2, 0, 0, map[int]int64{
				1: 1,
				3: 4,
				4: 5},
			},
			Vertex{3, 0, 0, map[int]int64{
				5: -10}},
			Vertex{4, 0, 0, map[int]int64{
				3: 1},
			},
			Vertex{5, 0, 0, map[int]int64{
				3: -10},
			},
		},
		newLinkedList(),
		map[string]int{},
		false,
		0,
	}
}

func getGGraph() (Graph, map[string]int) {
	return Graph{
			0, false,
			[]Vertex{
				Vertex{0, 0, 0, map[int]int64{
					1: 2},
				},
				Vertex{1, 0, 0, map[int]int64{
					2: 5},
				},
				Vertex{2, 0, 0, map[int]int64{
					0: 1,
					1: 1},
				},
			},
			newLinkedList(),
			map[string]int{},
			false,
			0,
		}, map[string]int{
			"A": 0,
			"B": 1,
			"C": 2,
		}
}
