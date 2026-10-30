package auth

import "slices"

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

func HasRole(roles []string, want Role) bool {
	return slices.Contains(roles, string(want))
}

func CanAccessTenant(claims *Claims, tenantID string) bool {
	if claims == nil {
		return false
	}
	if HasRole(claims.Roles, RoleAdmin) {
		return true
	}
	return claims.TenantID == tenantID
}
