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
	var createdAt, updatedAt string
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE email = ?`
	err := database.DB.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &createdAt, &updatedAt)
	if err != nil {
		return user, err
	}
	user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	user.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
	return user, nil
}

func UpdateUser(userID string, user *app.User) error {
	query := `UPDATE users SET username = ?, email = ?, password = ?, updated_at = ? WHERE id = ?`
	_, err := database.DB.Exec(query, user.Username, user.Email, user.Password, user.UpdatedAt, userID)
	return err
}

func DeleteUser(userID string) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := database.DB.Exec(query, userID)
	return err
}
