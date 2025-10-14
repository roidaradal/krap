package authz

import (
	"errors"

	"github.com/roidaradal/rdb/ze"
)

var (
	errMissingSchema      = errors.New("schema is not initialized")
	ErrUnauthorizedAccess = errors.New("public: Unauthorized access")
)

var (
	AccessSchema       *ze.Schema[Access]
	ScopedAccessSchema *ze.Schema[ScopedAccess]
)

type Access struct {
	Item string `fx:"upper"`
	Role string `fx:"upper"`
	Actions
}

type ScopedAccess struct {
	ScopeCode string `fx:"upper"`
	Access
}

type Actions struct {
	Rows   bool
	View   bool
	Add    bool
	Toggle bool
	Edit   bool
}

const (
	ROWS   string = "ROWS"
	VIEW   string = "VIEW"
	ADD    string = "ADD"
	TOGGLE string = "TOGGLE"
	EDIT   string = "EDIT"
)

var actionsList = []string{ROWS, VIEW, ADD, TOGGLE, EDIT}
