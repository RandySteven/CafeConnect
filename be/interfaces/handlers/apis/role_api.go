package api_interfaces

import "net/http"

type RoleApi interface {
	GetRoleList(w http.ResponseWriter, r *http.Request)
}
