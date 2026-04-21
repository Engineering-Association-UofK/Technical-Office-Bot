package middleware

import (
	"net/http"
	"slices"

	"github.com/danielgtaylor/huma/v2"
)

func HumaRequireRole(api huma.API, role Role) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		claims, ok := ctx.Context().Value(userKey).(*AdminClaims)
		if !ok {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized")
			return
		}

		hasAccess := slices.Contains(claims.Roles, RoleSuperAdmin) || slices.Contains(claims.Roles, role)
		if !hasAccess {
			huma.WriteErr(api, ctx, http.StatusForbidden, "Forbidden")
			return
		}

		next(ctx)
	}
}
