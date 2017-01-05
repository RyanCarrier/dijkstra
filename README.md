# dijkstra
Golangs fastest Dijkstra's shortest (and longest) path calculator, requires go 1.6 or above (for benchmarking).

## Need for speed
Benchmark comparisons to the other two top golang dijkstra implementations;

```go test -bench .```

![Wow so fast](/speed.png?raw=true "Benchmarks")
![Wow the multiply!](/mult.png?raw=true "Multiply")

Oddly the speed benefit seems to diminish as the nodes get higher, this could be
 due to the fact that the linked list places the to be visited node in the correct
 ordered position. This means worst case O(n). The generated test cases are worst
 case, with every need having access to every other node (weighting on the distances
   to ensure that the shortest path is through lots of nodes). This garuntees
  worst case scenario for my implementation, but maybe not for the others.

## Documentation
[godoc](https://godoc.org/github.com/RyanCarrier/dijkstra)

## How to
### Generate a graph
#### Importing from file

The package can import dijkstra files in the format;
```
0 1,1 2,1
1 0,1 2,2
2
```

using;
```go
graph, err := dijkstra.Import("path/to/file")
```

ie; node then each arc and it's weight. The default is to use nodes with numbers starting from 0, but the package will map string appropriatly.

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
