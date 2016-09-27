package collection

import (
	"os"
	"sync"
)

type Set struct {
	m map[interface{}]bool
	sync.RWMutex
}

type SortedFileArray []os.FileInfo

func (sf SortedFileArray) Reverse() {
	num := len(sf) / 2
	for i, j := 0, len(sf)-1; i < num; i, j = i+1, j-1 {
		sf[i], sf[j] = sf[j], sf[i]
	}
}

func (sf SortedFileArray) Len() int {
	return len(sf)
}

func (sf SortedFileArray) Swap(i, j int) {
	sf[i], sf[j] = sf[j], sf[i]
}

func (sf SortedFileArray) Less(i, j int) bool {
	return sf[i].Name() < sf[j].Name()
}

type SortedStringArray []string

func (sf SortedStringArray) Reverse() {
	num := len(sf) / 2
	for i, j := 0, len(sf)-1; i < num; i, j = i+1, j-1 {
		sf[i], sf[j] = sf[j], sf[i]
	}
}

func NewSet() *Set {
	m := make(map[interface{}]bool)
	return &Set{
		m: m,
	}
}

func newSet(m map[interface{}]bool) *Set {
	return &Set{
		m: m,
	}
}
func (s *Set) Copy() *Set {
	s.Lock()
	defer s.Unlock()

	m := make(map[interface{}]bool)
	for k, v := range s.m {
		m[k] = v
	}
	set := newSet(m)
	return set
}

func (s *Set) Add(item interface{}) {
	s.Lock()
	defer s.Unlock()
	s.m[item] = true
}
func (s *Set) Remove(item interface{}) {
	s.Lock()
	s.Unlock()
	delete(s.m, item)
}
func (s *Set) RemoveAll(set *Set) {
	s.Lock()
	s.Unlock()
	for _, e := range set.List() {
		delete(s.m, e)
	}
}

func (s *Set) Has(item interface{}) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[item]
	return ok
}
func (s *Set) Len() int {
	return len(s.List())
}

func (s *Set) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = map[interface{}]bool{}
}
func (s *Set) IsEmpty() bool {
	if s.Len() == 0 {
		return true
	}
	return false
}
func (s *Set) List() []interface{} {
	s.RLock()
	defer s.RUnlock()
	list := []interface{}{}
	for item := range s.m {
		list = append(list, item)
	}
	return list
}
