package task

import (
	"github.com/gin-gonic/gin"
	"github.com/roidaradal/fn/ds"
	"github.com/roidaradal/krap"
	"github.com/roidaradal/krap/authz"
	"github.com/roidaradal/krap/root"
	"github.com/roidaradal/rdb/ze"
)

type listConfig[T any, P any] struct {
	*baseTokenConfig[P]
	outputFn func(P, *ds.List[*T], *ze.Request, error)
}

type codedListConfig[A Actor, T any, P any] struct {
	*baseActorConfig[A, P]
	outputFn func(P, *ds.List[*T], *ze.Request, error)
}

type ListTask[T any] struct {
	*BaseDataTokenTask
	Fn ListFn[T]
}

type CodedListTask[A Actor, T any] struct {
	*BaseDataTask[A]
	Fn        ListFn[T]
	Validator HookFn[A]
	CodeIndex int
}

// Creates new ListTask
func NewListTask[T any](item string, fn ListFn[T]) *ListTask[T] {
	task := &ListTask[T]{}
	task.Item = item
	task.Fn = fn
	return task
}

// Creates new CodedListTask
func NewCodedListTask[A Actor, T any](item string, fn ListFn[T], codeIndex int) *CodedListTask[A, T] {
	task := &CodedListTask[A, T]{}
	task.Item = item
	task.Fn = fn
	task.CodeIndex = codeIndex
	return task
}

// Attach HookFn to CodedListTask
func (task *CodedListTask[A, T]) WithValidator(hookFn HookFn[A]) {
	task.Validator = hookFn
}

// ListTask CmdHandler
func (task ListTask[T]) CmdHandler() root.CmdHandler {
	cfg := &listConfig[T, []string]{}
	cfg.initialize = task.cmdInitialize
	cfg.errorFn = cmdDisplayError
	cfg.outputFn = func(args []string, list *ds.List[*T], rq *ze.Request, err error) {
		krap.DisplayList(list, rq, err)
	}
	return listTaskHandler(&task, cfg)
}

// ListTask WebHandler
func (task ListTask[T]) WebHandler() gin.HandlerFunc {
	cfg := &listConfig[T, *gin.Context]{}
	cfg.initialize = task.webInitialize
	cfg.errorFn = krap.SendDataError
	cfg.outputFn = krap.SendDataResponse
	return listTaskHandler(&task, cfg)
}

// Common: create ListTask Handler
func listTaskHandler[T any, P any](task *ListTask[T], cfg *listConfig[T, P]) func(P) {
	return func(p P) {
		// Initialize
		rq, params, authToken, err := cfg.initialize(p)
		if err != nil {
			cfg.errorFn(p, rq, err)
			return
		}
		// Check Authorization
		err = authz.CheckActionAllowedFor(rq, authToken.Type)
		var list *ds.List[*T]
		if err == nil {
			// Get data if authorized
			list, err = task.Fn(rq, params)
		}
		cfg.outputFn(p, list, rq, err)
	}
}

// CodedListTask CmdHandler
func (task CodedListTask[A, T]) CmdHandler() root.CmdHandler {
	cfg := &codedListConfig[A, T, []string]{}
	cfg.initialize = task.cmdInitialize
	cfg.errorFn = cmdDisplayError
	cfg.outputFn = func(args []string, list *ds.List[*T], rq *ze.Request, err error) {
		krap.DisplayList(list, rq, err)
	}
	codeFn := func(args []string) string {
		return getCode(args, task.CodeIndex)
	}
	return codedListTaskHandler(&task, cfg, codeFn)
}

// CodedListTask WebHandler
func (task CodedListTask[A, T]) WebHandler() gin.HandlerFunc {
	cfg := &codedListConfig[A, T, *gin.Context]{}
	cfg.initialize = task.webInitialize
	cfg.errorFn = krap.SendDataError
	cfg.outputFn = krap.SendDataResponse
	return codedListTaskHandler(&task, cfg, krap.WebCodeOption)
}

// Common: create CodedListTask Handler
func codedListTaskHandler[A Actor, T any, P any](task *CodedListTask[A, T], cfg *codedListConfig[A, T, P], codeFn func(P) string) func(P) {
	return func(p P) {
		// Initialize
		rq, params, actor, err := cfg.initialize(p)
		if err != nil {
			cfg.errorFn(p, rq, err)
			return
		}
		if task.Validator == nil {
			cfg.errorFn(p, rq, errMissingHook)
			return
		}
		// Get code and call validator
		code := codeFn(p)
		params, err = task.Validator(rq, params, actor, code)
		if err != nil {
			cfg.errorFn(p, rq, err)
			return
		}
		// Get list after passing all checks
		list, err := task.Fn(rq, params)
		cfg.outputFn(p, list, rq, err)
	}
}
