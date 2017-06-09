package dijkstra

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/davecgh/go-spew/spew"
)

type evaluation struct {
	*Graph
	*semWG
	src        int
	dest       int
	shortest   bool
	maxThreads int
}

//MultiShortest evaluates shortest path using multiple gorountines (maxing at
// amount of nodes or GOMAXPROCS)
func (g *Graph) MultiShortest(src, dest int) (BestPath, error) {
	return g.multiEvaluate(src, dest, runtime.GOMAXPROCS(0), true)
}

//DOES NOT DETECT INF LOOPS jkz it should now
func (g *Graph) multiEvaluate(src, dest, threads int, shortest bool) (BestPath, error) {
	if threads < 1 {
		return BestPath{}, errors.New("threads must be greater than 0")
	}
	if threads == 1 {
		return g.evaluate(src, dest, shortest)
	}
	eval := g.multiSetup(src, dest, threads, shortest)
	timer := time.NewTimer(30 * time.Second)
	cancel := make(chan bool)
	go func() {
		for {
			select {
			case <-timer.C:
				eval.Lock()
				spew.Dump(eval)
				fmt.Println("==========================================================================")
				for i := 0; i < 10; i++ {
					fmt.Println("")
				}
				s := spew.NewDefaultConfig()
				s.MaxDepth = 2
				s.Indent = "\t"
				s.Dump(eval.Verticies)
				log.Fatal("Timeout")
			case <-cancel:
				timer.Stop()
				timer = nil
				return
			default:
			}
		}
	}()
	eval.multiLoop()
	cancel <- true
	return eval.finally(src, dest)
}

func (eval *evaluation) getThreads() int {
	eval.RLock()
	t := eval.threads
	eval.RUnlock()
	return t
}

func (eval *evaluation) multiLoop() {
	MT := eval.maxThreads
	for eval.getThreads() > 0 || eval.getListLen() > 0 {
		for eval.getThreads() > 0 && eval.getListLen() == 0 {
		}
		if eval.getListLen() > 0 {
			//Visit the current lowest distanced Vertex
			for eval.getThreads() >= MT {
			}
			eval.visiting.Lock()
			if eval.visiting.len > 0 {
				eval.incr()
				eval.visiting.Lock()
				go eval.multiVisitNode()
			}
			eval.visiting.Unlock()
			eval.RLock()
		}
	}
}

//multiSetup sets up the graph for evaluation
func (g *Graph) multiSetup(src, dest, threads int, shortest bool) *evaluation {
	wg := semWG{}
	nodes := len(g.Verticies)
	//Setup graph
	wg.Lock()
	wg.threads = 0
	wg.Unlock()
	g.setup(shortest, src)
	if threads < nodes-1 && threads > 0 {
		return &evaluation{g, &wg, src, dest, shortest, threads}
	}
	return &evaluation{g, &wg, src, dest, shortest, nodes - 1}
}

//DOES NOT DETECT INF LOOPS
func (eval *evaluation) multiVisitNode() {
	defer eval.dec()
	//Check if been removed
	if eval.visiting.len == 0 {
		eval.visiting.Unlock()
		return
	}
	current := eval.getNextVertex()
	eval.visiting.Unlock()
	//don't have to lock cause writting never gets done to these areas
	current.RLock()
	//If the current distance is already worse than the best try another Vertex
	if (eval.shortest && current.distance >= eval.best) || (!eval.shortest && current.distance <= eval.best) {
		current.RUnlock()
		return
	}
	current.RUnlock()
	eval.checkArcs(current)
}

func (eval *evaluation) getNextVertex() *Vertex {
	if eval.shortest {
		return eval.visiting.popFront()
	}
	return eval.visiting.popBack()
}

func (eval *evaluation) checkArcs(current *Vertex) {
	current.setActive(true)
	defer current.setActive(false)
	/*spew.Dump(current)
	spew.Dump(eval)
	log.Fatal("FUk u")*/
	current.RLock()
	defer current.RUnlock()
	//	var u update
	for v, dist := range current.arcs {

		if v == current.ID {
			//could deadlock if arc to self lol
			continue
		}
		eval.Verticies[v].RLock()
		if (eval.shortest && current.distance+dist < eval.Verticies[v].distance) ||
			(!eval.shortest && current.distance+dist > eval.Verticies[v].distance) {

			eval.Verticies[v].swapToLock()
			//make sure we are still right...
			if (eval.shortest && current.distance+dist < eval.Verticies[v].distance) ||
				(!eval.shortest && current.distance+dist > eval.Verticies[v].distance) {
				eval.Verticies[v].distance = current.distance + dist
				eval.Verticies[v].bestVertex = current.ID
			}
			eval.Verticies[v].swapToRLock()
			if v == eval.dest {
				eval.best = current.distance + dist
				eval.visitedDest = true
			} else {
				eval.visiting.Lock()
				eval.visiting.pushOrdered(eval.Verticies[v])
				eval.visiting.Unlock()
			}
			//		}
		}
		eval.Verticies[v].RUnlock()
	}
}
