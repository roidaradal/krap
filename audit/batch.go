package audit

import (
	"fmt"

	"github.com/roidaradal/fn/clock"
	"github.com/roidaradal/fn/str"
	"github.com/roidaradal/rdb/ze"
)

// Creates new BatchLog
func NewBatchLog(action, details, actionGlue string) *BatchLog {
	now := clock.TimeNow()
	batchLog := &BatchLog{}
	batchLog.CreatedAt = clock.StandardFormat(now)
	batchLog.Code = fmt.Sprintf("%s-%s", clock.TimestampFormat(now), str.SplitInitials(action, actionGlue))
	batchLog.Action = action
	batchLog.Details = details
	return batchLog
}

// Creates rows of new BatchLog items
func NewBatchLogItems(batchCode string, detailsList []string) []*BatchLogItem {
	batchItems := make([]*BatchLogItem, len(detailsList))
	for i, details := range detailsList {
		batchItem := &BatchLogItem{}
		batchItem.Code = batchCode
		batchItem.Details = details
		batchItems[i] = batchItem
	}
	return batchItems
}

// Inserts BatchLog using transaction
func AddBatchLogTx(rqtx *ze.Request, batchLog *BatchLog) error {
	return BatchLogs.InsertTx(rqtx, &ze.AddParams[BatchLog]{
		Item: batchLog,
	})
}

// Inserts BatchLogItems using transaction
func AddBatchLogItemsTx(rqtx *ze.Request, batchItems []*BatchLogItem) error {
	return BatchLogItems.InsertTxRows(rqtx, &ze.AddParams[BatchLogItem]{
		Items: batchItems,
	})
}
