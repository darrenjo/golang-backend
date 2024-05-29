package controllers

import (
	"backend-api/app"
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
	userID := params["userId"]

	// Retrieve user ID from request context
	ctxUserID, ok := middlewares.GetUserContextKey(r)
	if !ok {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid user ID in token")
		return
	}

	// Check if the requested user ID matches the user ID from the token
	if strconv.Itoa(int(ctxUserID)) != userID {
		helpers.RespondWithError(w, http.StatusForbidden, "You can only update your own account")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user.UpdatedAt = time.Now()

	if err := models.UpdateUser(userID, &user); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error updating user")
		return
	}

	// Set the user ID to the ID from the URL
	user.ID = ctxUserID

	helpers.RespondWithJSON(w, http.StatusOK, user)
}

// DeleteUser handles deleting a user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userId"]

	// Retrieve user ID from request context
	ctxUserID, ok := middlewares.GetUserContextKey(r)
	if !ok {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid user ID in token")
		return
	}

	// Check if the requested user ID matches the user ID from the token
	if strconv.Itoa(int(ctxUserID)) != userID {
		helpers.RespondWithError(w, http.StatusForbidden, "You can only delete your own account")
		return
	}

	if err := models.DeleteUser(userID); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error deleting user")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
