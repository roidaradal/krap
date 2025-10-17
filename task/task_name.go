package task

import (
	"strings"

	"github.com/roidaradal/fn/str"
)

const actionGlue string = "-"

// Gets the core item name, removes trailing "%s" if any
func itemPrefix(item string) string {
	if isCompleteItem(item) {
		return item
	}
	parts := str.CleanSplit(item, actionGlue)
	return strings.Join(parts[:len(parts)-1], actionGlue)
}

// Common: Checks if name ends in "%s"
func isCompleteItem(item string) bool {
	return !strings.HasPrefix(item, "%s")
}
