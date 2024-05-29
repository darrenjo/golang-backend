package models

import (
	"backend-api/app"
	"backend-api/database"
	"time"
)

func CreateUser(user *app.User) error {
	query := `INSERT INTO users (username, password, email, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	_, err := database.DB.Exec(query, user.Username, user.Password, user.Email, user.CreatedAt, user.UpdatedAt)
	return err
}

func GetUserByEmail(email string) (app.User, error) {
	var user app.User
	var createdAt time.Time
	var updatedAt time.Time
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE email = ?`
	err := database.DB.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &createdAt, &updatedAt)
	if err != nil {
		return user, err
	}
	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt
	return user, nil
}

func UpdateUser(userID string, user *app.User) error {
	query := `UPDATE users SET username = ?, email = ?, password = ?, updated_at = ? WHERE id = ?`
	_, err := database.DB.Exec(query, user.Username, user.Email, user.Password, user.UpdatedAt, userID)
	return err
}

func DeleteUser(userID uint) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := database.DB.Exec(query, userID)
	return err
}
