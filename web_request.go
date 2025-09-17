package krap

import (
	"github.com/gin-gonic/gin"
)

const authHeaderKey string = "Authorization"

func WebAuthHeader(c *gin.Context) *AuthToken {
	authHeader := c.GetHeader(authHeaderKey)
	return NewAuthToken(authHeader)
}

func WebRequestOrigin(c *gin.Context) *RequestOrigin {
	browserInfo := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()
	return &RequestOrigin{
		BrowserInfo: &browserInfo,
		IPAddress:   &ipAddress,
	}
}
