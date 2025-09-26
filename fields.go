package krap

// Field, OldValue, NewValue
type FieldUpdate [3]string

func (f FieldUpdate) Tuple() (string, string, string) {
	return f[0], f[1], f[2]
}

func HasNullableUpdate[T comparable](oldItem *T, newItem *T, hasUpdate bool) bool {
	if !hasUpdate {
		return false
	}
	if oldItem == nil && newItem == nil {
		// Both nil = no update
		return false
	} else if oldItem != nil && newItem != nil {
		return *oldItem != *newItem
	} else {
		// One is nil, the other is not = has update
		return true
	}
}
