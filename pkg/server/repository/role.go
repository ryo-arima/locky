package repository

import (
	"errors"
	"sort"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

// RolePermission represents one permission definition (resource, action).
// resource = "users" / "groups" / "roles" etc., action = "read" / "write" etc.
// Corresponds to obj=obj(resource), act=action in Casbin policy (p, sub, obj, act).
type RolePermission struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

// RoleRepository uses Casbin policy as storage (currently file policy.csv + memory)
// for application-level Role CRUD abstraction.
// Note: Does not use DB here, Casbin Enforcer is the single source of truth.
// - ListRoles(): enumerate all subjects(roles) appearing in policy
// - GetRolePermissions(): list of (resource,action) pairs for 1 role
// - CreateRole(): check existing role duplication + add permissions
// - UpdateRole(): delete all existing policy role lines → recreate with new permissions
// - DeleteRole(): delete all role lines
// All call SavePolicy() for persistence to file (when adapter is file).
type Role interface {
	ListRoles(c *gin.Context) ([]string, error)
	GetRolePermissions(c *gin.Context, role string) ([]RolePermission, error)
	CreateRole(c *gin.Context, role string, perms []RolePermission) error
	UpdateRole(c *gin.Context, role string, perms []RolePermission) error
	DeleteRole(c *gin.Context, role string) error
}

type role struct {
	// App-wide permissions: etc/casbin/locky/*. Used for role management.
	appEnforcer *casbin.Enforcer
}

// NewRole: receives Enforcer and returns repository for role management.
//
//	appEnf -> etc/casbin/locky/model.conf + policy.csv (app-wide RBAC / CRUD target)
func NewRole(appEnf *casbin.Enforcer) Role {
	return &role{appEnforcer: appEnf}
}

// internal helper: returns the Enforcer that is currently the CRUD target.
func (rcvr *role) target() *casbin.Enforcer { return rcvr.appEnforcer }

// ListRoles: enumerate subjects in group policy
func (rcvr *role) ListRoles(c *gin.Context) ([]string, error) {
	subs, err := rcvr.target().GetAllSubjects()
	if err != nil {
		return nil, err
	}
	m := make(map[string]struct{})
	for _, s := range subs {
		// Skip if empty after trimming whitespace (for Casbin internal use, etc.)
		if strings.TrimSpace(s) == "" {
			continue
		}
		m[s] = struct{}{}
	}
	roles := make([]string, 0, len(m))
	for k := range m {
		roles = append(roles, k)
	}
	sort.Strings(roles)
	return roles, nil
}

// GetRolePermissions: extract (resource,action) from all policy lines for specified role.
// Policy storage format: p, <role>, <resource>, <action>
// Error if role is unspecified/blank.
func (rcvr *role) GetRolePermissions(c *gin.Context, role string) ([]RolePermission, error) {
	if strings.TrimSpace(role) == "" {
		return nil, errors.New("role required")
	}
	pols, err := rcvr.target().GetPolicy()
	if err != nil {
		return nil, err
	}
	res := []RolePermission{}
	for _, p := range pols {
		if len(p) < 3 { // Minimum length guarantee for p, sub, obj, act (depends on model)
			continue
		}
		if p[0] == role {
			res = append(res, RolePermission{Resource: p[1], Action: p[2]})
		}
	}
	return res, nil
}

// roleExists: simple check that role exists if 1 or more permissions exist.
// (A role with empty permissions is considered conceptually non-existent)
func (rcvr *role) roleExists(c *gin.Context, role string) bool {
	perms, _ := rcvr.GetRolePermissions(c, role)
	return len(perms) > 0
}

// CreateRole: add new role. Error if duplicate exists.
// If perms is empty, automatically grant roles:read as fallback,
// providing minimum permission to view self (roles list).
func (rcvr *role) CreateRole(c *gin.Context, role string, perms []RolePermission) error {
	role = strings.TrimSpace(role)
	if role == "" {
		return errors.New("role name required")
	}
	if rcvr.roleExists(c, role) {
		return errors.New("role already exists")
	}
	if len(perms) == 0 {
		perms = []RolePermission{{Resource: "group_info", Action: "read"}}
	}
	for _, pm := range perms {
		if _, err := rcvr.target().AddPolicy(role, pm.Resource, pm.Action); err != nil {
			return err
		}
	}
	return rcvr.target().SavePolicy()
}

// UpdateRole: delete all permissions for existing role with RemoveFilteredPolicy, then re-insert new permissions.
// When perms is empty, grant roles:read as minimum permission like Create.
// Note: Full replacement, not differential update.
func (rcvr *role) UpdateRole(c *gin.Context, role string, perms []RolePermission) error {
	role = strings.TrimSpace(role)
	if role == "" {
		return errors.New("role name required")
	}
	if _, err := rcvr.target().RemoveFilteredPolicy(0, role); err != nil {
		return err
	}
	if len(perms) == 0 {
		perms = []RolePermission{{Resource: "group_info", Action: "read"}}
	}
	for _, pm := range perms {
		if _, err := rcvr.target().AddPolicy(role, pm.Resource, pm.Action); err != nil {
			return err
		}
	}
	return rcvr.target().SavePolicy()
}

// DeleteRole: delete all policy lines for specified role. If non-existent, leave to RemoveFilteredPolicy result.
func (rcvr *role) DeleteRole(c *gin.Context, role string) error {
	role = strings.TrimSpace(role)
	if role == "" {
		return errors.New("role name required")
	}
	if _, err := rcvr.target().RemoveFilteredPolicy(0, role); err != nil {
		return err
	}
	return rcvr.target().SavePolicy()
}
