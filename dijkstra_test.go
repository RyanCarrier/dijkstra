package dijkstra

import (
	"os"
	"reflect"
	"strconv"
	"testing"

	ar "github.com/albertorestifo/dijkstra"

	pq "github.com/RyanCarrier/dijkstra-1"
	mm "github.com/mattomatic/dijkstra/dijkstra"
	mmg "github.com/mattomatic/dijkstra/graph"

	"github.com/RyanCarrier/dijkstra/testdata/bench"
)

//pq "github.com/Professorq/dijkstra"

func TestNoPath(t *testing.T) {
	testSolution(t, BestPath{}, ErrNoPath, "testdata/I.txt", 0, 4, true)
}

func TestLoop(t *testing.T) {
	testSolution(t, BestPath{}, newErrLoop(1, 2), "testdata/J.txt", 0, 4, true)
}

func TestCorrect(t *testing.T) {
	testSolution(t, getBSol(), nil, "testdata/B.txt", 0, 5, true)
	testSolution(t, getKSolLong(), nil, "testdata/K.txt", 0, 4, false)
	testSolution(t, getKSolShort(), nil, "testdata/K.txt", 0, 4, true)
}
func BenchmarkRyanCarrierNodes4(b *testing.B)    { benchmarkAlt(b, 4, 0) }
func BenchmarkRyanCarrierNodes10(b *testing.B)   { benchmarkAlt(b, 10, 0) }
func BenchmarkRyanCarrierNodes100(b *testing.B)  { benchmarkAlt(b, 100, 0) }
func BenchmarkRyanCarrierNodes1000(b *testing.B) { benchmarkAlt(b, 1000, 0) }

func BenchmarkProfessorQNodes4(b *testing.B)    { benchmarkAlt(b, 4, 1) }
func BenchmarkProfessorQNodes10(b *testing.B)   { benchmarkAlt(b, 10, 1) }
func BenchmarkProfessorQNodes100(b *testing.B)  { benchmarkAlt(b, 100, 1) }
func BenchmarkProfessorQNodes1000(b *testing.B) { benchmarkAlt(b, 1000, 1) }

func BenchmarkAlbertoNodes4(b *testing.B)    { benchmarkAlt(b, 4, 2) }
func BenchmarkAlbertoNodes10(b *testing.B)   { benchmarkAlt(b, 10, 2) }
func BenchmarkAlbertoNodes100(b *testing.B)  { benchmarkAlt(b, 100, 2) }
func BenchmarkAlbertoNodes1000(b *testing.B) { benchmarkAlt(b, 1000, 2) }

/*
//Mattomatics does not work.
func BenchmarkMattomaticNodes4(b *testing.B)    { benchmarkAlt(b, 4, 3) }
func BenchmarkMattomaticNodes10(b *testing.B)   { benchmarkAlt(b, 10, 3) }
func BenchmarkMattomaticNodes100(b *testing.B)  { benchmarkAlt(b, 100, 3) }
func BenchmarkMattomaticNodes1000(b *testing.B) { benchmarkAlt(b, 1000, 3) }
*/
func benchmarkAlt(b *testing.B, nodes, i int) {
	filename := "bench/" + strconv.Itoa(nodes) + ".txt"
	if err := os.Stat(filename); err != nil {
		bench.Generate(nodes)
	}
	switch i {
	case 0:
		benchmarkRC(b, filename)
	case 1:
		benchmarkProfQ(b, filename)
	case 2:
		benchmarkAR(b, filename)
	case 3:
		benchmarkMM(b, filename)
	default:
		b.Error("You're retarded")
	}
}

func benchmarkMM(b *testing.B, filename string) {
	rcg, _ := Import(filename)
	mg := mmg.LoadGraph(filename)
	rcsrc, rcdest := 0, len(rcg.Verticies)-1
	_, dest := 0, mg.Size()-1
	rcgot, _ := rcg.Shortest(rcsrc, rcdest)
	destN := mg.Search(dest)
	mgot := mm.Dijkstra(mg, destN)
	total := 0
	for _, d := range mgot {
		total += d
	}
	if rcgot.Distance != int64(total) {
		b.Fatal("Distances do not match, RC:", rcgot.Distance, " MM:", total)
	}
	rcgot, _ = rcg.Shortest(rcsrc, rcdest)
	total = 0
	mgot = mm.Dijkstra(mg, destN)
	for _, d := range mgot {
		total += d
	}
	if rcgot.Distance != int64(total) {
		b.Fatal("Distances do not match on iteration 2, RC:", rcgot.Distance, " MM:", total)
	}
	//====RESET TIMER BEFORE LOOP====
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mm.Dijkstra(mg, destN)
	}
}

func benchmarkAR(b *testing.B, filename string) {
	rcg, _ := Import(filename)
	arg := setupAR(rcg)
	rcsrc, rcdest := 0, len(rcg.Verticies)-1
	src, dest := "0", strconv.Itoa(rcdest)
	rcgot, _ := rcg.Shortest(rcsrc, rcdest)
	_, argot, _ := arg.Path("0", dest)
	if rcgot.Distance != int64(argot) {
		b.Fatal("Distances do not match, RC:", rcgot.Distance, " AR:", argot)
	}
	rcgot, _ = rcg.Shortest(rcsrc, rcdest)
	_, argot, _ = arg.Path("0", dest)
	if rcgot.Distance != int64(argot) {
		b.Fatal("Distances do not match on iteration 2, RC:", rcgot.Distance, " AR:", argot)
	}
	//====RESET TIMER BEFORE LOOP====
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		arg.Path(src, dest)
	}
}

func benchmarkProfQ(b *testing.B, filename string) {
	var g *pq.Graph
	rcg, _ := Import(filename)
	pqmap := setupPq(rcg)
	g = pq.NewGraph(pqmap)
	src, dest := 0, g.Len()-1
	rcsrc, rcdest := 0, len(rcg.Verticies)-1
	rcgot, _ := rcg.Shortest(rcsrc, rcdest)
	pqgot := g.ShortestPath(src, dest)
	if rcgot.Distance != int64(pqgot) {
		b.Fatal("Distances do not match, RC:", rcgot.Distance, " PQ:", pqgot)
	}
	rcgot, _ = rcg.Shortest(rcsrc, rcdest)
	pqgot = pq.NewGraph(pqmap).ShortestPath(src, dest)
	if rcgot.Distance != int64(pqgot) {
		b.Fatal("Distances do not match on iteration 2, RC:", rcgot.Distance, " PQ:", pqgot)
	}
	//====RESET TIMER BEFORE LOOP====
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pq.NewGraph(pqmap).ShortestPath(src, dest)
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

func setupAR(rcg Graph) ar.Graph {
	g := map[string]map[string]int{}
	for _, v := range rcg.Verticies {
		sv := strconv.Itoa(v.ID)
		g[sv] = map[string]int{}
		for key, val := range v.arcs {
			g[sv][strconv.Itoa(key)] = int(val)
		}
	}
	return g
}

func setupPq(rcg Graph) map[int]pq.Vertex {
	vs := map[int]pq.Vertex{}
	for _, v := range rcg.Verticies {
		temp := pq.Vertex{}
		temp.ID = v.ID
		temp.Arcs = map[int]int{}
		for key, val := range v.arcs {
			temp.Arcs[key] = int(val)
		}
		vs[temp.ID] = temp
	}
	return vs
}

func testSolution(t *testing.T, best BestPath, wanterr error, filename string, from, to int, shortest bool) {
	var err error
	var graph Graph
	graph, err = Import(filename)
	if err != nil {
		t.Fatal(err, filename)
	}
	var got BestPath
	if shortest {
		got, err = graph.Shortest(from, to)
	} else {
		got, err = graph.Longest(from, to)
	}
	testErrors(t, wanterr, err, filename)
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
