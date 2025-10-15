package authn

import (
	"github.com/roidaradal/fn/clock"
	"github.com/roidaradal/krap"
	"github.com/roidaradal/rdb/ze"
)

var Sessions *ze.Schema[Session]

const (
	sessionActive   string = "ACTIVE"
	sessionExtended string = "EXTENDED"
	sessionExpired  string = "EXPIRED"
	sessionLogout   string = "LOGOUT"
)

type Authable interface {
	GetID() ze.ID
	GetType() string
	GetPassword() string
}

type Session struct {
	ze.UniqueItem
	ze.CreatedItem
	Token
	AccountID     ze.ID  `json:"-"`
	AccountCode_  string `col:"-" json:"AccountCode"`
	LastUpdatedAt ze.DateTime
	ExpiresAt     ze.DateTime
	Status        string
	krap.RequestOrigin
}

type Params struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

type Token struct {
	Type string `validate:"required" fx:"upper"`
	Code string `validate:"required"`
}

func (a Token) String() string {
	return a.Type + authTokenGlue + a.Code
}

func (s Session) IsExpired() bool {
	return clock.CheckIfExpired(s.ExpiresAt)
}
