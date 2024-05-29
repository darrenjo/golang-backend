package middlewares

import (
	"backend-api/helpers"
	"context"
	"log"
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
			log.Println("Missing auth token")
			helpers.RespondWithError(w, http.StatusUnauthorized, "Missing auth token")
			return
		}

		// Remove "Bearer " prefix from token string
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		log.Printf("Token after trimming prefix: %s", tokenString)

		// Validate JWT token
		token, err := helpers.ValidateJWT(tokenString)
		if err != nil || !token.Valid {
			log.Printf("Token validation error: %v", err)
			helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Extract user ID from claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			log.Println("Invalid token claims")
			helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Check if user_id claim exists and is valid
		userIDClaim, ok := claims["user_id"]
		if !ok {
			log.Println("user_id claim is missing")
			helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		userID, ok := userIDClaim.(float64)
		if !ok {
			log.Printf("Invalid user_id type: %v", userIDClaim)
			helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		log.Printf("User ID from token: %f", userID)

		// Add user ID to request context
		ctx := context.WithValue(r.Context(), userContextKey, uint(userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserContextKey retrieves the user ID from request context
func GetUserContextKey(r *http.Request) (uint, bool) {
	userID, ok := r.Context().Value(userContextKey).(uint)
	log.Printf("Retrieved user ID from context: %d, exists: %v", userID, ok)
	return userID, ok
}
