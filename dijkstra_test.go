package dijkstra

import (
	"errors"
	"reflect"
	"slices"
	"strconv"
	"sync"
	"testing"
)

func TestErrors(t *testing.T) {
	t.Run("ErrNoPath", func(t *testing.T) {
		graph := Graph{
			[]Vertex{
				map[int]uint64{1: 2},
				map[int]uint64{2: 3},
				map[int]uint64{3: 4},
				map[int]uint64{2: 5},
				map[int]uint64{},
			},
		}
		_, err := graph.Shortest(0, 4)
		if err != ErrNoPath {
			testErrors(t, ErrNoPath, err, 0)
		}
	})
	t.Run("ErrLoop", func(t *testing.T) {
		graph := Graph{
			[]Vertex{
				map[int]uint64{1: 1, 2: 0},
				map[int]uint64{2: 5},
				map[int]uint64{1: 10, 3: 10},
				map[int]uint64{},
			},
		}
		_, err := graph.Longest(0, 3)
		testErrors(t, ErrLoopDetected, err, 0)
	})
	t.Run("ErrExists", func(t *testing.T) {
		graph := NewGraph()
		graph.AddEmptyVertex(0)
		graph.AddEmptyVertex(1)
		graph.AddArc(0, 1, 1)
		graph.AddEmptyVertex(10)
		err := graph.AddEmptyVertex(5)
		if err != nil {
			t.Error("Got error:", err)
		}
		err = graph.AddEmptyVertex(5)
		if !errors.Is(err, ErrVertexAlreadyExists) {
			t.Errorf("Got error: %v\nExpected error:%v", err, ErrVertexAlreadyExists)
		}
	})
}

func TestCorrect(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		g := NewGraph()
		for range 4 {
			g.AddNewEmptyVertex()
		}
		g.AddArc(0, 1, 1)
		g.AddArc(0, 2, 2)
		g.AddArc(1, 3, 1)
		g.AddArc(2, 3, 1)
		t.Run("Shortest", func(t *testing.T) {
			path, err := g.Shortest(0, 3)
			if err != nil {
				t.Fatal("Error in simple graph: ", err)
			}
			if path.Distance != 2 {
				t.Fatal("Incorrect distance in simple graph: ", path.Distance)
			}
			if !reflect.DeepEqual(path.Path, []int{0, 1, 3}) {
				t.Fatal("Incorrect path in simple graph: ", path.Path)
			}
		})
		t.Run("Longest", func(t *testing.T) {
			path, err := g.Longest(0, 3)
			if err != nil {
				t.Fatal("Error in simple graph: ", err)
			}
			if path.Distance != 3 {
				t.Fatal("Incorrect distance in simple graph: ", path.Distance)
			}
			if !reflect.DeepEqual(path.Path, []int{0, 2, 3}) {
				t.Fatal("Incorrect path in simple graph: ", path.Path)
			}
		})
	})
	t.Run("All", func(t *testing.T) {
		for i, test := range testGraphsCorrect {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				t.Run("Shortest", func(t *testing.T) {
					got, err := test.graph.ShortestAll(test.from, test.to)
					testErrors(t, test.evalErr, err, i)
					testResults(t, test.shortestSolution, got.SmallestPath(), true, i)
				})
				t.Run("Longest", func(t *testing.T) {
					got, err := test.graph.LongestAll(test.from, test.to)
					testErrors(t, test.evalErr, err, i)
					testResults(t, test.longestSolution, got.SmallestPath(), true, i)
				})
			})
		}

		shortText := []string{"Short", "Long_"}
		for i, text := range shortText {
			t.Run(text, func(t *testing.T) {
				g := NewGraph()
				for range 10 {
					g.AddNewEmptyVertex()
				}
				g.AddArc(0, 1, 1)
				g.AddArc(1, 2, 2)
				g.AddArc(1, 3, 2)
				g.AddArc(1, 4, 2)
				g.AddArc(2, 5, 1)
				g.AddArc(3, 5, 1)
				g.AddArc(4, 5, 1)
				g.AddArc(5, 6, 2)
				g.AddArc(5, 7, 2)
				g.AddArc(5, 8, 2)
				g.AddArc(6, 9, 1)
				g.AddArc(7, 9, 1)
				g.AddArc(8, 9, 1)
				g.AddArc(0, 9, 7)
				var err error
				var result BestPaths[int]
				if i == 0 {
					result, err = g.ShortestAll(0, 9)
				} else {
					result, err = g.LongestAll(0, 9)
				}
				if err != nil {
					t.Error("Error in simple graph: ", err)
				}
				expected := BestPaths[int]{
					Distance: 7,
					Paths: [][]int{
						{0, 9},
						//
						{0, 1, 4, 5, 8, 9},
						{0, 1, 3, 5, 8, 9},
						{0, 1, 2, 5, 8, 9},
						//
						{0, 1, 4, 5, 7, 9},
						{0, 1, 3, 5, 7, 9},
						{0, 1, 2, 5, 7, 9},
						//
						{0, 1, 4, 5, 6, 9},
						{0, 1, 3, 5, 6, 9},
						{0, 1, 2, 5, 6, 9},
						//
					}}
				if expected.Distance != result.Distance {
					t.Error("Incorrect distance in simple graph: ", result.Distance)
				}
				for _, path := range result.Paths {
					if !slices.ContainsFunc(expected.Paths, func(a []int) bool {
						if len(a) != len(path) {
							return false
						}
						for j := range len(a) {
							if a[j] != path[j] {
								return false
							}
						}
						return true
					}) {
						t.Errorf("Missing expected path in result: %v\ngot:%v", path, result.Paths)
					}
				}
			})
		}
	})

	t.Run("Default", func(t *testing.T) {
		var err error
		var got BestPath[int]
		t.Run("Shortest", func(t *testing.T) {
			for i, test := range testGraphsCorrect {
				got, err = test.graph.Shortest(test.from, test.to)
				testErrors(t, test.evalErr, err, i)
				testResults(t, test.shortestSolution, got, true, i)
			}
		})
		t.Run("Longest", func(t *testing.T) {
			for i, test := range testGraphsCorrect {
				got, err = test.graph.Longest(test.from, test.to)
				testErrors(t, test.evalErr, err, i)
				testResults(t, test.longestSolution, got, false, i)
			}
		})
	})
	t.Run("ForceList", func(t *testing.T) {
		var err error
		var got BestPath[int]
		var lists = [4]int{listShortPQ, listLongPQ, listShortLL, listLongLL}
		var titles = [4]string{"listShortPQ", "listLongPQ", "listShortLL", "listLongLL"}
		for i := range lists {
			t.Run(titles[i], func(t *testing.T) {
				t.Run("Shortest", func(t *testing.T) {
					for i, test := range testGraphsCorrect {
						got, err = test.graph.evaluate(test.from, test.to, true, lists[i])
						testErrors(t, test.evalErr, err, i)
						testResults(t, test.shortestSolution, got, true, i)
					}
				})
				t.Run("Longest", func(t *testing.T) {
					for i, test := range testGraphsCorrect {
						got, err = test.graph.evaluate(test.from, test.to, false, lists[i])
						testErrors(t, test.evalErr, err, i)
						testResults(t, test.longestSolution, got, false, i)
					}
				})
			})
		}
	})
	// t.Run("SolutionsAll", testCorrectSolutionsAll)
	t.Run("WorstCaseGenerated", func(t *testing.T) {
		for _, nodeAmount := range []int{3, 8, 16, 64} {
			t.Run(strconv.Itoa(nodeAmount)+"Nodes", func(t *testing.T) {
				gshort, expectedShort := generateWorstCaseShortest(nodeAmount)
				glong, expectedLong := generateWorstCaseLongest(nodeAmount)
				from, to := 0, nodeAmount-1
				resultShort, err := gshort.Shortest(from, to)
				resultLong, err := glong.Longest(from, to)
				tests := []struct {
					result, expected BestPath[int]
					graph            Graph
				}{
					{resultShort, expectedShort, gshort},
					{resultLong, expectedLong, glong},
				}
				for i, test := range tests {
					t.Run([]string{"Shortest", "Longest"}[i], func(t *testing.T) {
						if err != nil {
							str, _ := test.graph.ToString()
							t.Errorf("Run had error: %v\nfrom, to = %d, %d\n%s", err, from, to, str)
						}
						if test.result.Distance != test.expected.Distance {
							str, _ := test.graph.ToString()
							t.Errorf("Run had incorrect distance;(%d->%d)\n\texpected:%d\n\tgot:%d\n%s",
								from, to, test.expected.Distance, test.result.Distance, str)
						}
					})
				}
			})
		}
	})
	t.Run("Concurrent", func(t *testing.T) {
		baseNodes := 100
		threads := 16
		graph, expected := generateWorstCaseShortest(baseNodes)
		var results = make([]BestPath[int], threads)
		wg := sync.WaitGroup{}
		wg.Add(threads)
		for i := range threads {
			go func(i int) {
				results[i], _ = graph.Shortest(0, baseNodes-1)
				wg.Done()
			}(i)
		}
		wg.Wait()
		for _, result := range results {
			if result.Distance != expected.Distance {
				t.Error("Incorrect distance", result.Distance, "!=", expected.Distance)
			}
			if !slices.Equal(result.Path, expected.Path) {
				t.Error("Incorrect path", result.Path, "!=", expected.Path)
			}
		}

	})
	t.Run("SequentialRuns", func(t *testing.T) {
		nodeAmount := 10
		g, _ := generateWorstCaseShortest(nodeAmount)
		from, to := 0, nodeAmount-1
		initialResult, _ := g.Shortest(from, to)
		for i := 0; i < 10; i++ {
			result, err := g.Shortest(from, to)
			if err != nil {
				t.Error("Sequential runs had error: ", err)
			}
			if initialResult.Distance != result.Distance {
				t.Error("Sequential runs are not equal (distance) ", initialResult.Distance, result.Distance)
			}
			if !reflect.DeepEqual(initialResult.Path, result.Path) {
				t.Error("Sequential runs are not equal (path) ", initialResult.Path, result.Path)
			}
		}
		for i := 0; i < 10; i++ {
			result, err := g.Shortest(0, nodeAmount-1)
			if err != nil {
				t.Error("Sequential runs had error: ", err)
			}
			if initialResult.Distance != result.Distance {
				t.Error("Sequential runs are not equal (distance) ", initialResult.Distance, result.Distance)
			}
			if !reflect.DeepEqual(initialResult.Path, result.Path) {
				t.Error("Sequential runs are not equal (path) ", initialResult.Path, result.Path)
			}
		}
	})
}

func BenchmarkLists(b *testing.B) {
	var lists = [4]int{listShortPQ, listLongPQ, listShortLL, listLongLL}
	var titles = [4]string{"listShortPQ", "listLongPQ", "listShortLL", "listLongLL"}
	nodeIterations := 6
	shortest := false
	shortText := []string{"Short", "Long_"}
	var g Graph
	for ci, worstCase := range []bool{false, true} {
		b.Run([]string{"Default", "Worst__"}[ci], func(b *testing.B) {
			for z := range 2 {
				shortest = !shortest
				nodes := 1
				for j := 0; j < nodeIterations; j++ {
					nodes *= 4
					if worstCase {
						if nodes > 200 {
							continue
						}
						if shortest {
							g, _ = generateWorstCaseShortest(nodes)
						} else {
							g, _ = generateWorstCaseShortest(nodes)
						}
					} else {
						g = Generate(nodes)
					}
					for j, list := range lists {
						b.Run(shortText[z]+"/"+strconv.Itoa(nodes)+"Nodes"+"/"+titles[j], func(b *testing.B) {
							b.ResetTimer()
							for i := 0; i < b.N; i++ {
								g.evaluate(0, len(g.vertexArcs)-1, shortest, list)
							}
						})
					}
				}
			}
		})
	}
}

func BenchmarkAll(b *testing.B) {
	nodeIterations := 6
	nodes := 1
	for j := 0; j < nodeIterations; j++ {
		nodes *= 4
		g := Generate(nodes)
		b.Run(strconv.Itoa(nodes)+"Nodes", func(b *testing.B) {
			src, dest := 0, len(g.vertexArcs)-1
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				g.Shortest(src, dest)
			}
		})
	}
}

func testResults(t *testing.T, expected, got BestPath[int], shortest bool, testIndex int) {
	distmethod := "Shortest"
	if !shortest {
		distmethod = "Longest"
	}
	err := func(label string) {
		t.Errorf("%s %s\ntest %d\n got dist:%d\nwant dist:%d\n got path:%v\nwant path:%v", distmethod, label, testIndex, got.Distance, expected.Distance, got.Path, expected.Path)
	}
	if got.Distance != expected.Distance {
		err("distance incorrect")
	} else if !reflect.DeepEqual(got.Path, expected.Path) {
		err("path incorrect")
	}
}

type testGraph struct {
	stringRepresentation string
	graph                Graph
	from                 int
	to                   int
	shortestSolution     BestPath[int]
	longestSolution      BestPath[int]
	importErr            error
	evalErr              error
}
type testMappedGraph[T comparable] struct {
	stringRepresentation string
	graph                MappedGraph[T]
	from                 string
	to                   string
	shortestSolution     BestPath[T]
	longestSolution      BestPath[T]
	importErr            error
	evalErr              error
}

var testGraphsCorrect = []testGraph{
	{
		`0 1,4 2,2
1 3,2 4,3
2 1,1 3,4 4,4
3 5,10
4 3,1
5 3,10`,
		Graph{
			[]Vertex{
				{1: 4, 2: 2},
				{3: 2, 4: 3},
				{1: 1, 3: 4, 4: 4},
				{5: 10},
				{3: 1},
				{3: 10},
			},
		},
		0, 5,
		BestPath[int]{15, []int{0, 2, 1, 3, 5}},
		BestPath[int]{18, []int{0, 1, 4, 3, 5}},
		nil, nil,
	},
	{
		`0 2,1 3,1 1,10
1 3,10 2,1 3,1
2 4,10
3 2,10 4,1
4`,
		Graph{
			[]Vertex{
				{1: 10, 2: 1, 3: 1},
				{2: 1, 3: 1},
				{4: 10},
				{2: 10, 4: 1},
				{},
			},
		},
		0, 4,
		BestPath[int]{2, []int{0, 3, 4}},
		//0 1 2 4
		BestPath[int]{31, []int{0, 1, 3, 2, 4}},
		nil, nil,
	},
}

var testMappedGraphs = []testMappedGraph[string]{
	{
		`A B,4 C,2
B D,2 C,3 E,3
C B,1 D,4 E,5
D F,10
E D,1
F D,10`,
		MappedGraph[string]{
			graph: Graph{
				[]Vertex{
					{1: 4, 2: 2},
					{3: 2, 2: 3, 4: 3},
					{1: 1, 3: 4, 4: 5},
					{5: 10},
					{3: 1},
					{3: 10},
				},
			},
			mapping: map[string]int{"A": 0, "B": 1, "C": 2, "D": 3, "E": 4, "F": 5},
		},
		"0", "5",
		BestPath[string]{
			Distance: 15,
			Path:     []string{"0", "2", "1", "3", "5"},
		},
		BestPath[string]{
			Distance: 21,
			Path:     []string{"0", "1", "2", "3", "5"},
		},
		nil, nil,
	},
	{
		`A C,1 D,1 B,10
B D,10 C,1 D,1
C E,10
D C,10 E,1
E`,
		MappedGraph[string]{
			graph: Graph{
				[]Vertex{
					{1: 1, 2: 1, 3: 10},
					{4: 10},
					{1: 10, 4: 1},
					{1: 1, 2: 1},
					{},
				},
				//map looks weird as for importing we assign mapped values to when they
				// are found, and it checks arcs during each vertex
				// validation is checked that all arcs have a destination though
			},
			mapping: map[string]int{"A": 0, "B": 3, "C": 1, "D": 2, "E": 4},
		},
		"0", "4",
		BestPath[string]{
			2,
			[]string{"0", "3", "4"},
		},
		BestPath[string]{
			31,
			[]string{"0", "1", "3", "2", "4"},
		},
		nil, nil,
	},
}
