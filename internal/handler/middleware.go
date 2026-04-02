package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type Role string

type AdminClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Roles    []Role `json:"roles"`

	jwt.RegisteredClaims
}

type contextKey string

const userKey contextKey = "user_data"

// Open endpoint anyone can call
func Basic(next http.Handler) http.Handler {
	return recoveryMiddleware(
		loggingMiddleware(
			next,
		),
	)
}

// Protected endpoints that only someone with a valid JWT can call
func BasicProtected(next http.Handler) http.Handler {
	return Basic(
		authMiddleware(
			next,
		),
	)
}

// Only Admins can call this one
func Protected(next http.Handler) http.Handler {
	return Basic(
		authMiddleware(
			requireRole("sys:admin")(next),
		),
	)
}

// Role based access for admins
func RoleProtected(role Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return Protected(
			requireRole(role)(next),
		)
	}
}

func AnyRoleProtected(roles []Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return Protected(
			requireRoles(roles)(next),
		)
	}
}

func requireRole(role Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(userKey).(*AdminClaims)

			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if !slices.Contains(claims.Roles, role) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func requireRoles(roles []Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(userKey).(*AdminClaims)

			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if func(base []Role, toCheck []Role) bool {
				for b := range base {
					for t := range toCheck {
						if b == t {
							return true
						}
					}
				}
				return false
			}(claims.Roles, roles) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		slog.Debug("Logging middleware", "Method", r.Method, "Path", r.URL.Path, "Time took", time.Since(start))
	})
}

// Checks the validation of the JWT token made by the Spring backend
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		secretBytes := []byte(config.App.JwtSecret)

		// Get token and claims
		claims := &AdminClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return secretBytes, nil
		})

		// Check validity
		if err != nil || !token.Valid {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				msg := fmt.Sprintf("Recovered from panic: %v\n", err)
				slog.Error(msg)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
}
