package conk

import (
	"context"
	"time"

	"github.com/roidaradal/rdb/ze"
	"golang.org/x/sync/errgroup"
)

type ActionFn = func() error
type RequestFn = func(*ze.Request) error

// Perform actions (func() error) concurrently
func Actions(actions []ActionFn) error {
	return ActionsWithTimeout(actions, 0)
}

// Perform actions (func() error) concurrently, with timeout
func ActionsWithTimeout(actions []ActionFn, timeoutSeconds uint) error {
	ctx := context.Background()
	if timeoutSeconds > 0 {
		var cancel context.CancelFunc
		duration := time.Duration(timeoutSeconds) * time.Second
		ctx, cancel = context.WithTimeout(ctx, duration)
		defer cancel()
	}

	group, ctx := errgroup.WithContext(ctx)
	for _, action := range actions {
		group.Go(action)
	}

	return group.Wait()
}

// Perform requests (func(*Request) error) concurrently
func Requests(rq *ze.Request, requests []RequestFn) error {
	return RequestsWithTimeout(rq, requests, 0)
}

// Perform requests (func(*Request) error) concurrently, with timeout
func RequestsWithTimeout(rq *ze.Request, requests []RequestFn, timeoutSeconds uint) error {
	ctx := context.Background()
	if timeoutSeconds > 0 {
		var cancel context.CancelFunc
		duration := time.Duration(timeoutSeconds) * time.Second
		ctx, cancel = context.WithTimeout(ctx, duration)
		defer cancel()
	}

	group, ctx := errgroup.WithContext(ctx)
	for _, request := range requests {
		group.Go(func() error {
			srq := rq.SubRequest()
			err := request(srq)
			rq.MergeLogs(srq)
			return err
		})
	}

	return group.Wait()
}
