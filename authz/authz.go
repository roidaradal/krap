package authz

import (
	"fmt"
	"slices"
	"strings"

	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/krap"
	"github.com/roidaradal/rdb"
	"github.com/roidaradal/rdb/ze"
)

var (
	appAccess    dict.StringListMap                       // {fullAction => []roles}
	scopedAccess map[string]map[string]dict.StringListMap // map[table][scope][fullAction] => []roles
)

// Initialize authz package
func Initialize() error {
	errs := make([]error, 0)
	appAccess = make(dict.StringListMap)
	scopedAccess = make(map[string]map[string]dict.StringListMap)

	AccessSchema, errs = krap.AddSchema(&Access{}, "config_access", errs)
	ScopedAccessSchema, errs = krap.AddSharedSchema(&ScopedAccess{}, errs)

	if len(errs) > 0 {
		return fmt.Errorf("%d errors encountered: %w", len(errs), errs[0])
	}

	return nil
}

// Load app access list
func LoadAccess(rq *ze.Request) error {
	if AccessSchema == nil {
		rq.Status = ze.Err500
		return errMissingSchema
	}

	q := rdb.NewFullSelectRowsQuery(AccessSchema.Table, AccessSchema.Reader)
	access, err := q.Query(rq.DB)
	if err != nil {
		rq.AddLog("Failed to load app access from db")
		rq.Status = ze.Err500
		return err
	}

	if appAccess == nil {
		appAccess = make(dict.StringListMap)
	}

	for _, axs := range access {
		for _, action := range actionsList {
			fullAction := createFullAction(action, axs.Item)
			isActive := false
			switch action {
			case ROWS:
				isActive = axs.Rows
			case VIEW:
				isActive = axs.View
			case ADD:
				isActive = axs.Add
			case TOGGLE:
				isActive = axs.Toggle
			case EDIT:
				isActive = axs.Edit
			}
			if isActive {
				appAccess[fullAction] = append(appAccess[fullAction], axs.Role)
			}
		}
	}

	rq.Status = ze.OK200
	return nil
}

// Load scoped access list at table
func LoadScopedAccess(rq *ze.Request, table string) error {
	if ScopedAccessSchema == nil {
		rq.Status = ze.Err500
		return errMissingSchema
	}

	q := rdb.NewFullSelectRowsQuery(table, ScopedAccessSchema.Reader)
	access, err := q.Query(rq.DB)
	if err != nil {
		rq.AddFmtLog("Failed to load scoped access from '%s'", table)
		rq.Status = ze.Err500
		return err
	}

	if scopedAccess == nil {
		scopedAccess = make(map[string]map[string]dict.StringListMap)
	}
	scopedAccess[table] = make(map[string]dict.StringListMap)

	for _, axs := range access {
		scope := axs.ScopeCode
		if !dict.HasKey(scopedAccess[table], scope) {
			scopedAccess[table][scope] = make(dict.StringListMap)
		}
		for _, action := range actionsList {
			fullAction := createFullAction(action, axs.Item)
			isActive := false
			switch action {
			case ROWS:
				isActive = axs.Rows
			case VIEW:
				isActive = axs.View
			case ADD:
				isActive = axs.Add
			case TOGGLE:
				isActive = axs.Toggle
			case EDIT:
				isActive = axs.Edit
			}
			if isActive {
				scopedAccess[table][scope][fullAction] = append(scopedAccess[table][scope][fullAction], axs.Role)
			}
		}
	}

	rq.Status = ze.OK200
	return nil
}

// Get all app access
func GetAllAccess() dict.StringListMap {
	return appAccess
}

// Get {role => []fullActions}
func GetAllRoleAccess() dict.StringListMap {
	roleAccess := make(dict.StringListMap)
	for action, roles := range appAccess {
		for _, role := range roles {
			roleAccess[role] = append(roleAccess[role], action)
		}
	}
	for role, actions := range roleAccess {
		slices.Sort(actions)
		roleAccess[role] = actions
	}
	return roleAccess
}

// Get scoped access
func GetScopedAccess(table string, scopeCode string) dict.StringListMap {
	scopeAccess := make(dict.StringListMap)
	if !dict.HasKey(scopedAccess, table) {
		return scopeAccess
	}
	if !dict.HasKey(scopedAccess[table], scopeCode) {
		return scopeAccess
	}
	return scopedAccess[table][scopeCode]
}

// Check if role is allowed to do action on item
func CheckActionAllowedFor(rq *ze.Request, action, item, role string) error {
	role = strings.ToUpper(item)
	fullAction := createFullAction(strings.ToUpper(action), strings.ToUpper(item))
	allowedRoles := appAccess[fullAction]
	if !slices.Contains(allowedRoles, role) {
		rq.Status = ze.Err403
		return ErrUnauthorizedAccess
	}
	return nil
}

// Check if role is allowed to do scoped action on item
func CheckScopedActionAllowedFor(rq *ze.Request, table, scopeCode, action, item, role string) error {
	if !dict.HasKey(scopedAccess, table) {
		rq.Status = ze.Err403
		return ErrUnauthorizedAccess
	}
	scopeCode = strings.ToUpper(scopeCode)
	if !dict.HasKey(scopedAccess[table], scopeCode) {
		rq.Status = ze.Err403
		return ErrUnauthorizedAccess
	}
	role = strings.ToUpper(role)
	fullAction := createFullAction(strings.ToUpper(action), strings.ToUpper(item))
	allowedRoles := scopedAccess[table][scopeCode][fullAction]
	if !slices.Contains(allowedRoles, role) {
		rq.Status = ze.Err403
		return ErrUnauthorizedAccess
	}
	return nil
}

// Common: create Action-Item key
func createFullAction(action, item string) string {
	return fmt.Sprintf("%s-%s", action, item)
}
