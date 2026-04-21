package middleware

import (
	"net/http"
	"strings"

	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/config"
	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Role string

const (
	RoleAdmin       Role = "sys:admin"
	RoleSuperAdmin  Role = "sys:super_admin"
	RoleTechSupport Role = "sys:tech_support"
)

type AdminClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Roles    []Role `json:"roles"`

	jwt.RegisteredClaims
}

type contextKey string

const userKey contextKey = "user_data"

func HumaAuth(api huma.API) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		tokenString := ctx.Header("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized")
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		claims := &AdminClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.App.JwtSecret), nil
		})

		if err != nil || !token.Valid {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Invalid Token", err)
			return
		}

		// Update the context with claims for the handlers to use
		ctx = huma.WithValue(ctx, userKey, claims)
		next(ctx)
	}
}
