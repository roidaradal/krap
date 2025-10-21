package task

import (
	"github.com/gin-gonic/gin"
	"github.com/roidaradal/krap"
	"github.com/roidaradal/krap/authz"
	"github.com/roidaradal/krap/root"
	"github.com/roidaradal/rdb/ze"
)

type dataConfig[T any, P any] struct {
	*baseTokenConfig[P]
	outputFn func(P, *T, *ze.Request, error)
}

type codedDataConfig[A Actor, T any, P any] struct {
	*baseActorConfig[A, P]
	outputFn func(P, *T, *ze.Request, error)
}

type DataTask[T any] struct {
	*BaseDataTokenTask
	Fn DataFn[T]
}

type CodedDataTask[A Actor, T any] struct {
	*BaseDataTask[A]
	Fn        DataFn[T]
	Validator HookFn[A]
	CodeIndex int
}

// Creates new DataTask
func NewDataTask[T any](item string, fn DataFn[T]) *DataTask[T] {
	task := &DataTask[T]{}
	task.Item = item
	task.Fn = fn
	return task
}

// Creates new CodedDataTask
func NewCodedDataTask[A Actor, T any](item string, fn DataFn[T], codeIndex int) *CodedDataTask[A, T] {
	task := &CodedDataTask[A, T]{}
	task.Item = item
	task.Fn = fn
	task.CodeIndex = codeIndex
	return task
}

// Attach HookFn to CodedDataTask
func (task *CodedDataTask[A, T]) WithValidator(hookFn HookFn[A]) {
	task.Validator = hookFn
}

// DataTask CmdHandler
func (task DataTask[T]) CmdHandler() root.CmdHandler {
	cfg := &dataConfig[T, []string]{}
	cfg.initialize = task.cmdInitialize
	cfg.errorFn = cmdDisplayError
	cfg.outputFn = func(args []string, data *T, rq *ze.Request, err error) {
		krap.DisplayData(data, rq, err)
	}
	return dataTaskHandler(&task, cfg)
}

// DataTask WebHandler
func (task DataTask[T]) WebHandler() gin.HandlerFunc {
	cfg := &dataConfig[T, *gin.Context]{}
	cfg.initialize = task.webInitialize
	cfg.errorFn = krap.SendDataError
	cfg.outputFn = krap.SendDataResponse
	return dataTaskHandler(&task, cfg)
}

// Common: create ListTask Handler
func dataTaskHandler[T any, P any](task *DataTask[T], cfg *dataConfig[T, P]) func(P) {
	return func(p P) {
		// Initialize
		rq, params, authToken, err := cfg.initialize(p)
		if err != nil {
			cfg.errorFn(p, rq, err)
			return
		}
		// Check Authorization
		err = authz.CheckActionAllowedFor(rq, authToken.Type)
		var data *T
		if err == nil {
			// Get data if authorized
			data, err = task.Fn(rq, params)
		}
		cfg.outputFn(p, data, rq, err)
	}
}

// CodedListTask CmdHandler
func (task CodedDataTask[A, T]) CmdHandler() root.CmdHandler {
	cfg := &codedDataConfig[A, T, []string]{}
	cfg.initialize = task.cmdInitialize
	cfg.errorFn = cmdDisplayError
	cfg.outputFn = func(args []string, item *T, rq *ze.Request, err error) {
		krap.DisplayData(item, rq, err)
	}
	codeFn := func(args []string) string {
		return getCode(args, task.CodeIndex)
	}
	return codedDataTaskHandler(&task, cfg, codeFn)
}

// CodedListTask WebHandler
func (task CodedDataTask[A, T]) WebHandler() gin.HandlerFunc {
	cfg := &codedDataConfig[A, T, *gin.Context]{}
	cfg.initialize = task.webInitialize
	cfg.errorFn = krap.SendDataError
	cfg.outputFn = krap.SendDataResponse
	return codedDataTaskHandler(&task, cfg, krap.WebCodeOption)
}

// Common: create CodedListTask Handler
func codedDataTaskHandler[A Actor, T any, P any](task *CodedDataTask[A, T], cfg *codedDataConfig[A, T, P], codeFn func(P) string) func(P) {
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
		// Get data after passing all checks
		data, err := task.Fn(rq, params)
		cfg.outputFn(p, data, rq, err)
	}
}
