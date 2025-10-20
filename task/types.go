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

type BaseTokenTask struct {
	ze.Task
	CmdDecorator CmdTokenDecorator
	WebDecorator WebTokenDecorator
}

type BaseDataTask[A Actor] struct {
	ze.Task
	CmdDecorator CmdDataDecorator[A]
	WebDecorator WebDataDecorator[A]
}

type BaseDataTokenTask struct {
	ze.Task
	CmdDecorator CmdDataTokenDecorator
	WebDecorator WebDataTokenDecorator
}

type (
	CmdDecorator[A Actor]     = func([]string, Params) (Params, *A, error)
	WebDecorator[A Actor]     = func(*gin.Context, Params) (Params, *A, error)
	CmdTokenDecorator         = func([]string, Params) (Params, *authn.Token, error)
	WebTokenDecorator         = func(*gin.Context, Params) (Params, *authn.Token, error)
	CmdDataDecorator[A Actor] = func([]string, Params) (Params, *A, bool, error)
	WebDataDecorator[A Actor] = func(*gin.Context, Params) (Params, *A, bool, error)
	CmdDataTokenDecorator     = func([]string, Params) (Params, *authn.Token, bool, error)
	WebDataTokenDecorator     = func(*gin.Context, Params) (Params, *authn.Token, bool, error)
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
type HookFn[A Actor] = func(*ze.Request, Params, *A, string) (Params, error)

// Request, Params, Actor, Schema, Code, ID
type TypedHookFn[A Actor, T any] = func(*ze.Request, Params, *A, *ze.Schema[T], string) (Params, error)
