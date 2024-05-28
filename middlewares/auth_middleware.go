package middlewares

import (
	"backend-api/helpers"
	"context"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Missing Authorization header")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // No Bearer prefix
			helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid token format")
			return
		}

		token, err := helpers.ValidateJWT(tokenString)
		if err != nil {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid token: "+err.Error())
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid token claims: user_id not found")
			return
		}

		ctx := context.WithValue(r.Context(), "userID", uint(userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
