package dijkstra

import (
	"errors"
	"time"
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
	g.visiting.Unlock()
	current.setActive(true)
	defer current.setActive(false)
	//don't have to lock cause writting never gets done to these areas
	current.RLock()
	defer current.RUnlock()
	//If we have hit the destination set the flag, cheaper than checking it's
	// distance change at the end
	if current.ID == dest {
		return
	}
	//If the current distance is already worse than the best try another Vertex
	if (shortest && current.distance >= g.best) || (!shortest && current.distance <= g.best) {
		return
	}
	for v, dist := range current.arcs {
		current.RUnlock()
		time.Sleep(time.Millisecond * 1)
		current.RLock()
		select {
		case <-current.quit:
			return
		default:
		}
		if v == current.ID {
			//could deadlock if arc to self lol
			continue
		}
		g.Verticies[v].Lock()
		if (shortest && current.distance+dist < g.Verticies[v].distance) ||
			(!shortest && current.distance+dist > g.Verticies[v].distance) {
			//Check for loop
			if g.Verticies[v].active {
				g.Verticies[v].quit <- true
			}
			g.Verticies[v].distance = current.distance + dist
			g.Verticies[v].bestVertex = current.ID
			if v == dest {
				g.best = current.distance + dist
				g.visitedDest = true
			} else {
				g.visiting.Lock()
				g.visiting.pushOrdered(g.Verticies[v])
				g.visiting.Unlock()
			}
		}
		g.Verticies[v].Unlock()
	}
}
