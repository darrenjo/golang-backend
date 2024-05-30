package helpers

import (
	"backend-api/database"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// RespondWithError memberikan respon dengan pesan error dalam format JSON
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

// RespondWithJSON memberikan respon dengan payload dalam format JSON
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// HashPassword membuat hash dari password dalam plain-text
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash memeriksa apakah hashed password cocok dengan plain-text password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT menghasilkan token JWT untuk seorang pengguna
func GenerateJWT(userID uint) (string, error) {
	// Dapatkan JWT secret dari environment variable
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("variabel lingkungan JWT_SECRET tidak diatur")
	}

	// Buat token JWT baru dengan userID sebagai claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	// Tandatangani token dengan kunci rahasia
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("kesalahan menghasilkan token JWT: %v", err)
	}

	return tokenString, nil
}

// ValidateJWT memvalidasi token JWT
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode tanda tangan tidak terduga: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	// Periksa apakah token valid dan berisi klaim user_id
	if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return nil, errors.New("token tidak valid")
	}

	return token, nil
}

// DeletePhotoFileFromDB menghapus foto dari database MySQL berdasarkan ID foto.
func DeletePhotoFileFromDB(photoID uint) error {
	// Buatlah koneksi ke database Anda (database.DB)

	// Lakukan query untuk mendapatkan informasi foto berdasarkan ID
	query := "SELECT photo_url FROM photos WHERE id = ?"
	var photoURL string
	err := database.DB.QueryRow(query, photoID).Scan(&photoURL)
	if err != nil {
		return err
	}

	// Hapus foto dari database
	deleteQuery := "DELETE FROM photos WHERE id = ?"
	_, err = database.DB.Exec(deleteQuery, photoID)
	if err != nil {
		return err
	}

	return nil
}
