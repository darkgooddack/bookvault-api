package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/darkgooddack/bookvault-api/config"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const (
	UserIDKey ctxKey = "user_id"
)

var jwtSecret []byte

func Init(secret string) {
	jwtSecret = []byte(secret)
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		uid, ok := claims["user_id"]
		if !ok {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}

		var userID uint
		switch v := uid.(type) {
		case float64:
			userID = uint(v)
		case uint:
			userID = v
		case int:
			userID = uint(v)
		default:
			http.Error(w, "invalid user id in token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserIDFromContext(r *http.Request) (uint, bool) {
	v := r.Context().Value(UserIDKey)
	if v == nil {
		return 0, false
	}
	id, ok := v.(uint)
	return id, ok
}

func InitFromConfig(cfg *config.Config) {
	Init(cfg.JWTSecret)
}
