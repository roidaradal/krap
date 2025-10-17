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
	ActionFn[A Actor]       = func(*ze.Request, Params, *A) error
	ComboFn[A Actor, T any] = func(*ze.Request, Params, *A) (*T, error)
	DataFn[T any]           = func(*ze.Request, Params) (*T, error)
	ListFn[T any]           = func(*ze.Request, Params) (*ds.List[*T], error)
)

type BaseTask[A Actor] struct {
	ze.Task
	CmdDecorator CmdDecorator[A]
	WebDecorator WebDecorator[A]
}

type BaseTaskToken struct {
	ze.Task
	CmdDecorator CmdDecoratorToken
	WebDecorator WebDecoratorToken
}

type (
	CmdDecoratorToken     = func([]string, Params) (Params, *authn.Token, error)
	WebDecoratorToken     = func(*gin.Context, Params) (Params, *authn.Token, error)
	CmdDecorator[A Actor] = func([]string, Params) (Params, *A, error)
	WebDecorator[A Actor] = func(*gin.Context, Params) (Params, *A, error)
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
