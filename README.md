%# dijkstra
Golangs fastest Dijkstra's shortest (and longest) path calculator, requires go 1.6 or above (for benchmarking).

## Need for speed
Benchmark comparisons to the other two top golang dijkstra implementations;

```go test -bench .```

![Wow so fast](/speed.png?raw=true "Benchmarks")
![Wow the multiply!](/mult.png?raw=true "Multiply")

Speed benefit use to diminish due to linked lists sucking at high nodes in queue for checking. Since adding a priority queue (or linked list for small nodes), the benefits get even stronger.

Priority queues are used for nodes over 800, to gain this big increase in spead (about x15), it means there is some added overhead to initializing. Oddly this shouldn't be an issue at 16-256 nodes, but both still seem to take a significant performance hit.

## Documentation
[godoc](https://godoc.org/github.com/RyanCarrier/dijkstra)

## How to
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
graph, err := dijkstra.Import("path/to/file")
```

i.e. node then each arc and it's weight. The default is to use nodes with numbers starting from 0, but the package will map string appropriately.

#### Creating a graph

```go
package main

func main(){
  graph:=dijkstra.NewGraph()
  //Add the 3 verticies
  graph.AddVertex(0)
  graph.AddVertex(1)
  graph.AddVertex(2)
  //Add the arcs
  graph.AddArc(0,1,1)
  graph.AddArc(0,2,1)
  graph.AddArc(1,0,1)
  graph.AddArc(1,2,2)
}

```

### Finding paths

Once the graph is created, shortest or longest paths between two points can be generated.
```go

best, err := graph.Shortest(0,2)
if err!=nil{
  log.Fatal(err)
}
fmt.Println("Shortest distance ", best.Distance, " following path ", best.Path)

best, err := graph.Longest(0,2)
if err!=nil{
  log.Fatal(err)
}
fmt.Println("Longest distance ", best.Distance, " following path ", best.Path)

```
