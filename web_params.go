package krap

import (
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ANY_TYPE   string = "*"
	TOGGLE_ON  string = "on"
	TOGGLE_OFF string = "off"
	listAll    string = "all"
	listActive string = "active"
)

func WebCodeParam(c *gin.Context) string {
	code := c.Param("Code")
	return strings.ToUpper(code)
}

func WebTypeParam(c *gin.Context) string {
	typ := c.Param("Type")
	return strings.ToLower(typ)
}

func WebMustListBeActive(c *gin.Context) bool {
	option := c.DefaultQuery("list", listActive)
	return MustListBeActive(option)
}

// on/off, hasToggleOption
func WebToggleOption(c *gin.Context) (bool, bool) {
	option := c.DefaultQuery("toggle", "")
	return ToggleOption(option)
}

func WebCodeOption(c *gin.Context) string {
	code := c.DefaultQuery("code", "")
	return strings.ToUpper(code)
}

func WebTypeOption(c *gin.Context) string {
	typ := c.DefaultQuery("type", ANY_TYPE)
	return strings.ToUpper(typ)
}
