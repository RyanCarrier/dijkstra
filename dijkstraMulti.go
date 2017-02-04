package dijkstra

import (
	"errors"
	"fmt"
)

//DOES NOT DETECT INF LOOPS
func (g *Graph) multiEvaluate(src, dest, threads int, shortest bool) (BestPath, error) {
	if threads < 1 {
		return BestPath{}, errors.New("threads must be greater than 0")
	}
	wg := semWG{}
	var maxThreads int
	nodes := len(g.Verticies)
	if threads < nodes-1 && threads > 0 {
		maxThreads = threads
	} else {
		//Limit threads to nodes
		maxThreads = nodes - 1
	}
	//Setup graph
	wg.Lock()
	wg.threads = 0
	wg.Unlock()
	g.setup(shortest, src)
	wg.RLock()
	for wg.threads > 0 || g.getListLen() > 0 {
		for ; wg.threads > 0 && g.getListLen() == 0; wg.lockUnlock() {
		}
		for g.getListLen() > 0 {
			//Visit the current lowest distanced Vertex
			for ; wg.threads >= maxThreads; wg.lockUnlock() {
			}
			wg.RUnlock()
			g.visiting.Lock()
			if g.visiting.len > 0 {
				wg.incr()
				go g.multiVisitNode(dest, shortest, &wg)
			}
			g.visiting.Unlock()
			wg.RLock()
		}
	}
	return g.finally(src, dest)
}

//DOES NOT DETECT INF LOOPS
func (g *Graph) multiVisitNode(dest int, shortest bool, wg *semWG) {

	defer wg.dec()
	var current *Vertex
	g.visiting.Lock()
	if g.visiting.len == 0 {
		g.visiting.Unlock()
		return
	}
	if shortest {
		current = g.visiting.popFront()
	} else {
		current = g.visiting.popBack()
	}
	//spew.Dump(current)
	g.visiting.Unlock()
	current.setActive(true)
	defer current.setActive(false)
	fmt.Println("0")
	//don't have to lock cause writting never gets done to these areas
	cdist, cid := current.distance, current.ID
	//If we have hit the destination set the flag, cheaper than checking it's
	// distance change at the end
	if cid == dest {
		return
	}
	fmt.Println("1")
	//If the current distance is already worse than the best try another Vertex
	if (shortest && cdist >= g.best) || (!shortest && cdist <= g.best) {
		return
	}
	fmt.Println("2")
	for v, dist := range current.arcs {
		fmt.Println("3")
		select {
		case q := <-current.quit:
			fmt.Println(q)
			fmt.Println("dying")
			return
		default:
			fmt.Println("Loopin")
		}
		fmt.Println("4")
		if v == cid {
			//could deadlock if arc to self lol
			continue
		}
		fmt.Println("5")
		//Implement RWMutex instead
		g.Verticies[v].Lock()
		fmt.Println("6")
		if (shortest && cdist+dist < g.Verticies[v].distance) ||
			(!shortest && cdist+dist > g.Verticies[v].distance) {
			fmt.Println("7")
			//Check for loop
			if g.Verticies[v].active {
				//kill
				fmt.Println("KILL ", v)
				select {
				case <-g.Verticies[v].quit:
					fmt.Println("HE COULD QUIT")
				default:
					fmt.Println("HE CANT QUIT")
				}
				g.Verticies[v].quit <- true
				fmt.Println("KILLIN ", v)
			}
			fmt.Println("8")
			g.Verticies[v].distance = cdist + dist
			g.Verticies[v].bestVertex = cid
			fmt.Println("9")
			if v == dest {
				g.best = cdist + dist
				g.visitedDest = true
			} else {
				g.visiting.Lock()
				g.visiting.pushOrdered(g.Verticies[v])
				g.visiting.Unlock()
			}
			/*g.Verticies[v].Unlock()
			} else {
				g.Verticies[v].Unlock()
			*/
		}
		g.Verticies[v].Unlock()
	}
}
