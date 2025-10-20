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

// Create cmd taskConfig
func cmdTaskConfig[A Actor, T any](task *FullTask[A, T]) *taskConfig[A, T, []string] {
	cfg := &taskConfig[A, T, []string]{}
	cfg.initialize = task.cmdInitialize
	cfg.errorFn = cmdDisplayError
	cfg.outputFn = func(args []string, item *T, rq *ze.Request, err error) {
		krap.DisplayData(item, rq, err)
	}
	return cfg
}

// Create web taskConfig
func webTaskConfig[A Actor, T any](task *FullTask[A, T]) *taskConfig[A, T, *gin.Context] {
	cfg := &taskConfig[A, T, *gin.Context]{}
	cfg.initialize = task.webInitialize
	cfg.errorFn = krap.SendDataError
	cfg.outputFn = krap.SendDataResponse
	return cfg
}

// Creates new FullTask
func NewFullTask[A Actor, T any](action, item string, fn TaskFn[A, T], deferActionCheck bool) *FullTask[A, T] {
	task := &FullTask[A, T]{}
	task.Action = action
	task.Item = item
	task.Fn = fn
	task.DeferActionCheck = deferActionCheck
	return task
}

// FullTask CmdHandler
func (task FullTask[A, T]) CmdHandler() root.CmdHandler {
	return fullTaskHandler(&task, cmdTaskConfig(&task))
}

// FullTask WebHandler
func (task FullTask[A, T]) WebHandler() gin.HandlerFunc {
	return fullTaskHandler(&task, webTaskConfig(&task))
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
