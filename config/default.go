package config

import (
	"github.com/roidaradal/fn"
	"github.com/roidaradal/fn/conv"
	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/fn/str"
)

type Defaults struct {
	UintMap       map[string]uint
	IntMap        map[string]int
	StringMap     dict.StringMap
	StringListMap dict.StringListMap
}

// Tries to convert lookup[key] to uint, fallsback to defaultValue[key]
func uintOrDefault(lookup dict.StringMap, defaultValue map[string]uint, key string) uint {
	value := defaultValue[key]
	if lookupValue, ok := lookup[key]; ok {
		value = uint(conv.ParseInt(lookupValue))
	}
	return value
}

// Tries to convert lookup[key] to int, fallsback to defaultValue[key]
func intOrDefault(lookup dict.StringMap, defaultValue map[string]int, key string) int {
	value := defaultValue[key]
	if lookupValue, ok := lookup[key]; ok {
		value = conv.ParseInt(lookupValue)
	}
	return value
}

// Tries to get lookup[key], fallsback to defaultValue[key]
func stringOrDefault(lookup dict.StringMap, defaultValue dict.StringMap, key string) string {
	value, ok := lookup[key]
	return fn.Ternary(ok, value, defaultValue[key])
}

// Tries to convert lookup[key] to []string, fallsback to defaultValue[key]
func stringListOrDefault(lookup dict.StringMap, defaultValue dict.StringListMap, key string) []string {
	value, ok := lookup[key]
	return fn.Ternary(ok, str.CleanSplit(value, listGlue), defaultValue[key])
}
