package controllers

import (
	"backend-api/app"
	"backend-api/helpers"
	"backend-api/models"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

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

func GetPhotos(w http.ResponseWriter, r *http.Request) {
	photos, err := models.GetAllPhotos()
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error fetching photos")
		return
	}
	helpers.RespondWithJSON(w, http.StatusOK, photos)
}

func UpdatePhoto(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	photoID, err := strconv.Atoi(params["photoId"])
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid photo ID")
		return
	}

	var photo app.Photo
	if err := json.NewDecoder(r.Body).Decode(&photo); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID := r.Context().Value("userID").(uint)
	existingPhoto, err := models.GetPhotoByID(uint(photoID))
	if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, "Photo not found")
		return
	}

	if existingPhoto.UserID != userID {
		helpers.RespondWithError(w, http.StatusForbidden, "You are not allowed to update this photo")
		return
	}

	photo.ID = uint(photoID)
	photo.UserID = userID
	photo.UpdatedAt = time.Now()

	if err := models.UpdatePhoto(&photo); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error updating photo")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, photo)
}

func DeletePhoto(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	photoID, err := strconv.Atoi(params["photoId"])
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid photo ID")
		return
	}

	userID := r.Context().Value("userID").(uint)
	existingPhoto, err := models.GetPhotoByID(uint(photoID))
	if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, "Photo not found")
		return
	}

	if existingPhoto.UserID != userID {
		helpers.RespondWithError(w, http.StatusForbidden, "You are not allowed to delete this photo")
		return
	}

	if err := models.DeletePhoto(uint(photoID)); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error deleting photo")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func SetProfilePhoto(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	photoID, err := strconv.Atoi(params["photoId"])
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid photo ID")
		return
	}

	userID := r.Context().Value("userID").(uint)
	existingPhoto, err := models.GetPhotoByID(uint(photoID))
	if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, "Photo not found")
		return
	}

	if existingPhoto.UserID != userID {
		helpers.RespondWithError(w, http.StatusForbidden, "You are not allowed to set this photo as profile")
		return
	}

	if err := models.UnsetUserProfilePhotos(userID); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error unsetting profile photos")
		return
	}

	existingPhoto.IsProfile = true
	existingPhoto.UpdatedAt = time.Now()

	if err := models.UpdatePhoto(&existingPhoto); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error setting profile photo")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, existingPhoto)
}

func GetProfilePhoto(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	photos, err := models.GetUserProfilePhotos(userID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Error fetching profile photos")
		return
	}
	helpers.RespondWithJSON(w, http.StatusOK, photos)
}
