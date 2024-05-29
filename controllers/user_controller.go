package controllers

import (
	"backend-api/app"
	"backend-api/database"
	"backend-api/helpers"
	"backend-api/middlewares"
	"backend-api/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// RegisterUser handles the registration of a new user
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user app.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Hash the password
	hashedPassword, err := helpers.HashPassword(user.Password)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}
	user.Password = hashedPassword
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Save the user to the database
	if err := models.CreateUser(&user); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error saving user")
		return
	}

	helpers.RespondWithJSON(w, http.StatusCreated, user)
}

// LoginUser handles user login and JWT generation
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user app.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	log.Printf("Login attempt with email: %s", user.Email)

	storedUser, err := models.GetUserByEmail(user.Email)
	if err != nil {
		log.Printf("Error retrieving user by email: %v", err)
		helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	log.Printf("Stored user: %v", storedUser)

	if !helpers.CheckPasswordHash(user.Password, storedUser.Password) {
		log.Printf("Password mismatch for user: %s", user.Email)
		helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Assign the values of created_at and updated_at to storedUser
	var createdAt time.Time
	var updatedAt time.Time
	query := `SELECT created_at, updated_at FROM users WHERE email = ?`
	err = database.DB.QueryRow(query, user.Email).Scan(&createdAt, &updatedAt)
	if err != nil {
		log.Printf("Error retrieving created_at and updated_at: %v", err)
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error retrieving user data")
		return
	}
	storedUser.CreatedAt = createdAt
	storedUser.UpdatedAt = updatedAt

	token, err := helpers.GenerateJWT(storedUser.ID)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}

// UpdateUser handles updating user information
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user app.User
	params := mux.Vars(r)
	targetUserID := params["userId"]

	// Retrieve the user ID from the context
	userID := r.Context().Value(middlewares.UserContextKey).(uint)

	if strconv.FormatUint(uint64(userID), 10) != targetUserID {
		helpers.RespondWithError(w, http.StatusForbidden, "You can only update your own account")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user.UpdatedAt = time.Now()
	user.ID = userID

	if err := models.UpdateUser(targetUserID, &user); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error updating user")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, user)
}

// DeleteUser handles deleting a user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	params := mux.Vars(r)
	targetUserID, err := strconv.ParseUint(params["userId"], 10, 32)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if userID != uint(targetUserID) {
		helpers.RespondWithError(w, http.StatusForbidden, "You can only delete your own account")
		return
	}

	if err := models.DeleteUser(uint(targetUserID)); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error deleting user")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
