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
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Missing auth token")
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := helpers.ValidateJWT(tokenString)
		if err != nil || !token.Valid {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		userID := uint(claims["user_id"].(float64))

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
