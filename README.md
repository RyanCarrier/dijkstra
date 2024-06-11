# dijkstra

Fast golang Dijkstra's shortest (and longest) path finder

## Documentation

[godoc](https://pkg.go.dev/github.com/RyanCarrier/dijkstra/v2)

## How to

### Install

```bash
go get github.com/RyanCarrier/dijkstra/v2
```

### Generate a graph

#### Importing from file

The package can import dijkstra files in the format:

```
0 1,1 2,1
1 0,1 2,2
2
```

using;

```go
data, _ := os.ReadFile("path/to/file")
graph, err := dijkstra.Import(string(data))
```

i.e. node then each arc and it's weight. The default is to use nodes with
numbers starting from 0, but the package will map string appropriately (using
ImportStringMapped instead.

#### Creating a graph

```go
package main

import "github.com/RyanCarrier/dijkstra/v2"

func main(){
  graph := dijkstra.NewGraph()
  //Add the 3 verticies
  graph.AddEmptyVertex(0)
  graph.AddEmptyVertex(1)
  graph.AddEmptyVertex(2)
  //Add the arcs
  graph.AddArc(0, 1, 1)
  graph.AddArc(0, 2, 1)
  graph.AddArc(1, 2, 2)
}

```

### Finding paths

Once the graph is created, shortest or longest paths between two points can be generated.

```go

  best, err := graph.Shortest(0, 2)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("Shortest distance is", best.Distance, "following path ", best.Path)

  best, err = graph.Longest(0, 2)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("Longest distance is", best.Distance, "following path ", best.Path)


```

### Finding multiple paths

```go
graph := dijkstra.NewGraph()
//Add the 3 verticies
graph.AddVertexAndArcs(0, map[int]uint64{1: 1, 2: 1})
graph.AddVertexAndArcs(1, map[int]uint64{3: 1})
graph.AddVertexAndArcs(2, map[int]uint64{3: 1})
graph.AddVertexAndArcs(3, map[int]uint64{4: 1})

best, err := graph.ShortestAll(0, 4)
if err != nil {
  fmt.Println(graph)
  log.Fatal(err)
}
fmt.Println("Shortest distances are", best.Distance, "with paths; ", best.Paths)

best, err = graph.LongestAll(0, 4)
if err != nil {
  log.Fatal(err)
}
fmt.Println("Longest distances are", best.Distance, "following path ", best.Paths)

```
