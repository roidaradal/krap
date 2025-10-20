package task

import (
	"github.com/gin-gonic/gin"
	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/fn/ds"
	"github.com/roidaradal/krap/authn"
	"github.com/roidaradal/krap/root"
	"github.com/roidaradal/rdb/ze"
)

// Note: A type is for Actor

type Params = dict.Object

type Actor interface {
	GetRole() string
}

type (
	TaskFn[A Actor, T any] = func(*ze.Request, Params, *A) (*T, error)
	ActionFn[A Actor]      = func(*ze.Request, Params, *A) error
	DataFn[T any]          = func(*ze.Request, Params) (*T, error)
	ListFn[T any]          = func(*ze.Request, Params) (*ds.List[*T], error)
)

type BaseTask[A Actor] struct {
	ze.Task
	CmdDecorator CmdDecorator[A]
	WebDecorator WebDecorator[A]
}

type BaseDataTask struct {
	ze.Task
	CmdDecorator CmdDataDecorator
	WebDecorator WebDataDecorator
}

type (
	CmdDecorator[A Actor] = func([]string, Params) (Params, *A, error)
	WebDecorator[A Actor] = func(*gin.Context, Params) (Params, *A, error)
	CmdDataDecorator      = func([]string, Params) (Params, *authn.Token, bool, error)
	WebDataDecorator      = func(*gin.Context, Params) (Params, *authn.Token, bool, error)
)

type CmdHandler interface {
	CmdHandler() root.CmdHandler
}

type WebHandler interface {
	WebHandler() gin.HandlerFunc
}

type Handler interface {
	CmdHandler
	WebHandler
}

type (
	Router    = map[string]Handler
	AddRouter = map[[2]string]Handler
)

// Request, Params, Actor, Code, ID
type HookFn[A Actor] = func(*ze.Request, Params, *A, ze.Task, string, ze.ID) (Params, error)
