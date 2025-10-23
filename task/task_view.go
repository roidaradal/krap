package task

import (
	"github.com/gin-gonic/gin"
	"github.com/roidaradal/krap"
	"github.com/roidaradal/krap/authz"
	"github.com/roidaradal/krap/root"
	"github.com/roidaradal/rdb/ze"
)

type viewConfig[T any, P any] struct {
	*baseTokenConfig[P]
	outputFn func(P, *T, *ze.Request, error)
}

type codedViewConfig[A Actor, T any, P any] struct {
	*baseActorConfig[A, P]
	outputFn func(P, *T, *ze.Request, error)
}

type ViewTask[T any] struct {
	*BaseTokenTask
	Fn DataFn[T]
}

type CodedViewTask[A Actor, T any] struct {
	*BaseTask[A]
	Fn        DataFn[T]
	Validator HookFn[A]
	CodeIndex int
}

// Creates new ViewTask
func NewViewTask[T any](item string, fn DataFn[T]) *ViewTask[T] {
	task := &ViewTask[T]{
		BaseTokenTask: &BaseTokenTask{},
	}
	task.Action = authz.VIEW
	task.Item = item
	task.Fn = fn
	return task
}

// Creates new CodedViewTask
func NewCodedViewTask[A Actor, T any](item string, fn DataFn[T], codeIndex int) *CodedViewTask[A, T] {
	task := &CodedViewTask[A, T]{
		BaseTask: &BaseTask[A]{},
	}
	task.Action = authz.VIEW
	task.Item = item
	task.Fn = fn
	task.CodeIndex = codeIndex
	return task
}

// Attach HookFn to CodedViewTask
func (task *CodedViewTask[A, T]) WithValidator(hookFn HookFn[A]) {
	task.Validator = hookFn
}

// ViewTask CmdHandler
func (task ViewTask[T]) CmdHandler() root.CmdHandler {
	cfg := &viewConfig[T, []string]{
		baseTokenConfig: &baseTokenConfig[[]string]{},
	}
	cfg.initialize = task.cmdInitialize
	cfg.errorFn = cmdDisplayError
	cfg.outputFn = func(args []string, item *T, rq *ze.Request, err error) {
		krap.DisplayData(item, rq, err)
	}
	return viewTaskHandler(&task, cfg)
}

// ViewTask WebHandler
func (task ViewTask[T]) WebHandler() gin.HandlerFunc {
	cfg := &viewConfig[T, *gin.Context]{
		baseTokenConfig: &baseTokenConfig[*gin.Context]{},
	}
	cfg.initialize = task.webInitialize
	cfg.errorFn = krap.SendDataError
	cfg.outputFn = krap.SendDataResponse
	return viewTaskHandler(&task, cfg)
}

// Common: create ViewTask Handler
func viewTaskHandler[T any, P any](task *ViewTask[T], cfg *viewConfig[T, P]) func(P) {
	return func(p P) {
		// Initialize
		rq, authToken, err := cfg.initialize(p)
		if err != nil {
			cfg.errorFn(p, rq, err)
			return
		}
		// Check Authorization
		err = authz.CheckActionAllowedFor(rq, authToken.Type)
		var item *T
		if err == nil {
			// View if authorized
			item, err = task.Fn(rq)
		}
		cfg.outputFn(p, item, rq, err)
	}
}

// CodedViewTask CmdHandler
func (task CodedViewTask[A, T]) CmdHandler() root.CmdHandler {
	cfg := &codedViewConfig[A, T, []string]{
		baseActorConfig: &baseActorConfig[A, []string]{},
	}
	cfg.initialize = task.cmdInitialize
	cfg.errorFn = cmdDisplayError
	cfg.outputFn = func(args []string, item *T, rq *ze.Request, err error) {
		krap.DisplayData(item, rq, err)
	}
	codeFn := func(args []string) string {
		return getCode(args, task.CodeIndex)
	}
	return codedViewTaskHandler(&task, cfg, codeFn)
}

// CodedViewTask WebHandler
func (task CodedViewTask[A, T]) WebHandler() gin.HandlerFunc {
	cfg := &codedViewConfig[A, T, *gin.Context]{
		baseActorConfig: &baseActorConfig[A, *gin.Context]{},
	}
	cfg.initialize = task.webInitialize
	cfg.errorFn = krap.SendDataError
	cfg.outputFn = krap.SendDataResponse
	return codedViewTaskHandler(&task, cfg, krap.WebCodeParam)
}

// Common: create CodedViewTask Handler
func codedViewTaskHandler[A Actor, T any, P any](task *CodedViewTask[A, T], cfg *codedViewConfig[A, T, P], codeFn func(P) string) func(P) {
	return func(p P) {
		// Initialize
		rq, actor, err := cfg.initialize(p)
		if err != nil {
			cfg.errorFn(p, rq, err)
			return
		}
		// Check validator, if it exists
		if task.Validator != nil {
			code := codeFn(p)
			err = task.Validator(rq, actor, code)
			if err != nil {
				cfg.errorFn(p, rq, err)
				return
			}
		}
		// Get view
		item, err := task.Fn(rq)
		cfg.outputFn(p, item, rq, err)
	}
}
