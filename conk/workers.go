package conk

import "github.com/roidaradal/rdb/ze"

type WorkerFn[T any] = func(T) error
type RequestWorkerFn[T any] = func(*ze.Request, T) error
