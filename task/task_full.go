package task

import (
	"github.com/gin-gonic/gin"
	"github.com/roidaradal/krap"
	"github.com/roidaradal/krap/authz"
	"github.com/roidaradal/krap/root"
	"github.com/roidaradal/rdb/ze"
)

type taskConfig[A Actor, T any, P any] struct {
	*baseActorConfig[A, P]
	outputFn func(P, *T, *ze.Request, error)
}

type FullTask[A Actor, T any] struct {
	*BaseTask[A]
	Fn               TaskFn[A, T]
	DeferActionCheck bool
}

type CodedFullTask[A Actor, T any] struct {
	*BaseTask[A]
	Fn        TaskFn[A, T]
	Validator HookFn[A]
	CodeIndex int
}

// Create cmd taskConfig
func cmdTaskConfig[A Actor, T any](task *BaseTask[A]) *taskConfig[A, T, []string] {
	cfg := &taskConfig[A, T, []string]{}
	cfg.initialize = task.cmdInitialize
	cfg.errorFn = cmdDisplayError
	cfg.outputFn = func(args []string, item *T, rq *ze.Request, err error) {
		krap.DisplayData(item, rq, err)
	}
	return cfg
}

// Create web taskConfig
func webTaskConfig[A Actor, T any](task *BaseTask[A]) *taskConfig[A, T, *gin.Context] {
	cfg := &taskConfig[A, T, *gin.Context]{}
	cfg.initialize = task.webInitialize
	cfg.errorFn = krap.SendDataError
	cfg.outputFn = krap.SendDataResponse
	return cfg
}

// Creates new FullTask
func NewFullTask[A Actor, T any](action, item string, fn TaskFn[A, T], deferActionCheck bool) *FullTask[A, T] {
	task := &FullTask[A, T]{
		BaseTask: &BaseTask[A]{},
	}
	task.Action = action
	task.Item = item
	task.Fn = fn
	task.DeferActionCheck = deferActionCheck
	return task
}

// Creates new CodedFullTask
func NewCodedFullTask[A Actor, T any](action, item string, fn TaskFn[A, T], codeIndex int) *CodedFullTask[A, T] {
	task := &CodedFullTask[A, T]{
		BaseTask: &BaseTask[A]{},
	}
	task.Action = action
	task.Item = item
	task.Fn = fn
	task.CodeIndex = codeIndex
	return task
}

// Attach HookFn to CodedFullTask
func (task *CodedFullTask[A, T]) WithValidator(hookFn HookFn[A]) {
	task.Validator = hookFn
}

// FullTask CmdHandler
func (task FullTask[A, T]) CmdHandler() root.CmdHandler {
	return fullTaskHandler(&task, cmdTaskConfig[A, T](task.BaseTask))
}

// FullTask WebHandler
func (task FullTask[A, T]) WebHandler() gin.HandlerFunc {
	return fullTaskHandler(&task, webTaskConfig[A, T](task.BaseTask))
}

// Common: create FullTask Handler
func fullTaskHandler[A Actor, T any, P any](task *FullTask[A, T], cfg *taskConfig[A, T, P]) func(P) {
	return func(p P) {
		// Initialize
		rq, params, actor, err := cfg.initialize(p)
		if err != nil {
			cfg.errorFn(p, rq, err)
			return
		}
		// Check Authorization if not deferred
		if !task.DeferActionCheck {
			err = authz.CheckActionAllowedFor(rq, (*actor).GetRole())
		}
		var item *T = nil
		if err == nil {
			// Perform action if authorized
			item, err = task.Fn(rq, params, actor)
		}
		cfg.outputFn(p, item, rq, err)
	}
}

// CodedFullTask CmdHandler
func (task CodedFullTask[A, T]) CmdHandler() root.CmdHandler {
	codeFn := func(args []string) string {
		return getCode(args, task.CodeIndex)
	}
	return codedFullTaskHandler(&task, cmdTaskConfig[A, T](task.BaseTask), codeFn)
}

// CodedFullTask WebHandler
func (task CodedFullTask[A, T]) WebHandler() gin.HandlerFunc {
	return codedFullTaskHandler(&task, webTaskConfig[A, T](task.BaseTask), krap.WebCodeParam)
}

// Common: create CodedFullTask Handler
func codedFullTaskHandler[A Actor, T any, P any](task *CodedFullTask[A, T], cfg *taskConfig[A, T, P], codeFn func(P) string) func(P) {
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
		// Perform action
		item, err := task.Fn(rq, params, actor)
		cfg.outputFn(p, item, rq, err)
	}
}
