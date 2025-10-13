package audit

import (
	"fmt"

	"github.com/roidaradal/krap"
	"github.com/roidaradal/rdb/ze"
)

var (
	ActionLogs    *ze.Schema[ActionLog]
	BatchLogs     *ze.Schema[BatchLog]
	BatchLogItems *ze.Schema[BatchLogItem]
)

var errs []error = nil

// Initialize audit package
func Initialize() error {
	errs = make([]error, 0)

	ActionLogs, errs = krap.AddSharedSchema(&ActionLog{}, errs)
	BatchLogs, errs = krap.AddSchema(&BatchLog{}, "logs_batch", errs)
	BatchLogItems, errs = krap.AddSchema(&BatchLogItem{}, "logs_batch_items", errs)

	if len(errs) > 0 {
		return fmt.Errorf("%d errors encountered: %w", len(errs), errs[0])
	}

	return nil
}
