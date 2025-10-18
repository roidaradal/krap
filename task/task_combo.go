package task

import (
	"github.com/gin-gonic/gin"
	"github.com/roidaradal/krap"
	"github.com/roidaradal/krap/authz"
	"github.com/roidaradal/krap/root"
	"github.com/roidaradal/rdb/ze"
)

type comboConfig[A Actor, T any, P any] struct {
	*baseConfig[A, P]
	outputFn func(P, *T, *ze.Request, error)
}

type ComboTask[A Actor, T any] struct {
	*BaseTask[A]
	Fn               ComboFn[A, T]
	DeferActionCheck bool
}

// Create cmd comboConfig
func cmdComboConfig[A Actor, T any](task *ComboTask[A, T]) *comboConfig[A, T, []string] {
	cfg := &comboConfig[A, T, []string]{}
	cfg.initialize = task.cmdInitialize
	cfg.errorFn = cmdDisplayError
	cfg.outputFn = func(args []string, item *T, rq *ze.Request, err error) {
		krap.DisplayData(item, rq, err)
	}
	return cfg
}

// Create web comboConfig
func webComboConfig[A Actor, T any](task *ComboTask[A, T]) *comboConfig[A, T, *gin.Context] {
	cfg := &comboConfig[A, T, *gin.Context]{}
	cfg.initialize = task.webInitialize
	cfg.errorFn = krap.SendDataError
	cfg.outputFn = krap.SendDataResponse
	return cfg
}

// Creates new ComboTask
func NewComboTask[A Actor, T any](action, item string, fn ComboFn[A, T], deferActionCheck bool) *ComboTask[A, T] {
	task := &ComboTask[A, T]{}
	task.Action = action
	task.Item = item
	task.Fn = fn
	task.DeferActionCheck = deferActionCheck
	return task
}

// ComboTask CmdHandler
func (task ComboTask[A, T]) CmdHandler() root.CmdHandler {
	return comboTaskHandler(task, cmdComboConfig(&task))
}

// ComboTask WebHandler
func (task ComboTask[A, T]) WebHandler() gin.HandlerFunc {
	return comboTaskHandler(task, webComboConfig(&task))
}

// Common: create ComboTask Handler
func comboTaskHandler[A Actor, T any, P any](task ComboTask[A, T], cfg *comboConfig[A, T, P]) func(P) {
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
