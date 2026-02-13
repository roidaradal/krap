package web

import "github.com/gin-gonic/gin"

type EndpointHandlers map[string]gin.HandlerFunc

type Handler struct {
	Map  EndpointHandlers
	Verb string
}

// Register web handlers, returns the number of endpoints
func RegisterRoutes(server *gin.Engine, baseURL string, handlers []Handler, middlewares ...gin.HandlerFunc) int {
	count := 0
	router := server.Group(baseURL)
	if len(middlewares) > 0 {
		router.Use(middlewares...)
	}
	for _, handler := range handlers {
		for endpoint, handlerFn := range handler.Map {
			router.Handle(handler.Verb, endpoint, handlerFn)
			count += 1
		}
	}
	return count
}
