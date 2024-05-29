package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// RespondWithError responds with a JSON error message
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

// RespondWithJSON responds with a JSON payload
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// HashPassword hashes a plain-text password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash checks if the hashed password matches the plain-text password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT generates a JWT token for a user
func GenerateJWT(userID uint) (string, error) {
	// Get JWT secret from environment variable
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET environment variable is not set")
	}

	// Create a new JWT token with user ID as a claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("error generating JWT token: %v", err)
	}

	fmt.Printf("Generated JWT token: %s\n", tokenString)
	return tokenString, nil
}

// ValidateJWT validates the JWT token
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Printf("Error parsing token: %v\n", err)
		return nil, err
	}

	// Check if token is valid and contains the user_id claim
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["user_id"] == nil {
			fmt.Println("user_id claim is missing in the token")
			return nil, errors.New("invalid token: user_id claim is missing")
		}
		fmt.Printf("Token claims: %v\n", claims)
	}

	return token, nil
}
