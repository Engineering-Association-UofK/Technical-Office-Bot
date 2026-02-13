package utils

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

var errKeyExists = errors.New("key already exists in map")
var errMissingKey = errors.New("key Does not exist in map")

type ConcurrentMap[K comparable, V any] struct {
	state map[K]V
	mut   sync.RWMutex
}

// Constructor
func NewConcurrentMap[K comparable, V any]() *ConcurrentMap[K, V] {
	return &ConcurrentMap[K, V]{
		state: make(map[K]V),
	}
}

// CRUD
func (cMap *ConcurrentMap[K, V]) Add(key K, value V) error {
	cMap.mut.Lock()
	defer cMap.mut.Unlock()

	if _, exists := cMap.state[key]; exists {
		return errKeyExists
	}
	cMap.state[key] = value
	return nil
}

func (cMap *ConcurrentMap[K, V]) Update(key K, value V) error {
	cMap.mut.Lock()
	defer cMap.mut.Unlock()

	if _, exists := cMap.state[key]; exists {
		cMap.state[key] = value
		return nil
	}
	return errMissingKey
}

func (cMap *ConcurrentMap[K, V]) Delete(key K) error {
	cMap.mut.Lock()
	defer cMap.mut.Unlock()

	if _, exists := cMap.state[key]; exists {
		delete(cMap.state, key)
		return nil
	}
	return errMissingKey
}

func (cMap *ConcurrentMap[K, V]) Value(key K) (V, error) {
	cMap.mut.RLock()
	defer cMap.mut.RUnlock()

	value, ok := cMap.state[key]
	if ok {
		return value, nil
	}
	return value, errMissingKey
}

func (cMap *ConcurrentMap[K, V]) Len() int {
	cMap.mut.RLock()
	defer cMap.mut.RUnlock()
	l := len(cMap.state)
	return l
}

// Other
func (cMap *ConcurrentMap[K, V]) Empty() {
	cMap.mut.Lock()
	defer cMap.mut.Unlock()
	clear(cMap.state)

}

func (cMap *ConcurrentMap[K, V]) String() string {
	var sb strings.Builder
	sb.WriteString("{\n")
	for k, v := range cMap.state {
		s := fmt.Sprintf("   \"%v\": \"%v\",\n", k, v)
		sb.WriteString(s)
	}
	sb.WriteString("}")

	return sb.String()
}
