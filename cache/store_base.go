package cache

import (
	"github.com/roidaradal/fn/dict"
)

type codeable interface {
	GetCode() string
	SetIsActive(bool)
}

// T is expected to be a reference type
type Store[T codeable] struct {
	isActive bool
	codeMap  *dict.SyncMap[string, T]
}

// Create new Store
func NewStore[T codeable]() *Store[T] {
	return &Store[T]{
		isActive: true,
		codeMap:  dict.NewSyncMap[string, T](),
	}
}

// Create new disabled Store
func NewDisabledStore[T codeable]() *Store[T] {
	return &Store[T]{
		isActive: false,
		codeMap:  nil,
	}
}

func (s *Store[T]) IsDisabled() bool {
	return !useCache || !s.isActive
}

// Gets all stored objects
func (s *Store[T]) All() []T {
	if s.IsDisabled() {
		return nil
	}
	return s.codeMap.Values()
}

// Get item by code
func (s *Store[T]) GetByCode(code string) (T, bool) {
	if s.IsDisabled() {
		var t T
		return t, false
	}
	return s.codeMap.Get(code)
}

// Add items to store
func (s *Store[T]) AddItems(items []T) {
	if s.IsDisabled() {
		return
	}
	for _, item := range items {
		s.Add(item)
	}
}

// Add item to store
func (s *Store[T]) Add(item T) {
	if s.IsDisabled() {
		return
	}
	s.codeMap.Set(item.GetCode(), item)
}

// Update item in store
func (s *Store[T]) Update(item T) {
	if s.IsDisabled() {
		return
	}
	s.codeMap.Set(item.GetCode(), item)
}

// Toggle item in store by code
func (s *Store[T]) ToggleByCode(code string, isActive bool) {
	if s.IsDisabled() {
		return
	}
	item, ok := s.GetByCode(code)
	if !ok {
		return
	}
	item.SetIsActive(isActive)
	s.Update(item)
}

// Delete item in store by code
func (s *Store[T]) DeleteByCode(code string) {
	if s.IsDisabled() {
		return
	}
	s.codeMap.Delete(code)
}
