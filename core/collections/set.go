package collections

import (
	"sync"
)

type IntSet struct {
	m map[int64]bool
	sync.RWMutex
}

func NewIntSet() *IntSet {
	return &IntSet{
		m: map[int64]bool{},
	}
}

func (s *IntSet) Add(item int64) {
	s.Lock()
	defer s.Unlock()
	s.m[item] = true
}

func (s *IntSet) Remove(item int64) {
	s.Lock()
	s.Unlock()
	delete(s.m, item)
}

func (s *IntSet) Has(item int64) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[item]
	return ok
}

func (s *IntSet) Len() int {
	return len(s.List())
}

func (s *IntSet) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = map[int64]bool{}
}

func (s *IntSet) IsEmpty() bool {
	if s.Len() == 0 {
		return true
	}
	return false
}

func (s *IntSet) List() []int64 {
	s.RLock()
	defer s.RUnlock()
	list := []int64{}
	for item := range s.m {
		list = append(list, item)
	}
	return list
}
