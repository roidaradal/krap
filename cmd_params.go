package krap

import "strings"

const (
	ANY_TYPE   string = "*"
	TOGGLE_ON  string = "on"
	TOGGLE_OFF string = "off"
	listAll    string = "all"
	listActive string = "active"
)

func MustListBeActive(option string) bool {
	mustBeActive := strings.ToLower(option) != listAll
	return mustBeActive
}

// on/off, hasToggleOption
func ToggleOption(option string) (bool, bool) {
	switch strings.ToLower(option) {
	case TOGGLE_ON:
		return true, true
	case TOGGLE_OFF:
		return false, true
	default:
		return false, false
	}
}

func CmdTypeOption(params []string) string {
	typ := ANY_TYPE
	if len(params) > 0 {
		typ = strings.ToUpper(params[0])
	}
	return typ
}
