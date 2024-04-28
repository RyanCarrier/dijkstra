package dijkstra

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestWrongFormat(t *testing.T) {
	importDatas := []string{
		`
0 1,12
1 0 12
`, `
0 1,12
1 0,1Z2
		`}

	for i, importData := range importDatas {
		_, err := Import(importData)
		testErrors(t, ErrWrongFormat, err, i)
	}
}

func TestImport(t *testing.T) {
	t.Run("Correct", func(t *testing.T) {
		for i, test := range testGraphsCorrect {
			graph, err := Import(test.stringRepresentation)
			testErrors(t, nil, err, i)
			assertGraphsEqual(t, graph, test.graph, i)
		}
		t.Run("Mapped", func(t *testing.T) {
			for i, test := range testMappedGraphs {
				graph, err := ImportStringMapped(test.stringRepresentation)
				testErrors(t, nil, err, i)
				if !reflect.DeepEqual(graph.mapping, test.graph.mapping) {
					t.Fatal("maps are different (test ", i, ")",
						"\ngot:\n", fmt.Sprintf("%+v", graph.mapping),
						"\nwant:\n", fmt.Sprintf("%+v", test.graph.mapping))
				}
				assertGraphsEqual(t, graph.graph, test.graph.graph, i)
			}
		})
	})
}

type typeOptions[T comparable] struct {
	typeName string
	a        T
	b        T
	c        T
}
type customTestType struct {
	name string
	age  int
}

func (c customTestType) Less(other customTestType) int {
	ageCmp := c.age - other.age
	if ageCmp != 0 {
		return ageCmp
	}
	return strings.Compare(c.name, other.name)
}

func TestMaps(t *testing.T) {
	byteOption := typeOptions[byte]{"byte", 'a', 'b', 'c'}
	strOption := typeOptions[string]{"byte", "a", "b", "c"}
	floatOption := typeOptions[float32]{"byte", 0.1, 0.2, 0.3}
	float64Option := typeOptions[float64]{"byte", 0.1, 0.2, 0.3}
	customOption := typeOptions[customTestType]{
		"customTestType",
		customTestType{"a", 5},
		customTestType{"b", 1},
		customTestType{"c", 10},
	}
	testMap(t, byteOption)
	testMap(t, strOption)
	testMap(t, floatOption)
	testMap(t, float64Option)
	testMap(t, customOption)
}
func testMap[T comparable](t *testing.T, test typeOptions[T]) {
	var err error
	mappedGraph := NewMappedGraph[T]()
	err = mappedGraph.AddEmptyVertex(test.a)
	if err != nil {
		t.Fatal(err)
	}
	err = mappedGraph.AddEmptyVertex(test.b)
	if err != nil {
		t.Fatal(err)
	}
	err = mappedGraph.AddEmptyVertex(test.b)
	if err == nil {
		t.Fatal(err)
	}
	err = mappedGraph.AddEmptyVertex(test.c)
	if err != nil {
		t.Fatal(err)
	}
	err = mappedGraph.AddArc(test.a, test.b, 5)
	err = mappedGraph.AddArc(test.a, test.c, 10)
	err = mappedGraph.AddArc(test.c, test.b, 5)
	if err != nil {
		t.Fatalf("AddArc fail, err:%v, mapping:%v", err, mappedGraph.mapping)
	}
	result, err := mappedGraph.Shortest(test.a, test.b)
	if err != nil {
		t.Fatal(err)
	}
	if result.Distance != 5 {
		t.Fatal("wrong distance", result.Distance)
	}
	if len(result.Path) != 2 || result.Path[0] != test.a || result.Path[1] != test.b {
		t.Fatal("wrong path", result.Path)
	}
}

func TestToString(t *testing.T) {
	t.Run("Correct", func(t *testing.T) {
		testToString(t, testGraphsCorrect)
	})
	t.Run("Mapped", func(t *testing.T) {
		testToStringMapped(t, testMappedGraphs)
	})
}
func testToString(t *testing.T, testData []testGraph) {
	for i, test := range testData {
		got, _ := Import(test.stringRepresentation)
		result, err := got.Export()
		if err != nil {
			t.Fatal("Error in test", i, "\n"+err.Error())
		}
		imported, _ := Import(result)
		assertGraphsEqual(t, imported, test.graph, i)
	}
}

func testToStringMapped(t *testing.T, testData []testMappedGraph[string]) {
	for i, test := range testData {
		got, _ := Import(test.stringRepresentation)
		result, err := got.Export()
		if err != nil {
			t.Fatal("Error in test", i, "\n"+err.Error())
		}
		imported, _ := ImportStringMapped(result)
		result2, err := imported.Export()
		if err != nil {
			t.Fatal("Error in test", i, "\n"+err.Error())
		}
		if strings.TrimSpace(result) != strings.TrimSpace(result2) {
			t.Fatal("Error in test", i, "\n"+result, "\n", result2)
		}

	}
}

func testErrors(t *testing.T, wanterr, err error, testIndex int) {
	if wanterr == nil {
		if err != nil {
			t.Fatal("Error should be nil;\ntest ", testIndex, "\n"+err.Error())
		}
		return
	}
	if err == nil {
		t.Fatal("err should not be nil, want; ", wanterr.Error())
	}
	if !errors.Is(err, wanterr) {
		t.Fatal(
			"\nwant:", wanterr.Error(),
			"\n got:", err.Error(),
		)
	}
}

func assertGraphsEqual(t *testing.T, got, want Graph, i int) {
	if len(got.vertexArcs) != len(want.vertexArcs) {
		t.Fatal("Error in graph sizes (test", i, ") a size:", len(got.vertexArcs), ",b size:", len(want.vertexArcs))
	}
	for vi := range got.vertexArcs {
		if !reflect.DeepEqual(got.vertexArcs[vi], want.vertexArcs[vi]) {
			t.Fatal("Error in graph vertexArcs (test", i, "index", vi, ")\n got:", got.vertexArcs[vi], "\nwant:", want.vertexArcs[vi])
		}
	}
}
