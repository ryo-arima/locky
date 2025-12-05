package repository

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ryo-arima/locky/pkg/entity/response"
)

type RoleFilter struct{ ID string }

func TrimEndpoint(endpoint string) string {
	return strings.TrimRight(endpoint, "/")
}

// Table formatting helper (used by usecase)
func rolesTableString(res response.RoleResponse) string {
	if res.Code != "SUCCESS" {
		return fmt.Sprintf("Code: %s\nMessage: %s\n", res.Code, res.Message)
	}
	// Single role retrieval (Detail contains permissions)
	if res.Detail != nil {
		// Format permissions display
		permLines := []string{}
		if perms, ok := res.Detail.([]interface{}); ok {
			for _, p := range perms {
				permLines = append(permLines, fmt.Sprintf("  - %v", p))
			}
		}
		return fmt.Sprintf("Role: %v\nPermissions:\n%v\n", res.Roles, strings.Join(permLines, "\n"))
	}
	// List: Sort and enumerate Roles alphabetically
	switch v := res.Roles.(type) {
	case []string:
		cp := make([]string, len(v))
		copy(cp, v)
		sort.Strings(cp)
		return fmt.Sprintf("Roles (%d): %s\n", len(cp), strings.Join(cp, ", "))
	default:
		return fmt.Sprintf("Roles: %v\n", v)
	}
}

// RolesTableStringAlias public function (for display use from other packages)
func RolesTableStringAlias(res response.RoleResponse) string { return rolesTableString(res) }
