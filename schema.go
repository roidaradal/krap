package krap

import "github.com/roidaradal/rdb/ze"

// Creates a new schema with given table
func AddSchema[T any](item *T, table string, errs []error) (*ze.Schema[T], []error) {
	schema, err := ze.NewSchema(item, table)
	if err != nil {
		errs = append(errs, err)
	}
	return schema, errs
}

// Creates a new shared schema
func AddSharedSchema[T any](item *T, errs []error) (*ze.Schema[T], []error) {
	schema, err := ze.NewSharedSchema(item)
	if err != nil {
		errs = append(errs, err)
	}
	return schema, errs

}
