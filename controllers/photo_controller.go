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

// SetProfilePhoto sets a photo as the user's profile photo.
func SetProfilePhoto(w http.ResponseWriter, r *http.Request) {
	var photo app.Photo
	if err := json.NewDecoder(r.Body).Decode(&photo); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID, ok := middlewares.GetUserContextKey(r)
	if !ok {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid user ID in token")
		return
	}

	log.Printf("User ID from context: %d", userID)

	// Retrieve the existing profile photo of the user
	existingProfilePhotos, err := models.GetUserProfilePhotos(userID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error fetching profile photos")
		return
	}

	if len(existingProfilePhotos) > 0 {
		// Update the existing profile photo without changing created_at
		existingPhoto := existingProfilePhotos[0]
		existingPhoto.PhotoURL = photo.PhotoURL
		existingPhoto.IsProfile = true
		existingPhoto.UpdatedAt = time.Now()

		if err := models.UpdatePhoto(&existingPhoto); err != nil {
			helpers.RespondWithError(w, http.StatusInternalServerError, "Error updating profile photo")
			return
		}

		helpers.RespondWithJSON(w, http.StatusOK, existingPhoto)
		return
	}

	// Set the new photo as the profile photo
	photo.IsProfile = true
	photo.UserID = userID
	photo.CreatedAt = time.Now()
	photo.UpdatedAt = time.Now()
	if err := models.CreatePhoto(&photo); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error setting profile photo")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, photo)
}

// GetProfilePhoto retrieves the profile photo URL of the user.
func GetProfilePhoto(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserContextKey(r)
	if !ok {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid user ID in token")
		return
	}

	photoURL, err := models.GetUserProfilePhoto(userID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error fetching profile photo")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"photo_url": photoURL})
}

// DeletePhoto menangani penghapusan foto.
func DeletePhoto(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	photoID, err := strconv.ParseUint(params["photoId"], 10, 64)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid photo ID")
		return
	}

	userID, ok := middlewares.GetUserContextKey(r)
	if !ok {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid user ID in token")
		return
	}

	// Mengambil informasi foto dari database
	foundPhoto, err := models.GetPhotoByID(uint(photoID), userID)
	if err != nil {
		// Jika foto tidak ditemukan atau pengguna tidak memiliki hak akses, kirim pesan kesalahan
		helpers.RespondWithError(w, http.StatusNotFound, "Photo not found or user does not have access")
		return
	}

	// Hapus foto hanya jika pengguna memiliki akses ke foto tersebut
	if foundPhoto.UserID != userID {
		helpers.RespondWithError(w, http.StatusForbidden, "You do not have permission to delete this photo")
		return
	}

	// Hapus foto dengan ID tertentu untuk pengguna tertentu
	if err := models.DeletePhoto(uint(photoID), userID); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error deleting photo")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
