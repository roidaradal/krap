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

var accessControl ds.StringListMap = nil
var errUnauthorizedAccess = errors.New("public: Unauthorized access")

func AccessControl() ds.StringListMap {
	return accessControl
}

func ActorAccess() ds.StringListMap {
	actorAccess := make(ds.StringListMap)
	for action, actors := range accessControl {
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
	allowed := accessControl[action]
	if len(allowed) == 0 {
		return errUnauthorizedAccess
	}
	isAllowed := slices.Contains(allowed, role)
	return fn.Ternary(isAllowed, nil, errUnauthorizedAccess)
}

func LoadAccessControl(rq *krap.Request, table string, reader rdb.RowReader[rdb.Access]) error {
	a := rdb.Schema(rdb.Access{})
	q := rdb.NewFullSelectRowsQuery(table, reader)
	q.Where(rdb.Equal(&a.IsActive, true))

	access, err := q.Query(rq.DB)
	if err != nil {
		return err
	}

	accessControl = make(ds.StringListMap)
	for _, axs := range access {
		accessControl[axs.Action] = append(accessControl[axs.Action], axs.Role)
	}
	return nil
}
