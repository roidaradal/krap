package krap

import (
	"strings"

	"github.com/roidaradal/fn/ds"
	"github.com/roidaradal/fn/str"
)

const (
	DEFAULT_OPTION string = "."
	ANY_TYPE       string = "*"
	toggleOn       string = "on"
	toggleOff      string = "off"
	viewAll        string = "all"
	viewActive     string = "active"
	okMessage      string = "OK"
	envDev         string = "dev"
	envProd        string = "prod"
)

// Initializer function with name
type Initializer struct {
	Fn   func() error
	Name string
}

// Web request origin: BrowserInfo and IP address
type RequestOrigin struct {
	BrowserInfo *string
	IPAddress   *string
}

// Response for creating multiple items
type BulkCreateResult[T any] struct {
	BulkActionResult
	Items *ds.List[*T]
}

// Response for performing action on multiple items
type BulkActionResult struct {
	Success int
	Fail    int
	Fails   []string
}

// Checks if app env is 'dev' or 'prod'
func IsValidAppEnv(appEnv string) bool {
	return appEnv == envDev || appEnv == envProd
}

// Check if app env is 'prod'
func IsProdAppEnv(appEnv string) bool {
	return appEnv == envProd
}

// Get public error message from "public: <message>"
func publicErrorMessage(err error) (string, bool) {
	msg := err.Error()
	if strings.HasPrefix(msg, "public: ") {
		parts := str.CleanSplit(msg, ":")
		return parts[1], true
	}
	return "Unexpected error during request", false
}
