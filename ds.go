package krap

import "github.com/roidaradal/fn/ds"

type Identifiable interface {
	idable
	codeable
}

type idable interface {
	GetID() uint
}

type codeable interface {
	GetCode() string
}

type (
	IDCodeLookup           map[uint]string
	IDLookup[T idable]     map[uint]T
	CodeLookup[T codeable] map[string]T
	StringListMap          map[string][]string
)

type List[T any] struct {
	Items []T
	Count int
}

type ListLookup[T codeable, L any] struct {
	Items  []T
	Lookup map[string]L
	Count  int
}

type MapLookup[T codeable, L any] struct {
	Items  map[string]T
	Lookup map[string]L
	Count  int
}

func NewList[T any](items []T) *List[T] {
	if items == nil {
		items = []T{}
	}
	return &List[T]{Items: items, Count: len(items)}
}

func NewListLookup[T codeable, L any](items []T) *ListLookup[T, L] {
	if items == nil {
		items = []T{}
	}
	return &ListLookup[T, L]{
		Items:  items,
		Lookup: make(map[string]L),
		Count:  len(items),
	}
}

func NewMapLookup[T codeable, L any](items map[string]T) *MapLookup[T, L] {
	if items == nil {
		items = make(map[string]T)
	}
	return &MapLookup[T, L]{
		Items:  items,
		Lookup: make(map[string]L),
		Count:  len(items),
	}
}

func NewIDLookup[T idable](items []T) IDLookup[T] {
	idLookup := make(IDLookup[T])
	for _, item := range items {
		idLookup[item.GetID()] = item
	}
	return idLookup
}

func NewCodeLookup[T codeable](items []T) CodeLookup[T] {
	codeLookup := make(CodeLookup[T])
	for _, item := range items {
		codeLookup[item.GetCode()] = item
	}
	return codeLookup
}

func IDToCodeLookup[T Identifiable](idLookup IDLookup[T], validCodes *ds.Set[string]) CodeLookup[T] {
	codeLookup := make(CodeLookup[T])
	for _, item := range idLookup {
		code := item.GetCode()
		if validCodes != nil && !validCodes.Contains(code) {
			continue
		}
		codeLookup[code] = item
	}
	return codeLookup
}
