package krap

import "github.com/gin-gonic/gin"

type CmdConfig struct {
	Command   string
	MinParams int
	Params    string
	Handler   CmdHandler
}

type EndpointHandlers map[string]gin.HandlerFunc

type WebHandler struct {
	Map  EndpointHandlers
	Verb string
}

func NewCommand(command string, minParams int, params string, handler CmdHandler) *CmdConfig {
	return &CmdConfig{command, minParams, params, handler}
}

func NewCommandMap(cfgs ...*CmdConfig) map[string]*CmdConfig {
	commands := make(map[string]*CmdConfig)
	for _, cfg := range cfgs {
		commands[cfg.Command] = cfg
	}
	return commands
}

func RegisterRoutes(server *gin.Engine, cfg *WebConfig, handlers []WebHandler) int {
	count := 0
	router := server.Group(cfg.Base)
	for _, handler := range handlers {
		for endpoint, handlerFn := range handler.Map {
			router.Handle(handler.Verb, endpoint, handlerFn)
			count += 1
		}
	}
	return count
}
