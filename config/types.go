package config

import (
	"errors"

	"github.com/roidaradal/rdb/ze"
)

var (
	errMissingSchema      = errors.New("kv schema is not initialized")
	errUnavailableFeature = errors.New("public: Unavailable feature")
	errUnknownFeature     = errors.New("public: Unknown feature")
)

var (
	KVSchema       *ze.Schema[KV]
	Features       *ze.Schema[Feature]
	ScopedFeatures *ze.Schema[ScopedFeature]
)

type KV struct {
	Key           string `col:"AppKey"`
	Value         string `col:"AppValue"`
	LastUpdatedAt ze.DateTime
}

type Feature struct {
	ze.ActiveItem
	Name string `fx:"upper" col:"Feature"`
}

type ScopedFeature struct {
	Feature
	ScopeCode string `fx:"upper"`
}
