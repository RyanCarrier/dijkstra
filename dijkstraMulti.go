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
	current.stageDebug = 1
	eval.visiting.Unlock()
	current.setActive(true)
	current.stageDebug = 2
	defer current.setActive(false)
	//don't have to lock cause writting never gets done to these areas
	current.RLock()
	defer current.RUnlock()
	current.stageDebug = 4
	//If the current distance is already worse than the best try another Vertex
	if (eval.shortest && current.distance >= eval.best) || (!eval.shortest && current.distance <= eval.best) {
		current.stageDebug = -1
		return
	}
	current.stageDebug = 5
	eval.checkArcs(current)
}

func (eval *evaluation) getNextVertex() *Vertex {
	if eval.shortest {
		return eval.visiting.popFront()
	}
	return eval.visiting.popBack()
}

func (eval *evaluation) checkArcs(current *Vertex) {
	current.stageDebug = 6
	for v, dist := range current.arcs {
		current.stageDebug = 7
		current.RUnlock()
		//	time.Sleep(time.Millisecond * 1)
		current.stageDebug = 8
		current.RLock()
		current.stageDebug = 9
		select {
		case <-current.quit:
		//	return
		default:
		}
		current.stageDebug = 10
		if v == current.ID {
			current.stageDebug = -2
			//could deadlock if arc to self lol
			continue
		}
		current.stageDebug = 10 + 10000*v
		eval.Verticies[v].RLock()
		current.stageDebug = 11
		if (eval.shortest && current.distance+dist < eval.Verticies[v].distance) ||
			(!eval.shortest && current.distance+dist > eval.Verticies[v].distance) {
			//Check for loop
			current.stageDebug = 12
			if eval.Verticies[v].active {
				current.stageDebug = 1211111
				select {
				case <-eval.Verticies[v].quit:
				default:
				}
				current.stageDebug = 13
				eval.Verticies[v].quit <- true
			}
			current.stageDebug = 14
			eval.Verticies[v].RUnlock()
			eval.Verticies[v].Lock()
			eval.Verticies[v].distance = current.distance + dist
			eval.Verticies[v].bestVertex = current.ID
			eval.Verticies[v].Unlock()
			eval.Verticies[v].RLock()
			current.stageDebug = 15
			if v == eval.dest {
				current.stageDebug = 16
				eval.best = current.distance + dist
				eval.visitedDest = true
			} else {
				current.stageDebug = 17
				eval.visiting.Lock()
				eval.visiting.pushOrdered(eval.Verticies[v])
				eval.visiting.Unlock()
			}
		}
		current.stageDebug = 18
		eval.Verticies[v].RUnlock()
	}
	current.stageDebug = -3
}
