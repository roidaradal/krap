package task

import (
	"github.com/gin-gonic/gin"
	"github.com/roidaradal/krap"
	"github.com/roidaradal/krap/authz"
	"github.com/roidaradal/krap/root"
	"github.com/roidaradal/rdb/ze"
)

type actionConfig[A Actor, P any] struct {
	*baseConfig[A, P]
	outputFn func(P, *ze.Request, error)
}

type ActionTask[A Actor] struct {
	*BaseTask[A]
	Fn ActionFn[A]
}

type CodedActionTask[A Actor] struct {
	*ActionTask[A]
	Validator HookFn[A]
	CodeIndex int
}

// Create cmd actionConfig
func cmdActionConfig[A Actor](task *ActionTask[A]) *actionConfig[A, []string] {
	cfg := &actionConfig[A, []string]{}
	cfg.initialize = task.cmdInitialize
	cfg.errorFn = cmdDisplayError
	cfg.outputFn = func(args []string, rq *ze.Request, err error) {
		krap.DisplayOutput(rq, err)
	}
	return cfg
}

// Create web actionConfig
func webActionConfig[A Actor](task *ActionTask[A]) *actionConfig[A, *gin.Context] {
	cfg := &actionConfig[A, *gin.Context]{}
	cfg.initialize = task.webInitialize
	cfg.errorFn = krap.SendActionError
	cfg.outputFn = krap.SendActionResponse
	return cfg
}

// Creates new ActionTask
func NewActionTask[A Actor](action, item string, fn ActionFn[A]) *ActionTask[A] {
	task := &ActionTask[A]{}
	task.Action = action
	task.Item = item
	task.Fn = fn
	return task
}

// Creates new CodedActionTask
func NewCodedActionTask[A Actor](action, item string, fn ActionFn[A], codeIndex int) *CodedActionTask[A] {
	task := &CodedActionTask[A]{}
	task.Action = action
	task.Item = item
	task.Fn = fn
	task.CodeIndex = codeIndex
	return task
}

// Attach HookFn to CodedActonTask
func (task *CodedActionTask[A]) WithValidator(hookFn HookFn[A]) {
	task.Validator = hookFn
}

// ActionTask CmdHandler
func (task ActionTask[A]) CmdHandler() root.CmdHandler {
	return actionTaskHandler(&task, cmdActionConfig(&task))
}

// ActionTask WebHandler
func (task ActionTask[A]) WebHandler() gin.HandlerFunc {
	return actionTaskHandler(&task, webActionConfig(&task))
}

// Common: create ActionTask Handler
func actionTaskHandler[A Actor, P any](task *ActionTask[A], cfg *actionConfig[A, P]) func(P) {
	return func(p P) {
		// Initialize
		rq, params, actor, err := cfg.initialize(p)
		if err != nil {
			cfg.errorFn(p, rq, err)
			return
		}
		// Check Authorization
		err = authz.CheckActionAllowedFor(rq, (*actor).GetRole())
		if err == nil {
			// Perform action if authorized
			err = task.Fn(rq, params, actor)
		}
		cfg.outputFn(p, rq, err)
	}
}

// CodedActionTask CmdHandler
func (task CodedActionTask[A]) CmdHandler() root.CmdHandler {
	codeFn := func(args []string) string {
		return getCode(args, task.CodeIndex)
	}
	return codedActionTaskHandler(&task, cmdActionConfig(task.ActionTask), codeFn)
}

// CodedActionTask WebHandler
func (task CodedActionTask[A]) WebHandler() gin.HandlerFunc {
	return codedActionTaskHandler(&task, webActionConfig(task.ActionTask), krap.WebCodeParam)
}

// Common: create CodedActionTask Handler
func codedActionTaskHandler[A Actor, P any](task *CodedActionTask[A], cfg *actionConfig[A, P], codeFn func(P) string) func(P) {
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
		params, err = task.Validator(rq, params, actor, task.Task, code, 0)
		if err != nil {
			cfg.errorFn(p, rq, err)
			return
		}
		// Perform action
		err = task.Fn(rq, params, actor)
		cfg.outputFn(p, rq, err)
	}
}
