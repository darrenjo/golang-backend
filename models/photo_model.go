package models

import (
	"backend-api/app"
	"backend-api/database"
	"log"
	"time"
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

// UpdatePhoto memperbarui informasi foto yang sudah ada.
func UpdatePhoto(photo *app.Photo) error {
	query := `UPDATE photos SET photo_url = ?, is_profile = ?, updated_at = ? WHERE id = ? AND user_id = ?`
	_, err := database.DB.Exec(query, photo.PhotoURL, photo.IsProfile, photo.UpdatedAt, photo.ID, photo.UserID)
	return err
}

// DeletePhoto menghapus foto dengan ID tertentu untuk pengguna tertentu.
func DeletePhoto(photoID, userID uint) error {
	query := `DELETE FROM photos WHERE id = ? AND user_id = ?`
	_, err := database.DB.Exec(query, photoID, userID)
	return err
}

// GetPhotoByID retrieves the details of a photo by photo ID and user ID.
func GetPhotoByID(photoID, userID uint) (*app.Photo, error) {
	var photo app.Photo
	var createdAtStr, updatedAtStr string
	query := `SELECT id, photo_url, user_id, is_profile, created_at, updated_at FROM photos WHERE id = ? AND user_id = ? LIMIT 1`
	err := database.DB.QueryRow(query, photoID, userID).Scan(&photo.ID, &photo.PhotoURL, &photo.UserID, &photo.IsProfile, &createdAtStr, &updatedAtStr)
	if err != nil {
		log.Printf("Error fetching photo with ID %d for user ID %d: %v", photoID, userID, err)
		return nil, err
	}

	// Convert createdAtStr and updatedAtStr to time.Time
	photo.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		log.Printf("Error parsing created_at for photo with ID %d: %v", photoID, err)
		return nil, err
	}
	photo.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAtStr)
	if err != nil {
		log.Printf("Error parsing updated_at for photo with ID %d: %v", photoID, err)
		return nil, err
	}

	return &photo, nil
}

// GetUserProfilePhotos mengembalikan semua foto profil untuk pengguna tertentu.
func GetUserProfilePhotos(userID uint) ([]app.Photo, error) {
	rows, err := database.DB.Query("SELECT id, photo_url, user_id, is_profile, created_at, updated_at FROM photos WHERE user_id = ? AND is_profile = true", userID)
	if err != nil {
		log.Printf("Error fetching profile photos for user ID %d: %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var photos []app.Photo
	for rows.Next() {
		var photo app.Photo
		var createdAtStr, updatedAtStr string
		if err := rows.Scan(&photo.ID, &photo.PhotoURL, &photo.UserID, &photo.IsProfile, &createdAtStr, &updatedAtStr); err != nil {
			log.Printf("Error scanning profile photo for user ID %d: %v", userID, err)
			return nil, err
		}

		// Convert createdAtStr and updatedAtStr to time.Time
		photo.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			log.Printf("Error parsing created_at for photo with ID %d: %v", photo.ID, err)
			return nil, err
		}
		photo.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAtStr)
		if err != nil {
			log.Printf("Error parsing updated_at for photo with ID %d: %v", photo.ID, err)
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
		log.Printf("Error fetching profile photo URL for user ID %d: %v", userID, err)
		return "", err
	}
	return photoURL, nil
}
