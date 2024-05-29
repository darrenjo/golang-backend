package middlewares

import (
	"backend-api/helpers"
	"context"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// Define the key type for user ID in context
type key int

const userContextKey key = 0

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Missing auth token")
			return
		}

		// Remove "Bearer " prefix from token string
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Validate JWT token
		token, err := helpers.ValidateJWT(tokenString)
		if err != nil || !token.Valid {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Extract user ID from claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}
		userID := uint(claims["user_id"].(float64))

		// Add user ID to request context
		ctx := context.WithValue(r.Context(), userContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserContextKey retrieves the user ID from request context
func GetUserContextKey(r *http.Request) (uint, bool) {
	userID, ok := r.Context().Value(userContextKey).(uint)
	return userID, ok
}
