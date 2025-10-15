package authn

import (
	"strings"

	"github.com/roidaradal/fn/clock"
	"github.com/roidaradal/fn/hash"
	"github.com/roidaradal/fn/str"
	"github.com/roidaradal/krap"
	"github.com/roidaradal/rdb"
	"github.com/roidaradal/rdb/ze"
)

// Function hook for post-authentication actions to account
type PostAuthHook[T any] = func(*ze.Request, *T)

// Authenticate account
func authenticateAccount[T Authable](rq *ze.Request, params *Params, schema *ze.Schema[T], condition rdb.Condition) (*T, error) {
	if condition == nil {
		rq.AddLog("Missing condition for authenticate account")
		rq.Status = ze.Err500
		return nil, ze.ErrMissingParams
	}

	q := rdb.NewFullSelectRowQuery(schema.Table, schema.Reader)
	q.Where(condition)

	account, err := q.QueryRow(rq.DB)
	if err != nil {
		rq.AddErrorLog(err)
		rq.Status = ze.Err401
		return nil, ErrNotFoundAccount
	}

	hashPassword := (*account).GetPassword()
	if ok := hash.MatchPassword(params.Password, hashPassword); !ok {
		rq.AddLog("Password doesnt match")
		rq.Status = ze.Err401
		return nil, ErrFailedAuth
	}

	return account, nil
}

// Creates new session
func newSession[T Authable](rq *ze.Request, accountRef *T, origin *krap.RequestOrigin) (*Session, error) {
	if Sessions == nil {
		rq.Status = ze.Err500
		return nil, ze.ErrMissingSchema
	}

	if accountRef == nil {
		rq.Status = ze.Err400
		return nil, ze.ErrMissingParams
	}
	account := *accountRef
	accountID := account.GetID()

	// Prepare session object
	var browserInfo *string = nil
	var ipAddress *string = nil
	if origin != nil {
		browserInfo = str.NonEmptyRefString(origin.BrowserInfo)
		ipAddress = str.NonEmptyRefString(origin.IPAddress)
	}
	now, expiry := clock.DateTimeNowWithExpiry(sessionDuration)

	s := &Session{}
	s.ID = 0
	s.CreatedAt = now
	s.Type = strings.ToUpper(account.GetType())
	s.Code = str.RandomString(sessionCodeLength, true, true, true)
	s.AccountID = accountID
	s.LastUpdatedAt = now
	s.ExpiresAt = expiry
	s.Status = sessionActive
	s.BrowserInfo = browserInfo
	s.IPAddress = ipAddress

	// Insert session
	id, err := Sessions.InsertID(rq, s)
	if err != nil {
		return nil, err
	}
	s.ID = id
	storeAddSession(s) // add to session store

	return s, nil
}
