package config

import (
	"fmt"
	"strings"

	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/fn/dyn"
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
		return nil, errMissingSchema
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
