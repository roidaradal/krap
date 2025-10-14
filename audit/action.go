package audit

import (
	"strings"

	"github.com/roidaradal/fn/clock"
	"github.com/roidaradal/rdb/ze"
)

const detailsGlue string = "|"

// Creates new details string from items separated by |
func NewDetails(items ...string) string {
	return strings.Join(items, detailsGlue)
}

// Creates new ActionLog
func NewActionLog(actorID ze.ID, action, details string) *ActionLog {
	actionLog := &ActionLog{}
	actionLog.CreatedAt = clock.DateTimeNow()
	actionLog.ActorID = actorID
	actionLog.Action = action
	actionLog.Details = details
	return actionLog
}

// Creates list of new ActionLogs
func NewActionLogs(actorID ze.ID, actionDetails [][2]string) []*ActionLog {
	actionLogs := make([]*ActionLog, len(actionDetails))
	for i, pair := range actionDetails {
		action, details := pair[0], pair[1]
		actionLogs[i] = NewActionLog(actorID, action, details)
	}
	return actionLogs
}

// Inserts ActionLog using transaction at given table
func AddActionLogTx(rqtx *ze.Request, actionLog *ActionLog, table string) error {
	if ActionLogs == nil {
		rqtx.Status = ze.Err500
		return errMissingSchema
	}
	return ActionLogs.InsertTxAt(rqtx, actionLog, table)
}

// Inserts ActionLog rows using transaction at given table
func AddActionLogsTx(rqtx *ze.Request, actionLogs []*ActionLog, table string) error {
	if ActionLogs == nil {
		rqtx.Status = ze.Err500
		return errMissingSchema
	}
	return ActionLogs.InsertTxRowsAt(rqtx, actionLogs, table)
}
