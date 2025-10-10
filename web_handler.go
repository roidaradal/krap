package krap

import "github.com/gin-gonic/gin"

type EndpointHandlers map[string]gin.HandlerFunc

type WebHandler struct {
	Map  EndpointHandlers
	Verb string
}

// Register web handlers, returns the number of endpoints
func RegisterRoutes(server *gin.Engine, baseURL string, handlers []WebHandler) int {
	count := 0
	router := server.Group(baseURL)
	for _, handler := range handlers {
		for endpoint, handlerFn := range handler.Map {
			router.Handle(handler.Verb, endpoint, handlerFn)
			count += 1
		}
	}
	return count
}
