package task

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/roidaradal/krap"
	"github.com/roidaradal/krap/root"
)

// Create new CmdConfig using task.CmdHandler()
func Cmd[T CmdHandler](command string, minParams int, docs string, task T) *root.CmdConfig {
	return &root.CmdConfig{
		Command:   command,
		MinParams: minParams,
		Docs:      docs,
		Handler:   task.CmdHandler(),
	}
}

// Create new CmdConfig Router
func CmdRoute[T CmdHandler](command string, minParams int, docs string, router map[string]T) *root.CmdConfig {
	// Build the handlers of each router option
	handlerOf := make(map[string]root.CmdHandler)
	for key, task := range router {
		handlerOf[key] = task.CmdHandler()
	}
	routerHandler := func(args []string) {
		option := strings.ToLower(args[0])
		handler, ok := handlerOf[option]
		if !ok {
			krap.DisplayError(ErrInvalidOption)
			return
		}
		handler(args)
	}
	return &root.CmdConfig{
		Command:   command,
		MinParams: minParams,
		Docs:      docs,
		Handler:   routerHandler,
	}
}

// Create gin.HandlerFunc from task.WebHandler()
func Web[T WebHandler](task T) gin.HandlerFunc {
	return task.WebHandler()
}

// Create gin.HandlerFunc Router
func Fork[T WebHandler](router map[string]T, response *krap.ResponseType) gin.HandlerFunc {
	// Build handlers of each router option
	handlerOf := make(map[string]gin.HandlerFunc)
	for key, task := range router {
		handlerOf[key] = task.WebHandler()
	}
	return func(c *gin.Context) {
		option := krap.WebForkParam(c)
		handler, ok := handlerOf[option]
		if !ok {
			response.SendErrorFn(c, nil, ErrInvalidOption)
			return
		}
		handler(c)
	}
}
