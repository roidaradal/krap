package authn

import (
	"time"

	"github.com/roidaradal/fn/clock"
	"github.com/roidaradal/rdb"
	"github.com/roidaradal/rdb/ze"
)

// Common: archives the given session, and deletes from the sessions table
func archiveSession(rq *ze.Request, session *Session, status string) error {
	err := rq.StartTransaction(2)
	if err != nil {
		rq.Status = ze.Err500
		return err
	}
	rqtx := rq
	rqtx.Checker = rdb.AssertRowsAffected(1)
	sessionID := session.ID

	// 1) INSERT to sessions_archive
	session.ID = 0 // for auto-increment
	session.Status = status
	err = Sessions.InsertTxAt(rqtx, session, tableSessionsArchive)
	if err != nil {
		return err
	}

	// 2) DELETE from sessions
	if Sessions == nil {
		rq.Status = ze.Err500
		return ze.ErrMissingSchema
	}
	condition := rdb.Equal(&Sessions.Ref.ID, sessionID)
	err = Sessions.DeleteTx(rqtx, condition)
	if err != nil {
		return err
	}

	// 3) Commit transaction
	err = rqtx.CommitTransaction()
	if err != nil {
		return err
	}

	storeDeleteSession(session)
	return nil
}

// Archives the expired sessions
func archiveExpiredSessions() (*ze.Request, error) {
	rq, err := ze.NewRequest("ArchiveExpiredSessions")
	if err != nil {
		return rq, err
	}

	start := clock.TimeNow()
	defer func() {
		rq.AddDurationLog(start)
	}()

	if Sessions == nil {
		rq.Status = ze.Err500
		return rq, ze.ErrMissingSchema
	}

	condition := rdb.LessEqual(&Sessions.Ref.ExpiresAt, clock.DateTimeNow())
	expired, err := Sessions.GetRows(rq, condition)
	if err != nil {
		rq.AddLog("Failed to get expired sessions")
		rq.Status = ze.Err500
		return rq, err
	}

	numExpired := len(expired)
	rq.AddFmtLog("Expired: %d", numExpired)
	if numExpired == 0 {
		return rq, nil // end early if no expired
	}

	success, fail := 0, 0
	for i, session := range expired {
		rq.AddFmtLog("%.2d / %.2d: ArchiveExpiredSession: %s", i+1, numExpired, session.Token.String())

		err = archiveSession(rq, session, sessionExpired)
		if err != nil {
			fail += 1
			rq.AddFmtLog("Failed: %s", err.Error())
		} else {
			success += 1
		}
	}

	rq.AddFmtLog("Success: %d, Fail: %d", success, fail)
	return rq, nil
}

// Deletes the older archived sessions
func deleteArchivedSessions(marginDays uint) (*ze.Request, error) {
	rq, err := ze.NewRequest("DeleteArchivedSessions")
	if err != nil {
		return rq, err
	}

	now := clock.TimeNow()
	defer func() {
		rq.AddDurationLog(now)
	}()

	margin := time.Duration(marginDays) * 24 * time.Hour
	limitTime := clock.StandardFormat(now.Add(-margin))

	if Sessions == nil {
		rq.Status = ze.Err500
		return rq, ze.ErrMissingSchema
	}

	condition := rdb.LessEqual(&Sessions.Ref.ExpiresAt, limitTime)
	err = Sessions.DeleteAt(rq, condition, tableSessionsArchive)
	if err != nil {
		return rq, err
	}

	return rq, nil
}
