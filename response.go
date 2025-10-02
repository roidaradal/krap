package krap

import "github.com/roidaradal/fn/ds"

type BulkCreate[T any] struct {
	Items   *ds.List[*T]
	Success int
	Fail    int
	Fails   []string
}

type BulkAction struct {
	Success int
	Fail    int
	Fails   []string
}
