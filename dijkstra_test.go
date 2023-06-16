package dijkstra

import (
	"log"
	"os"
	"reflect"
	"strconv"
	"sync"
	"testing"
)

func TestErrors(t *testing.T) {
	t.Run("ErrNoPath", func(t *testing.T) {
		testSolution(t, BestPath{}, ErrNoPath, "testdata/I.txt", 0, 4, true, -1)
	})
	t.Run("ErrLoop", func(t *testing.T) {
		testSolution(t, BestPath{}, newErrLoop(2, 1), "testdata/J.txt", 0, 4, true, -1)
	})
	t.Run("ErrAlreadyCalculating", func(t *testing.T) {
		nodeAmount := 5
		g := GenerateWorstCase(nodeAmount)
		g.running = true
		_, err := g.Shortest(0, nodeAmount-1)
		if err != ErrAlreadyCalculating {
			t.Error("gotErr:", err, "wantErr:", ErrAlreadyCalculating)
		}
	})

}

func TestCorrect(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		testSolution(t, getBSol(), nil, "testdata/B.txt", 0, 5, true, -1)
		testSolution(t, getKSolLong(), nil, "testdata/K.txt", 0, 4, false, -1)
		testSolution(t, getKSolShort(), nil, "testdata/K.txt", 0, 4, true, -1)
	})
	t.Run("SolutionsAll", testCorrectSolutionsAll)
	t.Run("Lists", func(t *testing.T) {
		for i := 0; i <= 3; i++ {
			testSolution(t, getBSol(), nil, "testdata/B.txt", 0, 5, true, i)
			testSolution(t, getKSolLong(), nil, "testdata/K.txt", 0, 4, false, i)
			testSolution(t, getKSolShort(), nil, "testdata/K.txt", 0, 4, true, i)
		}
	})
	t.Run("AutoLargeList", func(t *testing.T) {
		g := NewGraph()
		for i := 0; i < 2000; i++ {
			v := g.AddNewVertex()
			v.AddArc(i+1, 1)
		}
		g.AddNewVertex()
		_, err := g.Shortest(0, 2000)
		testErrors(t, nil, err, "manual test")
		_, err = g.Longest(0, 2000)
		testErrors(t, nil, err, "manual test")
		_, err = g.ShortestSafe(0, 2000)
		testErrors(t, nil, err, "manual test")
		_, err = g.LongestSafe(0, 2000)
		testErrors(t, nil, err, "manual test")
	})
	t.Run("Concurrent", testConcurrentSolving)
	t.Run("SequentialRuns", func(t *testing.T) {
		nodeAmount := 10
		g := GenerateWorstCase(nodeAmount)
		initialResult, _ := g.Shortest(0, nodeAmount-1)
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
		for i := 0; i < 10; i++ {
			result, err := g.ShortestSafe(0, nodeAmount-1)
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
func testConcurrentSolving(t *testing.T) {
	baseNodes := 100
	graphAmount := 8
	graphs := []Graph{}
	for i := 0; i < graphAmount; i++ {
		graphs = append(graphs, GenerateWorstCase(baseNodes+i))
	}
	var results = make([]BestPath, len(graphs))
	wg := sync.WaitGroup{}
	wg.Add(len(graphs))
	for i, g := range graphs {
		go func(i int, g Graph) {
			results[i], _ = g.Shortest(0, baseNodes-1+i)
			wg.Done()
		}(i, g)
	}
	wg.Wait()
	for i, result := range results {
		if len(result.Path) != baseNodes+i {
			t.Error("Incorrect path length", len(result.Path), "!=", baseNodes+i)
		}
		for j, stepNode := range result.Path {
			if stepNode != j {
				t.Error("Incorrect path step", stepNode, "!=", j)
			}
		}
	}

}
func testCorrectSolutionsAll(t *testing.T) {
	graph := NewGraph()
	//Add the 3 verticies
	graph.AddVertex(0)
	graph.AddVertex(1)
	graph.AddVertex(2)
	graph.AddVertex(3)

	//Add the arcs
	graph.AddArc(0, 1, 1)
	graph.AddArc(0, 2, 1)
	graph.AddArc(1, 3, 0)
	graph.AddArc(2, 3, 0)
	testGraphSolutionAll(t, BestPaths{BestPath{1, []int{0, 2, 3}}, BestPath{1, []int{0, 1, 3}}}, nil, *graph, 0, 3, true)

	graph = NewGraph()
	graph.AddVertex(0)
	graph.AddVertex(1)
	graph.AddVertex(2)
	graph.AddVertex(3)
	graph.AddVertex(4)
	graph.AddVertex(5)
	graph.AddArc(0, 1, 1)
	graph.AddArc(0, 3, 1)
	graph.AddArc(1, 2, 1)
	graph.AddArc(1, 4, 1)
	graph.AddArc(2, 5, 1)
	graph.AddArc(3, 4, 1)
	graph.AddArc(4, 2, 1)
	graph.AddArc(4, 5, 1)

	testGraphSolutionAll(t, BestPaths{BestPath{3, []int{0, 3, 4, 5}}, BestPath{3, []int{0, 1, 4, 5}}, BestPath{3, []int{0, 1, 2, 5}}}, nil, *graph, 0, 5, true)
}

var benchNames = []string{"github.com/RyanCarrier-ALL", "github.com/RyanCarrier", "github.com/RyanCarrierSafe"}
var listNames = []string{"PQShort", "PQLong", "LLShort", "LLLong"}

func BenchmarkSetup(b *testing.B) {
	nodeIterations := 6
	nodes := 1
	for j := 0; j < nodeIterations; j++ {
		nodes *= 4
		b.Run("setup/"+strconv.Itoa(nodes)+"Nodes", func(b *testing.B) {
			filename := "testdata/bench/" + strconv.Itoa(nodes) + ".txt"
			if _, err := os.Stat(filename); err != nil {
				g := Generate(nodes)
				err := g.ExportToFile(filename)
				if err != nil {
					log.Fatal(err)
				}
			}
			g, _ := Import(filename)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				g.setup(true, 0, -1)
			}
		})
	}
}

func BenchmarkLists(b *testing.B) {
	nodeIterations := 6
	shortest := false
	shortText := []string{"Short", "Long"}
	for z := 0; z < 2; z++ {
		shortest = !shortest
		for i, n := range listNames {
			nodes := 1
			for j := 0; j < nodeIterations; j++ {
				nodes *= 4
				b.Run(shortText[z]+"/"+n+"/"+strconv.Itoa(nodes)+"Nodes", func(b *testing.B) {
					benchmarkList(b, nodes, i, shortest)
				})
			}
		}
	}
}

func benchmarkList(b *testing.B, nodes, list int, shortest bool) {

	filename := "testdata/bench/" + strconv.Itoa(nodes) + ".txt"
	if _, err := os.Stat(filename); err != nil {
		g := Generate(nodes)
		err := g.ExportToFile(filename)
		if err != nil {
			log.Fatal(err)
		}
	}
	graph, _ := Import(filename)
	src, dest := 0, len(graph.Verticies)-1
	//====RESET TIMER BEFORE LOOP====
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.setup(shortest, src, list)
		graph.postSetupEvaluate(src, dest, shortest)
	}
}

func BenchmarkAll(b *testing.B) {
	nodeIterations := 6
	for i, n := range benchNames {
		nodes := 1
		for j := 0; j < nodeIterations; j++ {
			nodes *= 4
			b.Run(n+"/"+strconv.Itoa(nodes)+"Nodes", func(b *testing.B) {
				benchmarkAlt(b, nodes, i)
			})

		}
	}
	//Cleanup
	nodes := 1
	for j := 0; j < nodeIterations; j++ {
		nodes *= 4
		os.Remove("testdata/bench/" + strconv.Itoa(nodes) + ".txt")
	}
}

func benchmarkAlt(b *testing.B, nodes, i int) {
	filename := "testdata/bench/" + strconv.Itoa(nodes) + ".txt"
	if _, err := os.Stat(filename); err != nil {
		g := Generate(nodes)
		err := g.ExportToFile(filename)
		if err != nil {
			log.Fatal(err)
		}
	}
	switch i {
	case 0:
		benchmarkRCall(b, filename)
	case 1:
		benchmarkRC(b, filename)
	case 2:
		benchmarkRCa(b, filename)
	default:
		b.Error("You're retarded")
	}
}

func benchmarkRC(b *testing.B, filename string) {
	graph, _ := Import(filename)
	src, dest := 0, len(graph.Verticies)-1
	//====RESET TIMER BEFORE LOOP====
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.Shortest(src, dest)
	}
}
func benchmarkRCa(b *testing.B, filename string) {
	graph, _ := Import(filename)
	src, dest := 0, len(graph.Verticies)-1
	//====RESET TIMER BEFORE LOOP====
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.ShortestSafe(src, dest)
	}
}
func benchmarkRCall(b *testing.B, filename string) {
	graph, _ := Import(filename)
	src, dest := 0, len(graph.Verticies)-1
	//====RESET TIMER BEFORE LOOP====
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.ShortestAll(src, dest)
	}
}

func testSolution(t *testing.T, best BestPath, wanterr error, filename string, from, to int, shortest bool, list int) {
	var err error
	var graph Graph
	graph, err = Import(filename)
	if err != nil {
		t.Fatal(err, filename)
	}
	var got BestPath
	var gotAll BestPaths
	if list >= 0 {
		graph.setup(shortest, from, list)
		got, err = graph.postSetupEvaluate(from, to, shortest)
	} else if shortest {
		got, err = graph.Shortest(from, to)
	} else {
		got, err = graph.Longest(from, to)
	}
	testErrors(t, wanterr, err, filename)
	testResults(t, got, best, shortest, filename)
	if shortest {
		got, err = graph.ShortestSafe(from, to)
	} else {
		got, err = graph.LongestSafe(from, to)
	}
	testErrors(t, wanterr, err, filename)
	testResults(t, got, best, shortest, filename)
	if list >= 0 {
		graph.setup(shortest, from, list)
		gotAll, err = graph.postSetupEvaluateAll(from, to, shortest)
	} else if shortest {
		gotAll, err = graph.ShortestAll(from, to)
	} else {
		gotAll, err = graph.LongestAll(from, to)
	}
	testErrors(t, wanterr, err, filename)
	if len(gotAll) == 0 {
		gotAll = BestPaths{BestPath{}}
	}
	testResults(t, gotAll[0], best, shortest, filename)
}

func testGraphSolutionAll(t *testing.T, best BestPaths, wanterr error, graph Graph, from, to int, shortest bool) {
	var err error
	var gotAll BestPaths
	if shortest {
		gotAll, err = graph.ShortestAll(from, to)
	} else {
		gotAll, err = graph.LongestAll(from, to)
	}
	testErrors(t, wanterr, err, "From graph")
	if len(gotAll) == 0 {
		gotAll = BestPaths{BestPath{}}
	}
	testResultsGraphAll(t, gotAll, best, shortest)
}

func testResultsGraphAll(t *testing.T, got, best BestPaths, shortest bool) {
	distmethod := "Shortest"
	if !shortest {
		distmethod = "Longest"
	}
	if len(got) != len(best) {
		t.Error(distmethod, " amount of solutions incorrect\ngot: ", len(got), "\nwant: ", len(best))
		return
	}
	for i := range got {
		if got[i].Distance != best[i].Distance {
			t.Error(distmethod, " distance incorrect\ngot: ", got[i].Distance, "\nwant: ", best[i].Distance)
		}
	}
	for i := range got {
		found := false
		j := -1
		for j = range best {
			if reflect.DeepEqual(got[i].Path, best[j].Path) {
				//delete found result
				best = append(best[:j], best[j+1:]...)
				found = true
				break
			}
		}
		if found == false {
			t.Error(distmethod, " could not find path in solution\ngot:", got[i].Path)
		}
	}
}

func testResults(t *testing.T, got, best BestPath, shortest bool, filename string) {
	distmethod := "Shortest"
	if !shortest {
		distmethod = "Longest"
	}
	if got.Distance != best.Distance {
		t.Error(distmethod, " distance incorrect\n", filename, "\ngot: ", got.Distance, "\nwant: ", best.Distance)
	}
	if !reflect.DeepEqual(got.Path, best.Path) {
		t.Error(distmethod, " path incorrect\n\n", filename, "got: ", got.Path, "\nwant: ", best.Path)
	}
}

func getKSolLong() BestPath {
	return BestPath{
		31,
		[]int{
			0, 1, 3, 2, 4,
		},
	}
}
func getKSolShort() BestPath {
	return BestPath{
		2,
		[]int{
			0, 3, 4,
		},
	}
}
