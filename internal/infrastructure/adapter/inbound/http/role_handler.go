package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/apex20/backend/internal/domain/permission"
)

type roleEntry struct {
	Value int32  `json:"value"`
	Name  string `json:"name"`
}

type listRolesOutput struct {
	Body struct {
		Roles []roleEntry `json:"roles"`
	}
}

// RegisterRoleHandler registers the /admin/roles route on the given API.
func RegisterRoleHandler(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "list-roles",
		Method:      http.MethodGet,
		Path:        "/admin/roles",
		Summary:     "List Roles",
		Tags:        []string{"Admin"},
	}, func(_ context.Context, _ *struct{}) (*listRolesOutput, error) {
		roles := []permission.Role{
			permission.RoleGM,
			permission.RolePlayer,
			permission.RoleTrusted,
		}
		out := &listRolesOutput{}
		out.Body.Roles = make([]roleEntry, len(roles))
		for i, r := range roles {
			out.Body.Roles[i] = roleEntry{
				Value: int32(r),
				Name:  r.String(),
			}
		}
		return out, nil
	})
}
