// Package konfig contains app config related functions and types
package konfig

import (
	"github.com/roidaradal/fn/fail"
	"github.com/roidaradal/rdb/ze"
)

// Initialize the config package
func Initialize() error {
	errs := make([]error, 0)

	KVSchema, errs = ze.AddSchema(&KV{}, "config_app", errs)

	if len(errs) > 0 {
		return fail.FromErrors("konfig.Initialize", errs)
	}

	return nil
}
