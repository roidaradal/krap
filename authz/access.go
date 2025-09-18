package authz

import (
	"errors"
	"slices"
	"strings"

	"github.com/roidaradal/fn"
	"github.com/roidaradal/fn/ds"
	"github.com/roidaradal/krap"
	"github.com/roidaradal/rdb"
)

var accessList ds.StringListMap = nil
var errUnauthorizedAccess = errors.New("public: Unauthorized access")

func AccessList() ds.StringListMap {
	return accessList
}

func ActorAccess() ds.StringListMap {
	actorAccess := make(ds.StringListMap)
	for action, actors := range accessList {
		for _, actor := range actors {
			actorAccess[actor] = append(actorAccess[actor], action)
		}
	}
	for actor, actions := range actorAccess {
		slices.Sort(actions)
		actorAccess[actor] = actions
	}
	return actorAccess
}

func CheckActionAllowed(action string, role string) error {
	role = strings.ToUpper(role)
	action = strings.ToUpper(action)
	allowed := accessList[action]
	if len(allowed) == 0 {
		return errUnauthorizedAccess
	}
	isAllowed := slices.Contains(allowed, role)
	return fn.Ternary(isAllowed, nil, errUnauthorizedAccess)
}

func LoadAccessList(rq *krap.Request, AccessSchema *rdb.Schema[rdb.Access]) error {
	a := AccessSchema.Ref
	q := rdb.NewFullSelectRowsQuery(AccessSchema.Table, AccessSchema.Reader)
	q.Where(rdb.Equal(&a.IsActive, true))

	access, err := q.Query(rq.DB)
	if err != nil {
		rq.AddLog("Failed to load access list from db")
		return err
	}

	accessList = make(ds.StringListMap)
	for _, axs := range access {
		accessList[axs.Action] = append(accessList[axs.Action], axs.Role)
	}
	return nil
}
