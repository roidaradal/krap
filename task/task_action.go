package task

import (
	"github.com/gin-gonic/gin"
	"github.com/roidaradal/krap"
	"github.com/roidaradal/krap/authz"
	"github.com/roidaradal/krap/root"
)

type ActionTask[A Actor] struct {
	BaseTask[A]
	Fn ActionFn[A]
}

type CodedActionTask[A Actor] struct {
	ActionTask[A]
	Validator HookFn[A]
	CodeIndex int
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

// Attach CmdDecorator to ActionTask, Return task to be chainable
func (t *ActionTask[A]) WithCmd(cmdDecorator CmdDecorator[A]) *ActionTask[A] {
	t.CmdDecorator = cmdDecorator
	return t
}

// Attach WebDecorator to ActionTask, Return task to be chainable
func (t *ActionTask[A]) WithWeb(webDecorator WebDecorator[A]) *ActionTask[A] {
	t.WebDecorator = webDecorator
	return t
}

// Attach CmdDecorator to CodedActionTask, Return task to be chainable
func (t *CodedActionTask[A]) WithCmd(cmdDecorator CmdDecorator[A]) *CodedActionTask[A] {
	t.CmdDecorator = cmdDecorator
	return t
}

// Attach WebDecorator to CodedActionTask, Return task to be chainable
func (t *CodedActionTask[A]) WithWeb(webDecorator WebDecorator[A]) *CodedActionTask[A] {
	t.WebDecorator = webDecorator
	return t
}

// Attach HookFn to CodedActonTask, Return task to be chainable
func (t *CodedActionTask[A]) WithValidator(hookFn HookFn[A]) *CodedActionTask[A] {
	t.Validator = hookFn
	return t
}

// ActionTask CmdHandler
func (task ActionTask[A]) CmdHandler() root.CmdHandler {
	return func(args []string) {
		// Initialize
		rq, params, actor, err := task.cmdInitialize(args)
		if err != nil {
			krap.DisplayError(err)
			return
		}
		// Check Authorization
		err = authz.CheckActionAllowedFor(rq, (*actor).GetRole())
		if err == nil {
			// Perform action if authorized
			err = task.Fn(rq, params, actor)
		}
		krap.DisplayOutput(rq, err)
	}
}

// ActionTask WebHandler
func (task ActionTask[A]) WebHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Initialize
		rq, params, actor, err := task.webInitialize(c)
		if err != nil {
			krap.SendActionError(c, rq, err)
			return
		}
		// Check Authorization
		err = authz.CheckActionAllowedFor(rq, (*actor).GetRole())
		if err == nil {
			// Perform action if authorized
			err = task.Fn(rq, params, actor)
		}
		krap.SendActionResponse(c, rq, err)
	}
}

// CodedActionTask CmdHandler
func (task CodedActionTask[A]) CmdHandler() root.CmdHandler {
	return func(args []string) {
		// Initialize
		rq, params, actor, err := task.cmdInitialize(args)
		if err != nil {
			krap.DisplayError(err)
			return
		}
		if task.Validator == nil {
			krap.DisplayError(errMissingHook)
			return
		}
		// Get code and call hook
		code := getCode(args, task.CodeIndex)
		params, err = task.Validator(rq, params, actor, task.Task, code, 0)
		if err != nil {
			krap.DisplayError(err)
			return
		}
		// Perform action
		err = task.Fn(rq, params, actor)
		krap.DisplayOutput(rq, err)
	}
}

// CodedActionTask WebHandler
func (task CodedActionTask[A]) WebHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Initialize
		rq, params, actor, err := task.webInitialize(c)
		if err != nil {
			krap.SendActionError(c, rq, err)
			return
		}
		if task.Validator == nil {
			krap.SendActionError(c, rq, errMissingHook)
			return
		}
		// Get code and call hook
		// Import: Code needs to be WebCodeParam
		code := krap.WebCodeParam(c)
		params, err = task.Validator(rq, params, actor, task.Task, code, 0)
		if err != nil {
			krap.SendActionError(c, rq, err)
			return
		}
		// Perform action
		err = task.Fn(rq, params, actor)
		krap.SendActionResponse(c, rq, err)
	}
}
