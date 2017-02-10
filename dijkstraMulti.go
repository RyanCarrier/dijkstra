package dijkstra

import (
	"errors"
	"fmt"
	"log"
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

//DOES NOT DETECT INF LOOPS
func (g *Graph) multiEvaluate(src, dest, threads int, shortest bool) (BestPath, error) {
	if threads < 1 {
		return BestPath{}, errors.New("threads must be greater than 0")
	}
	eval := g.multiSetup(src, dest, threads, shortest)
	timer := time.NewTimer(2 * time.Second)
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

func (eval *evaluation) multiLoop() {
	eval.RLock()
	for eval.threads > 0 || eval.getListLen() > 0 {
		for ; eval.threads > 0 && eval.getListLen() == 0; eval.lockUnlock() {
		}
		for eval.getListLen() > 0 {
			//Visit the current lowest distanced Vertex
			for ; eval.threads >= eval.maxThreads; eval.lockUnlock() {
			}
			eval.RUnlock()
			eval.visiting.Lock()
			if eval.visiting.len > 0 {
				eval.incr()
				go eval.multiVisitNode()
			}
			eval.visiting.Unlock()
			eval.RLock()
		}
	}
	eval.RUnlock()
}

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
	eval.visiting.Lock()
	if eval.visiting.len == 0 {
		eval.visiting.Unlock()
		return
	}
	current := eval.getNextVertex()
	eval.visiting.Unlock()
	current.setActive(true)
	defer current.setActive(false)
	//don't have to lock cause writting never gets done to these areas
	current.RLock()
	defer current.RUnlock()
	//If the current distance is already worse than the best try another Vertex
	if (eval.shortest && current.distance >= eval.best) || (!eval.shortest && current.distance <= eval.best) {
		return
	}
	eval.checkArcs(current)
}

func (eval *evaluation) getNextVertex() *Vertex {
	if eval.shortest {
		return eval.visiting.popFront()
	}
	return eval.visiting.popBack()
}

func (eval *evaluation) checkArcs(current *Vertex) {
	for v, dist := range current.arcs {
		select {
		case <-current.quit:
			return
		default:
		}
		if v == current.ID {
			//could deadlock if arc to self lol
			continue
		}
		eval.Verticies[v].RLock()
		if (eval.shortest && current.distance+dist < eval.Verticies[v].distance) ||
			(!eval.shortest && current.distance+dist > eval.Verticies[v].distance) {
			//Check for loop
			if eval.Verticies[v].active {
				//ensures only 1 sitting in the queue
				select {
				case <-eval.Verticies[v].quit:
				default:
				}
				eval.Verticies[v].quit <- true
			}
			eval.Verticies[v].swapToLock()
			eval.Verticies[v].distance = current.distance + dist
			eval.Verticies[v].bestVertex = current.ID
			eval.Verticies[v].swapToRLock()
			if v == eval.dest {
				eval.best = current.distance + dist
				eval.visitedDest = true
			} else {
				eval.visiting.Lock()
				eval.visiting.pushOrdered(eval.Verticies[v])
				eval.visiting.Unlock()
			}
		}
		eval.Verticies[v].RUnlock()
	}
}
