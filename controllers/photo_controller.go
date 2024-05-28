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

// CreatePhoto handles the creation of a new photo
func CreatePhoto(w http.ResponseWriter, r *http.Request) {
	var photo app.Photo
	if err := json.NewDecoder(r.Body).Decode(&photo); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID := r.Context().Value("userID").(uint)
	photo.UserID = userID
	photo.CreatedAt = time.Now()
	photo.UpdatedAt = time.Now()

	if err := models.CreatePhoto(&photo); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error saving photo")
		return
	}

	helpers.RespondWithJSON(w, http.StatusCreated, photo)
}

// GetPhotos handles fetching all photos
func GetPhotos(w http.ResponseWriter, r *http.Request) {
	photos, err := models.GetAllPhotos()
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error fetching photos")
		return
	}
	helpers.RespondWithJSON(w, http.StatusOK, photos)
}

// UpdatePhoto handles updating a photo
func UpdatePhoto(w http.ResponseWriter, r *http.Request) {
	var photo app.Photo
	params := mux.Vars(r)
	photoID := params["photoId"]

	if err := json.NewDecoder(r.Body).Decode(&photo); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID := r.Context().Value("userID").(uint)
	storedPhoto, err := models.GetPhotoByID(photoID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, "Photo not found")
		return
	}

	if storedPhoto.UserID != userID {
		helpers.RespondWithError(w, http.StatusForbidden, "You are not authorized to update this photo")
		return
	}

	photo.ID = storedPhoto.ID
	photo.UserID = storedPhoto.UserID
	photo.CreatedAt = storedPhoto.CreatedAt
	photo.UpdatedAt = time.Now()

	if err := models.UpdatePhoto(&photo); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error updating photo")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, photo)
}

// DeletePhoto handles deleting a photo
func DeletePhoto(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	photoID := params["photoId"]

	userID := r.Context().Value("userID").(uint)
	storedPhoto, err := models.GetPhotoByID(photoID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, "Photo not found")
		return
	}

	if storedPhoto.UserID != userID {
		helpers.RespondWithError(w, http.StatusForbidden, "You are not authorized to delete this photo")
		return
	}

	if err := models.DeletePhoto(photoID); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error deleting photo")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
