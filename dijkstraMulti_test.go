package dijkstra

import (
	"fmt"
	"log"
	"math"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

//pq "github.com/Professorq/dijkstra"
func TestMultiLowThreads(t *testing.T) {
	g, _ := Import("testdata/B.txt")
	_, err := g.multiEvaluate(0, 1, 0, true)
	if err.Error() != "threads must be greater than 0" {
		t.Error("Threads error test failure")
	}
}
func TestMultiNoPath(t *testing.T) {
	testSolution(t, BestPath{}, ErrNoPath, "testdata/I.txt", 0, 4, true)
}

func TestMultiLoop(t *testing.T) {
	testSolution(t, BestPath{}, newErrLoop(1, 2), "testdata/J.txt", 0, 4, true)
}

func TestMultiCorrect(t *testing.T) {

	testMultiSolution(t, "testdata/B.txt", 0, 5, true)
	//testMultiSolution(t, "testdata/K.txt", 0, 4, false)
	testMultiSolution(t, "testdata/K.txt", 0, 4, true)

}

func testMultiSolution(t *testing.T, filename string, from, to int, shortest bool) {
	var err error
	var singleGraph Graph
	var distmethod string
	singleGraph, err = Import(filename)
	multiGraph, _ := Import(filename)
	if err != nil {
		t.Fatal(err, filename)
	}
	var got BestPath
	if shortest {
		distmethod = "Shortest"
		got, _ = singleGraph.Shortest(from, to)
	} else {
		distmethod = "Longest"
		got, _ = singleGraph.Longest(from, to)
	}
	//Run 1k times
	for j := 0; j < 100; j++ {
		multiGraph, _ = Import(filename)
		got2, _ := multiGraph.MultiShortest(from, to)
		checkMultiResults(t, got, got2, distmethod, filename, singleGraph, multiGraph)
		multiGraph, _ = Import(filename)
		//Run with different amounts of threads, 1,2,4....
		for i := 1; i <= math.MaxInt32; i *= 2 {
			got2, _ := multiGraph.multiEvaluate(from, to, i, shortest)
			checkMultiResults(t, got, got2, distmethod, filename, singleGraph, multiGraph)
		}
	}
}

func checkMultiResults(t *testing.T, single, multi BestPath, distmethod, filename string, singleGraph, multiGraph Graph) {
	if multi.Distance != single.Distance || !reflect.DeepEqual(multi.Path, single.Path) {
		t.Error(distmethod, " distance incorrect\n", filename, "\ngot: ", multi.Distance, "\nwant: ", single.Distance, " path incorrect\n\n", filename, "got: ", multi.Path, "\nwant: ", single.Path)
		c := spew.NewDefaultConfig()
		c.MaxDepth = 4
		c.Indent = "   "
		fmt.Println("\n", distmethod, " distance incorrect", filename, "got: ", multi.Distance, "want: ", single.Distance)
		fmt.Println(multi.Path, "==========", single.Path)

		fmt.Println("")
		for a := range multiGraph.Verticies {
			fmt.Println(a, ":", multiGraph.Verticies[a].distance, "===", a, ":", singleGraph.Verticies[a].distance)
		}
		fmt.Println("====", filename, "====", single.Path)
		//c.Dump(multiGraph)
		log.Fatal("exit")
	}
}
