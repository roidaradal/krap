package krap

import "strings"

// False if option is 'all', otherwise true
func MustBeActiveOption(option string) bool {
	mustBeActive := strings.ToLower(option) != viewAll
	return mustBeActive
}

// Return toggle on/off (boolean), hasToggleOption (ok flag)
func ToggleOption(option string) (bool, bool) {
	switch strings.ToLower(option) {
	case toggleOn:
		return true, true
	case toggleOff:
		return false, true
	default:
		return false, false
	}
}

// Returns uppercase type at params[limit] if it exists,
// Defaults to ANY_TYPE (*)
func CmdTypeOption(params []string, limit int) string {
	typ := ANY_TYPE
	if len(params) > limit {
		typ = strings.ToUpper(params[limit])
	}
	return typ
}
