package konfig

import (
	"fmt"
	"strings"

	"github.com/roidaradal/fn"
	"github.com/roidaradal/fn/conv"
	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/fn/dyn"
	"github.com/roidaradal/fn/str"
	"github.com/roidaradal/krap"
	"github.com/roidaradal/rdb"
	"github.com/roidaradal/rdb/ze"
)

const (
	keyGlue  string = "."
	listGlue string = "|"
)

// Initialize config package
func Initialize() error {
	errs := make([]error, 0)
	appFeatures = make(dict.BoolMap)
	scopedFeatures = make(map[string]dict.StringListMap)

	KVSchema, errs = krap.AddSchema(&KV{}, "config_app", errs)
	Features, errs = krap.AddSchema(&Feature{}, "config_features", errs)
	ScopedFeatures, errs = krap.AddSharedSchema(&ScopedFeature{}, errs)

	if len(errs) > 0 {
		return fmt.Errorf("%d errors encountered: %w", len(errs), errs[0])
	}

	return nil
}

// Load Config lookup from database
func Lookup(rq *ze.Request, appKeys []string) (dict.StringMap, error) {
	if KVSchema == nil {
		rq.Status = ze.Err500
		return nil, ze.ErrMissingSchema
	}
	kv := KVSchema.Ref
	q := rdb.NewLookupQuery[KV](KVSchema.Table, &kv.Key, &kv.Value)
	q.Where(rdb.In(&kv.Key, appKeys))
	lookup, err := q.Lookup(rq.DB)
	if err != nil {
		rq.AddLog("Failed to load app config from db")
		rq.Status = ze.Err500
		return nil, err
	}
	rq.Status = ze.OK200
	return lookup, nil
}

// Decorates a Config object with the contents of lookup
func Create[T any](cfg *T, lookup dict.StringMap, defaults *Defaults) *T {
	for key := range defaults.UintMap {
		value := uintOrDefault(lookup, defaults.UintMap, key)
		dyn.SetFieldValue(cfg, getKey(key), value)
	}
	for key := range defaults.IntMap {
		value := intOrDefault(lookup, defaults.IntMap, key)
		dyn.SetFieldValue(cfg, getKey(key), value)
	}
	for key := range defaults.StringMap {
		value := stringOrDefault(lookup, defaults.StringMap, key)
		dyn.SetFieldValue(cfg, getKey(key), value)
	}
	for key := range defaults.StringListMap {
		value := stringListOrDefault(lookup, defaults.StringListMap, key)
		dyn.SetFieldValue(cfg, getKey(key), value)
	}
	return cfg
}

// Extract the second part of <Domain>.<Key>
func getKey(fullKey string) string {
	parts := strings.Split(fullKey, keyGlue)
	if len(parts) != 2 {
		return fullKey
	}
	return parts[1]
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
