package models

import (
	"backend-api/app"
	"backend-api/database"
)

func CreatePhoto(photo *app.Photo) error {
	query := `INSERT INTO photos (title, caption, photo_url, user_id, is_profile, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := database.DB.Exec(query, photo.Title, photo.Caption, photo.PhotoURL, photo.UserID, photo.IsProfile, photo.CreatedAt, photo.UpdatedAt)
	return err
}

func GetPhotoByID(photoID uint) (app.Photo, error) {
	var photo app.Photo
	query := `SELECT id, title, caption, photo_url, user_id, is_profile, created_at, updated_at FROM photos WHERE id = ?`
	err := database.DB.QueryRow(query, photoID).Scan(&photo.ID, &photo.Title, &photo.Caption, &photo.PhotoURL, &photo.UserID, &photo.IsProfile, &photo.CreatedAt, &photo.UpdatedAt)
	return photo, err
}

func UpdatePhoto(photo *app.Photo) error {
	query := `UPDATE photos SET title = ?, caption = ?, photo_url = ?, is_profile = ?, updated_at = ? WHERE id = ? AND user_id = ?`
	_, err := database.DB.Exec(query, photo.Title, photo.Caption, photo.PhotoURL, photo.IsProfile, photo.UpdatedAt, photo.ID, photo.UserID)
	return err
}

func DeletePhoto(photoID uint) error {
	query := `DELETE FROM photos WHERE id = ?`
	_, err := database.DB.Exec(query, photoID)
	return err
}

func UnsetUserProfilePhotos(userID uint) error {
	query := `UPDATE photos SET is_profile = false WHERE user_id = ?`
	_, err := database.DB.Exec(query, userID)
	return err
}

func GetAllPhotos() ([]app.Photo, error) {
	rows, err := database.DB.Query("SELECT id, title, caption, photo_url, user_id, is_profile, created_at, updated_at FROM photos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []app.Photo
	for rows.Next() {
		var photo app.Photo
		if err := rows.Scan(&photo.ID, &photo.Title, &photo.Caption, &photo.PhotoURL, &photo.UserID, &photo.IsProfile, &photo.CreatedAt, &photo.UpdatedAt); err != nil {
			return nil, err
		}
		photos = append(photos, photo)
	}

	return photos, nil
}

func GetUserProfilePhotos(userID uint) ([]app.Photo, error) {
	rows, err := database.DB.Query("SELECT id, title, caption, photo_url, user_id, is_profile, created_at, updated_at FROM photos WHERE user_id = ? AND is_profile = true", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []app.Photo
	for rows.Next() {
		var photo app.Photo
		if err := rows.Scan(&photo.ID, &photo.Title, &photo.Caption, &photo.PhotoURL, &photo.UserID, &photo.IsProfile, &photo.CreatedAt, &photo.UpdatedAt); err != nil {
			return nil, err
		}
		photos = append(photos, photo)
	}

	return photos, nil
}
