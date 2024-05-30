package middlewares

import (
	"backend-api/helpers"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// Define the key type for user ID in context
type key int

const userContextKey key = 0

// JWTAuth middleware validates JWT tokens and extracts user ID from claims
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
		secret := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			log.Printf("Token validation error: %v", err)
			helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Extract user ID from claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("Invalid token claims")
			helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Convert user ID claim to uint
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			log.Println("Invalid user_id type in token claims")
			helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		userID := uint(userIDFloat)

		log.Printf("User ID from token: %d", userID)

		// Add user ID to request context
		ctx := context.WithValue(r.Context(), userContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserContextKey retrieves the user ID from request context
func GetUserContextKey(r *http.Request) (uint, bool) {
	userID, ok := r.Context().Value(userContextKey).(uint)
	log.Printf("Retrieved user ID from context: %d, exists: %v", userID, ok)
	return userID, ok
}
