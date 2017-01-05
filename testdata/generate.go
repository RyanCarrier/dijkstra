package main

import (
	"flag"
	"log"
	"strconv"

	"github.com/RyanCarrier/dijkstra/testdata/bench"
)

var n int

func init() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatal("Pls specify amount of nodes")
	}
	var err error
	n, err = strconv.Atoi(flag.Args()[0])
	if err != nil {
		log.Fatal("Pls use numbers not words")
	}
}

func main() {
	bench.Generate(n)
}
