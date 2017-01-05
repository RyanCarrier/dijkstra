package bench

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"

	"github.com/RyanCarrier/dijkstra"
)

//Generate generates file with the amount of nodes specified
func Generate(nodes int) {
	filename := "bench/" + strconv.Itoa(nodes) + ".txt"
	fmt.Println("Generating file "+filename+" with nodes ", nodes)
	graph := dijkstra.Graph{}
	var i int
	for i = 0; i < nodes; i++ {
		v := dijkstra.NewVertex(i)
		for j := 0; j <= nodes; j++ {
			if j == i {
				continue
			}
			v.AddArc(j, int64(2*nodes-j)+rand.Int63n(int64(nodes)*int64(nodes-j+1)))
		}
		graph.AddVerticies(*v)
	}
	err := graph.ExportToFile(filename)
	if err != nil {
		log.Fatal(err)
	}
}
