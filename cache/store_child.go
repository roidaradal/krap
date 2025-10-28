package cache

import (
	"slices"

	"github.com/roidaradal/fn"
	"github.com/roidaradal/rdb/ze"
)

type child interface {
	GetParentID() ze.ID
}

type codeableChild interface {
	codeable
	child
}

type idcodeableChild interface {
	idcodeable
	child
}

// T is expected to be a reference type
type ChildStore[T codeableChild] struct {
	*Store[T]
}

// T is expected to be a reference type
type ChildIDStore[T idcodeableChild] struct {
	*IDStore[T]
}

// Create new ChildStore
func NewChildStore[T codeableChild]() *ChildStore[T] {
	return &ChildStore[T]{
		Store: NewStore[T](),
	}
}

// Create new disabled ChildStore
func NewDisabledChildStore[T codeableChild]() *ChildStore[T] {
	return &ChildStore[T]{
		Store: NewDisabledStore[T](),
	}
}

// Create new ChildIDStore
func NewChildIDStore[T idcodeableChild]() *ChildIDStore[T] {
	return &ChildIDStore[T]{
		IDStore: NewIDStore[T](),
	}
}

// Create new disabled ChildIDStore
func NewDisabledChildIDStore[T idcodeableChild]() *ChildIDStore[T] {
	return &ChildIDStore[T]{
		IDStore: NewDisabledIDStore[T](),
	}
}

// Get items with parent IDs
func (s *ChildStore[T]) FromParentIDs(parentIDs ...ze.ID) []T {
	if s.isDisabled() || len(parentIDs) == 0 {
		return nil
	}
	items := fn.Filter(s.All(), func(item T) bool {
		return slices.Contains(parentIDs, item.GetParentID())
	})
	return items
}

// Get items with parent IDs
func (s *ChildIDStore[T]) FromParentIDs(parentIDs ...ze.ID) []T {
	if s.isDisabled() || len(parentIDs) == 0 {
		return nil
	}
	items := fn.Filter(s.All(), func(item T) bool {
		return slices.Contains(parentIDs, item.GetParentID())
	})
	return items
}
