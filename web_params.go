package krap

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// Gets the web request origin
func WebRequestOrigin(c *gin.Context) *RequestOrigin {
	browserInfo := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()
	return &RequestOrigin{
		BrowserInfo: &browserInfo,
		IPAddress:   &ipAddress,
	}
}

// Gets the lowercase :Fork param
func WebForkParam(c *gin.Context) string {
	return strings.ToLower(c.Param("Fork"))
}

// Gets the uppercase :Code param
func WebCodeParam(c *gin.Context) string {
	return strings.ToUpper(c.Param("Code"))
}

// Gets the lowercase :Type param
func WebTypeParam(c *gin.Context) string {
	return strings.ToLower(c.Param("Type"))
}

// False if option is 'all', otherwise true,
// From ?view=option query string
func WebMustBeActiveOption(c *gin.Context) bool {
	option := c.DefaultQuery("view", viewActive)
	return MustBeActiveOption(option)
}

// Return toggle on/off (boolean), hasToggleOption (ok flag),
// From ?toggle=option query string
func WebToggleOption(c *gin.Context) (bool, bool) {
	option := c.Query("toggle")
	return ToggleOption(option)
}

// Gets the uppercase code option from ?code=option query string
func WebCodeOption(c *gin.Context) string {
	return strings.ToUpper(c.Query("code"))
}

// Gets the uppercase type option from ?type=option query string,
// Defaults to ANY_TYPE (*)
func WebTypeOption(c *gin.Context) string {
	typ := c.DefaultQuery("type", ANY_TYPE)
	return strings.ToUpper(typ)
}
