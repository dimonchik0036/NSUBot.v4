package main

import "sync"

type Set struct {
	Mux sync.RWMutex    `json:"-"`
	Set map[string]bool `json:"set"`
}

func (s *Set) Add(item string) {
	s.Mux.Lock()
	defer s.Mux.Unlock()
	if s.Set == nil {
		s.Set = map[string]bool{}
	}
	s.Set[item] = true
}

func (s *Set) Del(item string) {
	s.Mux.Lock()
	defer s.Mux.Unlock()
	delete(s.Set, item)
}

func (s *Set) Check(item string) bool {
	s.Mux.RLock()
	defer s.Mux.RUnlock()
	return s.Set[item]
}

func (s *Set) Change(item string) {
	s.Mux.Lock()
	defer s.Mux.Unlock()
	if s.Set == nil {
		s.Set = map[string]bool{}
		s.Set[item] = true
		return
	}

	if s.Set[item] {
		delete(s.Set, item)
	} else {
		s.Set[item] = true
	}
}

func (s *Set) GetAll() (result []string) {
	s.Mux.RLock()
	defer s.Mux.RUnlock()

	for k := range s.Set {
		result = append(result, k)
	}

	return
}
