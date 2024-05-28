package models

import (
	"backend-api/app"
	"backend-api/database"
	"database/sql"
	"errors"
)

// CreateUser creates a new user in the database
func CreateUser(user *app.User) error {
	// Use prepared statement for SQL query
	query := `INSERT INTO users (username, password, email, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	_, err := database.DB.Exec(query, user.Username, user.Password, user.Email, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return errors.New("failed to create user: " + err.Error())
	}
	return nil
}

// GetUserByEmail retrieves a user by email from the database
func GetUserByEmail(email string) (app.User, error) {
	var user app.User
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE email = ?`
	err := database.DB.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("user not found")
		}
		return user, errors.New("failed to get user by email: " + err.Error())
	}
	return user, nil
}

// UpdateUser updates a user in the database
func UpdateUser(userID string, user *app.User) error {
	// Use prepared statement for SQL query
	query := `UPDATE users SET username = ?, email = ?, password = ?, updated_at = ? WHERE id = ?`
	_, err := database.DB.Exec(query, user.Username, user.Email, user.Password, user.UpdatedAt, userID)
	if err != nil {
		return errors.New("failed to update user: " + err.Error())
	}
	return nil
}

// DeleteUser deletes a user from the database
func DeleteUser(userID string) error {
	// Use prepared statement for SQL query
	query := `DELETE FROM users WHERE id = ?`
	_, err := database.DB.Exec(query, userID)
	if err != nil {
		return errors.New("failed to delete user: " + err.Error())
	}
	return nil
}
