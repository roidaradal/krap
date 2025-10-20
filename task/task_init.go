package task

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/roidaradal/fn"
	"github.com/roidaradal/fn/str"
	"github.com/roidaradal/krap/authn"
	"github.com/roidaradal/krap/authz"
	"github.com/roidaradal/rdb/ze"
)

const actionGlue string = "-"

var (
	ErrInvalidActor  = errors.New("public: Invalid actor")
	ErrInvalidOption = errors.New("public: Invalid option")
	errMissingHook   = errors.New("missing hook")
)

// Attach CmdDecorator to BaseTask
func (t *BaseTask[A]) WithCmd(cmdDecorator CmdDecorator[A]) {
	t.CmdDecorator = cmdDecorator
}

// Attach WebDecorator to BaseTask
func (t *BaseTask[A]) WithWeb(webDecorator WebDecorator[A]) {
	t.WebDecorator = webDecorator
}

// Attach CmdDecorator to BaseTokenTask
func (t *BaseTokenTask) WithCmd(cmdDecorator CmdTokenDecorator) {
	t.CmdDecorator = cmdDecorator
}

// Attach WebDecorator to BaseTokenTask
func (t *BaseTokenTask) WithWeb(webDecorator WebTokenDecorator) {
	t.WebDecorator = webDecorator
}

// Attach CmdDecorator to BaseDataTask
func (t *BaseDataTask[A]) WithCmd(cmdDecorator CmdDataDecorator[A]) {
	t.CmdDecorator = cmdDecorator
}

// Attach WebDecorator to BaseDataTask
func (t *BaseDataTask[A]) WithWeb(webDecorator WebDataDecorator[A]) {
	t.WebDecorator = webDecorator
}

// Attach CmdDecorator to BaseDataTokenTask
func (t *BaseDataTokenTask) WithCmd(cmdDecorator CmdDataTokenDecorator) {
	t.CmdDecorator = cmdDecorator
}

// Attach WebDecorator to BaseDataTokenTask
func (t *BaseDataTokenTask) WithWeb(webDecorator WebDataTokenDecorator) {
	t.WebDecorator = webDecorator
}

// Common BaseTask initialization process (cmd or web)
func initialize[A Actor](task BaseTask[A], params Params, actor *A, err error) (*ze.Request, Params, *A, error) {
	if err == nil && actor == nil {
		err = ErrInvalidActor
	}
	if err != nil {
		return nil, nil, nil, err
	}
	// Create request
	name := itemPrefix(task.Item)
	rq, err := ze.NewRequest(name)
	if err != nil {
		return rq, nil, nil, err
	}
	// Attach action, item to request
	rq.Action = task.Action
	rq.Item = task.Item
	return rq, params, actor, nil
}

// Common BaseTokenTask initialization process (cmd or web)
func initializeToken(task BaseTokenTask, params Params, authToken *authn.Token, err error) (*ze.Request, Params, *authn.Token, error) {
	if err == nil && authToken == nil {
		err = authn.ErrInvalidSession
	}
	if err != nil {
		return nil, nil, nil, err
	}
	// Create request
	name := itemPrefix(task.Item)
	rq, err := ze.NewRequest(name)
	if err != nil {
		return rq, nil, nil, err
	}
	// Attach action, item to request
	rq.Action = task.Action
	rq.Item = task.Item
	return rq, params, authToken, nil
}

// Common BaseDataTask initialize process (cmd or web)
func initializeData[A Actor](task BaseDataTask[A], params Params, actor *A, mustBeActive bool, err error) (*ze.Request, Params, *A, error) {
	if err == nil && actor == nil {
		err = ErrInvalidActor
	}
	if err != nil {
		return nil, nil, nil, err
	}
	// Create request
	name := itemPrefix(task.Item)
	rq, err := ze.NewRequest(name)
	if err != nil {
		return rq, nil, nil, err
	}
	// Attach action, item to request
	rq.Action = fn.Ternary(mustBeActive, authz.VIEW, authz.ROWS)
	rq.Item = task.Item
	return rq, params, actor, nil
}

// Common BaseDataTokenTask initialize process (cmd or web)
func initializeDataToken(task BaseDataTokenTask, params Params, authToken *authn.Token, mustBeActive bool, err error) (*ze.Request, Params, *authn.Token, error) {
	if err == nil && authToken == nil {
		err = authn.ErrInvalidSession
	}
	if err != nil {
		return nil, nil, nil, err
	}
	// Create request
	name := itemPrefix(task.Item)
	rq, err := ze.NewRequest(name)
	if err != nil {
		return rq, nil, nil, err
	}
	// Attach action, item to request
	rq.Action = fn.Ternary(mustBeActive, authz.VIEW, authz.ROWS)
	rq.Item = task.Item
	return rq, params, authToken, nil
}

// Initialize for BaseTask CmdHandler
func (task BaseTask[A]) cmdInitialize(args []string) (*ze.Request, Params, *A, error) {
	// Decorate the params
	params := make(Params)
	params, actor, err := task.CmdDecorator(args, params)
	return initialize(task, params, actor, err)
}

// Initialize for BaseTask WebHandler
func (task BaseTask[A]) webInitialize(c *gin.Context) (*ze.Request, Params, *A, error) {
	// Decorate the params
	params := make(Params)
	params, actor, err := task.WebDecorator(c, params)
	return initialize(task, params, actor, err)
}

// Initialize for BaseTokenTask CmdHandler
func (task BaseTokenTask) cmdInitialize(args []string) (*ze.Request, Params, *authn.Token, error) {
	// Decorate the params
	params := make(Params)
	params, authToken, err := task.CmdDecorator(args, params)
	return initializeToken(task, params, authToken, err)
}

// Initialize for BaseTokenTask WebHandler
func (task BaseTokenTask) webInitialize(c *gin.Context) (*ze.Request, Params, *authn.Token, error) {
	// Decorate the params
	params := make(Params)
	params, authToken, err := task.WebDecorator(c, params)
	return initializeToken(task, params, authToken, err)
}

// Initialize for BaseDataTask CmdHandler
func (task BaseDataTask[A]) cmdInitialize(args []string) (*ze.Request, Params, *A, error) {
	// Decorate the params
	params := make(Params)
	params, actor, mustBeActive, err := task.CmdDecorator(args, params)
	return initializeData(task, params, actor, mustBeActive, err)
}

// Initialize for BaseDataTask WebHandler
func (task BaseDataTask[A]) webInitialize(c *gin.Context) (*ze.Request, Params, *A, error) {
	// Decorate the params
	params := make(Params)
	params, actor, mustBeActive, err := task.WebDecorator(c, params)
	return initializeData(task, params, actor, mustBeActive, err)
}

// Initialize for BaseDataTokenTask CmdHandler
func (task BaseDataTokenTask) cmdInitialize(args []string) (*ze.Request, Params, *authn.Token, error) {
	// Decorate the params
	params := make(Params)
	params, authToken, mustBeActive, err := task.CmdDecorator(args, params)
	return initializeDataToken(task, params, authToken, mustBeActive, err)
}

// Initialize for BaseDataTokenTask WebHandler
func (task BaseDataTokenTask) webInitialize(c *gin.Context) (*ze.Request, Params, *authn.Token, error) {
	// Decorate the params
	params := make(Params)
	params, authToken, mustBeActive, err := task.WebDecorator(c, params)
	return initializeDataToken(task, params, authToken, mustBeActive, err)
}

// Get code from args string list on index
func getCode(args []string, index int) string {
	code := ""
	if index < len(args) {
		code = strings.ToUpper(args[index])
	}
	return code
}

// Gets the core item name, removes trailing "%s" if any
func itemPrefix(item string) string {
	if isCompleteItem(item) {
		return item
	}
	parts := str.CleanSplit(item, actionGlue)
	return strings.Join(parts[:len(parts)-1], actionGlue)
}

// Common: Checks if name ends in "%s"
func isCompleteItem(item string) bool {
	return !strings.HasPrefix(item, "%s")
}
