package config

import (
	"slices"
	"strings"

	"github.com/roidaradal/fn"
	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/rdb"
	"github.com/roidaradal/rdb/ze"
)

var (
	appFeatures    dict.BoolMap
	scopedFeatures map[string]dict.StringListMap
)

// Load app features
func LoadFeatures(rq *ze.Request) error {
	if Features == nil {
		rq.Status = ze.Err500
		return errMissingSchema
	}
	f := Features.Ref
	q := rdb.NewLookupQuery[Feature](Features.Table, &f.Name, &f.IsActive)
	lookup, err := q.Lookup(rq.DB)
	if err != nil {
		rq.AddLog("Failed to load features from db")
		rq.Status = ze.Err500
		return err
	}
	appFeatures = lookup
	rq.Status = ze.OK200
	return nil
}

// Load active scoped features at table
func LoadScopedFeatures(rq *ze.Request, table string) error {
	if ScopedFeatures == nil {
		rq.Status = ze.Err500
		return errMissingSchema
	}

	f := ScopedFeatures.Ref
	q := rdb.NewFullSelectRowsQuery(table, ScopedFeatures.Reader)
	q.Where(rdb.Equal(&f.IsActive, true))

	features, err := q.Query(rq.DB)
	if err != nil {
		rq.AddFmtLog("Failed to load scoped features from '%s'", table)
		rq.Status = ze.Err500
		return err
	}

	if scopedFeatures == nil {
		scopedFeatures = make(map[string]dict.StringListMap)
	}
	scopedFeatures[table] = make(dict.StringListMap)

	for _, f := range features {
		scope, feature := f.ScopeCode, f.Name
		scopedFeatures[table][scope] = append(scopedFeatures[table][scope], feature)
	}

	rq.Status = ze.OK200
	return nil
}

// Get map[Feature]IsActive
func GetAllFeatures() dict.BoolMap {
	return appFeatures
}

// Get list of active app features
func GetActiveFeatures() []string {
	if appFeatures == nil {
		return nil
	}
	activeFeatures := fn.FilterMap(appFeatures, func(feature string, isActive bool) bool {
		return isActive
	})
	return dict.Keys(activeFeatures)
}

// Get all active {scope => []features} at table
func GetAllScopedFeatures(table string) dict.StringListMap {
	if !dict.HasKey(scopedFeatures, table) {
		return nil
	}
	return scopedFeatures[table]
}

// Get all active {feature => []scopes} at table
func GetAllFeatureScopes(table string) dict.StringListMap {
	featureScopes := make(dict.StringListMap)
	for scope, features := range scopedFeatures[table] {
		for _, feature := range features {
			featureScopes[feature] = append(featureScopes[feature], scope)
		}
	}
	for feature, scopes := range featureScopes {
		slices.Sort(scopes)
		featureScopes[feature] = scopes
	}
	return featureScopes
}

// Get all active scoped features at table for given scopeCodes
func GetScopedFeatures(table string, scopeCodes ...string) dict.StringListMap {
	miniScopedFeatures := make(dict.StringListMap)
	if !dict.HasKey(scopedFeatures, table) {
		return miniScopedFeatures
	}
	for _, scope := range scopeCodes {
		scope = strings.ToUpper(scope)
		features := scopedFeatures[table][scope]
		if len(features) == 0 {
			features = []string{}
		}
		slices.Sort(features)
		miniScopedFeatures[scope] = features
	}
	return miniScopedFeatures
}

// Check if feature is available
func CheckFeature(feature string) error {
	feature = strings.ToUpper(feature)
	isActive, ok := appFeatures[feature]
	if !ok {
		return errUnknownFeature
	}
	return fn.Ternary(isActive, nil, errUnavailableFeature)
}

// Check if scoped feature at table is available
func CheckScopedFeature(table, scope, feature string) error {
	if !dict.HasKey(scopedFeatures, table) {
		return errUnavailableFeature
	}
	scope = strings.ToUpper(scope)
	feature = strings.ToUpper(feature)
	enabled := scopedFeatures[table][scope]
	if len(enabled) == 0 {
		return errUnavailableFeature
	}
	isEnabled := slices.Contains(enabled, feature)
	return fn.Ternary(isEnabled, nil, errUnavailableFeature)
}
