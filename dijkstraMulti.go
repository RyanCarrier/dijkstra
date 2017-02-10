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
	} /*
		retry := false
		for a := range eval.Verticies {
			select {
			case u := <-eval.Verticies[a].quit:
				fmt.Fprintln(os.Stderr, "WASTED UPDATEEE=========================================================")
				if (eval.shortest && u.newBestDist < eval.Verticies[a].distance) ||
					(!eval.shortest && u.newBestDist > eval.Verticies[a].distance) {
					eval.Verticies[a].Lock()
					eval.Verticies[a].bestVertex = u.newBestVertex
					eval.Verticies[a].distance = u.newBestDist
					eval.Verticies[a].Unlock()
					eval.visiting.pushOrdered(eval.Verticies[a])
					retry = true
				}
			default:
				//fmt.Fprintln(os.Stderr, "NO WASTED UPDATE", a)
			}
			if eval.Verticies[a].active {
				retry = true
				fmt.Fprintln(os.Stderr, "WTF I'm NOT DONE YET!!+!+!+!+!+!+")
			}
		}

		eval.RUnlock()
		if retry {
			eval.multiLoop()
		}*/
	//	time.Sleep(time.Second * 2)
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
		/*
			select {
			case u = <-current.quit:
				fmt.Fprintln(os.Stderr, "===================GOT UPDATE============", fmt.Sprintf("%+v", u))
				if (eval.shortest && u.newBestDist < current.distance) ||
					(!eval.shortest && u.newBestDist > current.distance) {
					current.swapToLock()
					current.bestVertex = u.newBestVertex
					current.distance = u.newBestDist
					current.swapToRLock()
					eval.checkArcs(current)
					return
				}
			default:
			}*/
		if v == current.ID {
			//could deadlock if arc to self lol
			continue
		}
		eval.Verticies[v].RLock()
		if (eval.shortest && current.distance+dist < eval.Verticies[v].distance) ||
			(!eval.shortest && current.distance+dist > eval.Verticies[v].distance) {
			//Check for loop
			//if eval.Verticies[v].active {

			//ensures only 1 sitting in the queue
			/*
				select {
				case u = <-eval.Verticies[v].quit:
					fmt.Fprintln(os.Stderr, "===================INTERCEPTED UPDATE============", fmt.Sprintf("%+v", u))
					if (eval.shortest && u.newBestDist > current.distance+dist) ||
						(!eval.shortest && u.newBestDist < current.distance+dist) {
						u = update{current.ID, current.distance + dist}
					}
				default:
					u = update{current.ID, current.distance + dist}
				}
				fmt.Fprintln(os.Stderr, "SENDING UPDATE", fmt.Sprintf("%+v", u), " to ", v)
				//go func(v int, u update) {
				eval.Verticies[v].quit <- u
				//}(v, u)*/
			//		} //else {
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
