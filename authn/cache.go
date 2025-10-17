package authn

import (
	"github.com/roidaradal/fn"
	"github.com/roidaradal/fn/clock"
	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/rdb"
	"github.com/roidaradal/rdb/ze"
)

var (
	sessionStore   = dict.NewSyncMap[string, *Session]()
	sessionUpdates = dict.NewSyncMap[string, [2]ze.DateTime]()
)

// Initialize the authn session store and updates map
func InitializeStore() error {
	rq, err := ze.NewRequest("LoadSessions")
	if err != nil {
		return err
	}

	if Sessions == nil {
		return ze.ErrMissingSchema
	}

	s := Sessions.Ref
	condition := rdb.Greater(&s.ExpiresAt, clock.DateTimeNow())

	sessions, err := Sessions.GetRows(rq, condition)
	if err != nil {
		return err
	}

	storeAddSessions(sessions)
	return nil
}

// Adds sessions to session store
func storeAddSessions(sessions []*Session) {
	for _, session := range sessions {
		storeAddSession(session)
	}
}

// Add one session to session store
func storeAddSession(session *Session) {
	if session == nil {
		return
	}
	authToken := session.Token.String()
	sessionStore.Set(authToken, session)
}

// Get session from session store using authToken string as key
func storeGetSession(authToken string) *Session {
	session, ok := sessionStore.Get(authToken)
	return fn.Ternary(ok, session, nil)
}

// Deletes session from session store
func storeDeleteSession(session *Session) {
	if session == nil {
		return
	}
	authToken := session.Token.String()
	sessionStore.Delete(authToken)
	sessionUpdates.Delete(authToken)
}

// Adds a session extension to session updates
func storeExtendSession(authToken string, now, expiry ze.DateTime) error {
	sessionUpdates.Set(authToken, [2]ze.DateTime{now, expiry})
	session, ok := sessionStore.Get(authToken)
	if ok {
		session.ExpiresAt = expiry
		sessionStore.Set(authToken, session)
	}
	return nil
}
