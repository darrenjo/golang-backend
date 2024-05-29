package app

import "time"

// Photo represents a user's photo in the system
type Photo struct {
	ID        uint      `json:"id"`
	PhotoURL  string    `json:"photo_url" binding:"required"`
	UserID    uint      `json:"user_id"`
	IsProfile bool      `json:"is_profile"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
