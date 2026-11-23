package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"Server-rental/internal/auth"
	"Server-rental/internal/domain"
)

type contextKey string

const UserContextKey contextKey = "user"

type AuthUser struct {
	ID   int
	Role domain.Role
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "требуется токен авторизации"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "неверный формат токена"})
			return
		}

		claims, err := auth.ParseToken(parts[1])
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "недействительный токен"})
			return
		}

		userID := int(claims["user_id"].(float64))
		role := domain.Role(claims["role"].(string))

		ctx := context.WithValue(r.Context(), UserContextKey, AuthUser{
			ID:   userID,
			Role: role,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
