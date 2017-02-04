package dijkstra

import "sync"

type semWG struct {
	sync.RWMutex
	threads int
}

func (s *semWG) lockUnlock() {
	s.RUnlock()
	s.RLock()
}

func (s *semWG) dec() {
	s.Lock()
	defer s.Unlock()
	s.threads--
}

func (s *semWG) incr() {
	s.Lock()
	defer s.Unlock()
	s.threads++
}
