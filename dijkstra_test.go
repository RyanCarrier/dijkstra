package dijkstra

import (
	"math"
	"reflect"
	"strconv"
	"testing"

	ar "github.com/albertorestifo/dijkstra"

	mm "github.com/mattomatic/dijkstra/dijkstra"
	mmg "github.com/mattomatic/dijkstra/graph"
)

//pq "github.com/Professorq/dijkstra"

func TestFailure(t *testing.T) {
	testSolution(t, BestPath{}, ErrNoPath, "testdata/I.txt", 0, 4)
}

func TestCorrect(t *testing.T) {
	testSolution(t, getBSol(), nil, "testdata/B.txt", 0, 5)
}
func BenchmarkRyanCarrierNodes4(b *testing.B)    { benchmarkAlt(b, "testdata/4.txt", 0) }
func BenchmarkRyanCarrierNodes10(b *testing.B)   { benchmarkAlt(b, "testdata/10.txt", 0) }
func BenchmarkRyanCarrierNodes100(b *testing.B)  { benchmarkAlt(b, "testdata/100.txt", 0) }
func BenchmarkRyanCarrierNodes1000(b *testing.B) { benchmarkAlt(b, "testdata/1000.txt", 0) }

/*
func BenchmarkProfessorQNodes4(b *testing.B)    { benchmarkAlt(b, "testdata/4.txt", 1) }
func BenchmarkProfessorQNodes10(b *testing.B)   { benchmarkAlt(b, "testdata/10.txt", 1) }
func BenchmarkProfessorQNodes100(b *testing.B)  { benchmarkAlt(b, "testdata/100.txt", 1) }
func BenchmarkProfessorQNodes1000(b *testing.B) { benchmarkAlt(b, "testdata/1000.txt", 1) }
*/

func BenchmarkAlbertoNodes4(b *testing.B)    { benchmarkAlt(b, "testdata/4.txt", 2) }
func BenchmarkAlbertoNodes10(b *testing.B)   { benchmarkAlt(b, "testdata/10.txt", 2) }
func BenchmarkAlbertoNodes100(b *testing.B)  { benchmarkAlt(b, "testdata/100.txt", 2) }
func BenchmarkAlbertoNodes1000(b *testing.B) { benchmarkAlt(b, "testdata/1000.txt", 2) }

func BenchmarkMattomaticNodes4(b *testing.B)    { benchmarkAlt(b, "testdata/4.txt", 3) }
func BenchmarkMattomaticNodes10(b *testing.B)   { benchmarkAlt(b, "testdata/10.txt", 3) }
func BenchmarkMattomaticNodes100(b *testing.B)  { benchmarkAlt(b, "testdata/100.txt", 3) }
func BenchmarkMattomaticNodes1000(b *testing.B) { benchmarkAlt(b, "testdata/1000.txt", 3) }

func benchmarkAlt(b *testing.B, filename string, i int) {
	switch i {
	case 0:
		benchmarkRC(b, filename)
		//case 1:
		//	benchmarkProfQ(b, filename)
	case 2:
		benchmarkAR(b, filename)
	case 3:
		benchmarkMM(b, filename)
	default:
		b.Error("You're retarded")
	}
}

func benchmarkMM(b *testing.B, filename string) {
	rcg, _, _ := Import(filename)
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
	rcg, _, _ := Import(filename)
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

/*
func benchmarkProfQ(b *testing.B, filename string) {
	var g *pq.Graph
	rcg, _, _ := Import(filename)
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
}*/

func benchmarkRC(b *testing.B, filename string) {
	graph, _, _ := Import(filename)
	src, dest := 0, len(graph.Verticies)-1
	//====RESET TIMER BEFORE LOOP====
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.SetDefaults(int64(math.MaxInt64), -1)
		graph.Shortest(src, dest)
	}
}

func setupAR(rcg Graph) ar.Graph {
	g := map[string]map[string]int{}
	for _, v := range rcg.Verticies {
		sv := strconv.Itoa(v.ID)
		g[sv] = map[string]int{}
		for key, val := range v.Arcs {
			g[sv][strconv.Itoa(key)] = int(val)
		}
	}
	return g
}

/*
func setupPq(rcg Graph) map[int]pq.Vertex {
	vs := map[int]pq.Vertex{}
	for _, v := range rcg.Verticies {
		temp := pq.Vertex{}
		temp.ID = v.ID
		temp.Arcs = map[int]int{}
		for key, val := range v.Arcs {
			temp.Arcs[key] = int(val)
		}
		vs[temp.ID] = temp
	}
	return vs
}*/

func testSolution(t *testing.T, best BestPath, wanterr error, filename string, from, to int) {
	graph, _, err := Import(filename)
	if err != nil {
		t.Error(err)
	}
	got, err := graph.Shortest(from, to)
	testErrors(t, wanterr, err)
	if got.Distance != best.Distance {
		t.Error("Distance incorrect\ngot: ", got.Distance, "\nwant: ", best.Distance)
	}
	if !reflect.DeepEqual(got.Path, best.Path) {
		t.Error("Path incorrect\ngot: ", got.Path, "\nwant: ", best.Path)
	}
}
