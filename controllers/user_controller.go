package controllers

import (
	"backend-api/app"
	"backend-api/helpers"
	"backend-api/models"
	"encoding/json"
	"net/http"
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

	storedUser, err := models.GetUserByEmail(user.Email)
	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if !helpers.CheckPasswordHash(user.Password, storedUser.Password) {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := helpers.GenerateJWT(storedUser.ID)
	if err != nil {
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

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user.UpdatedAt = time.Now()

	if err := models.UpdateUser(userID, &user); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error updating user")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, user)
}

// DeleteUser handles deleting a user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userId"]

	if err := models.DeleteUser(userID); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error deleting user")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}