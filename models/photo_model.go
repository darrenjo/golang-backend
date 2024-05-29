package models

import (
	"backend-api/app"
	"backend-api/database"
)

// CreatePhoto menambahkan foto profil baru untuk pengguna.
func CreatePhoto(photo *app.Photo) error {
	query := `INSERT INTO photos (photo_url, user_id, is_profile, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	result, err := database.DB.Exec(query, photo.PhotoURL, photo.UserID, photo.IsProfile, photo.CreatedAt, photo.UpdatedAt)
	if err != nil {
		return err
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	photo.ID = uint(lastInsertID)
	return nil
}

// DeletePhoto menghapus foto dengan ID tertentu untuk pengguna tertentu.
func DeletePhoto(photoID, userID uint) error {
	query := `DELETE FROM photos WHERE id = ? AND user_id = ?`
	_, err := database.DB.Exec(query, photoID, userID)
	return err
}

// GetUserProfilePhotos mengembalikan semua foto profil untuk pengguna tertentu.
func GetUserProfilePhotos(userID uint) ([]app.Photo, error) {
	rows, err := database.DB.Query("SELECT id, photo_url, user_id, is_profile, created_at, updated_at FROM photos WHERE user_id = ? AND is_profile = true", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []app.Photo
	for rows.Next() {
		var photo app.Photo
		if err := rows.Scan(&photo.ID, &photo.PhotoURL, &photo.UserID, &photo.IsProfile, &photo.CreatedAt, &photo.UpdatedAt); err != nil {
			return nil, err
		}
		photos = append(photos, photo)
	}

	return photos, nil
}

// GetUserProfilePhoto mengembalikan URL foto profil untuk pengguna tertentu.
func GetUserProfilePhoto(userID uint) (string, error) {
	var photoURL string
	query := `SELECT photo_url FROM photos WHERE user_id = ? AND is_profile = true LIMIT 1`
	err := database.DB.QueryRow(query, userID).Scan(&photoURL)
	if err != nil {
		return "", err
	}
	return photoURL, nil
}
