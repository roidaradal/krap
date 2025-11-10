package krap

import (
	"strings"

	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/fn/io"
)

// Reads the patch object from path as type T, then convert to dict.Object
func CmdReadPatchObject[T any](path string) (dict.Object, error) {
	patchItem, err := io.ReadJSON[T](path)
	if err != nil {
		return nil, err
	}
	return dict.ToObject(patchItem)
}

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

// Return list type (current, archive, future)
func ListTypeOption(option string) string {
	option = strings.ToLower(option)
	listType := LIST_CURRENT
	switch option {
	case LIST_ARCHIVE:
		listType = option
	case LIST_FUTURE:
		listType = option
	}
	return listType
}
