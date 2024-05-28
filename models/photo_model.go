package models

import (
	"backend-api/app"
	"backend-api/database"
)

// CreatePhoto saves a new photo in the database
func CreatePhoto(photo *app.Photo) error {
	query := `INSERT INTO photos (title, caption, photo_url, user_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := database.DB.Exec(query, photo.Title, photo.Caption, photo.PhotoURL, photo.UserID, photo.CreatedAt, photo.UpdatedAt)
	return err
}

// GetAllPhotos retrieves all photos from the database
func GetAllPhotos() ([]app.Photo, error) {
	var photos []app.Photo
	rows, err := database.DB.Query(`SELECT id, title, caption, photo_url, user_id, created_at, updated_at FROM photos`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var photo app.Photo
		if err := rows.Scan(&photo.ID, &photo.Title, &photo.Caption, &photo.PhotoURL, &photo.UserID, &photo.CreatedAt, &photo.UpdatedAt); err != nil {
			return nil, err
		}
		photos = append(photos, photo)
	}
	return photos, nil
}

// GetPhotoByID retrieves a photo by its ID
func GetPhotoByID(id string) (*app.Photo, error) {
	var photo app.Photo
	query := `SELECT id, title, caption, photo_url, user_id, created_at, updated_at FROM photos WHERE id = ?`
	err := database.DB.QueryRow(query, id).Scan(&photo.ID, &photo.Title, &photo.Caption, &photo.PhotoURL, &photo.UserID, &photo.CreatedAt, &photo.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &photo, nil
}

// UpdatePhoto updates an existing photo in the database
func UpdatePhoto(photo *app.Photo) error {
	query := `UPDATE photos SET title = ?, caption = ?, photo_url = ?, updated_at = ? WHERE id = ?`
	_, err := database.DB.Exec(query, photo.Title, photo.Caption, photo.PhotoURL, photo.UpdatedAt, photo.ID)
	return err
}

// DeletePhoto deletes a photo from the database
func DeletePhoto(id string) error {
	query := `DELETE FROM photos WHERE id = ?`
	_, err := database.DB.Exec(query, id)
	return err
}
