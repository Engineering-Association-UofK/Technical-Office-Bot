package handler

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type AdminClaims struct {
	ID    float64  `json:"id"`
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Type  string   `json:"type"`
	Roles []string `json:"roles"`

	jwt.RegisteredClaims
}

type contextKey string

const userKey contextKey = "user_data"

// Open endpoint anyone can call
func Basic(next http.Handler) http.Handler {
	return recoveryMiddleware(
		loggingMiddleware(
			corsMiddleware(
				next,
			),
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
			requireType("admin")(next),
		),
	)
}

// Role based access for admins
func RoleProtected(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return Protected(
			requireRole(role)(next),
		)
	}
}

func AnyRoleProtected(roles []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return Protected(
			requireRoles(roles)(next),
		)
	}
}

func requireRole(role string) func(http.Handler) http.Handler {
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

func requireRoles(roles []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(userKey).(*AdminClaims)

			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if func(base []string, toCheck []string) bool {
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

func requireType(Type string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(userKey).(*AdminClaims)

			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if claims.Type != Type {
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

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
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

		secretBytes, err := base64.RawStdEncoding.DecodeString(config.App.JwtSecret)
		if err != nil {
			slog.Error("Failed to decode JWT secret from Base64", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

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
			slog.Debug("Invalid Token: " + err.Error())
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
