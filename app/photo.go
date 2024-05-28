package app

import "time"

// Photo represents a photo uploaded by a user
type Photo struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title" binding:"required"`
	Caption   string    `json:"caption"`
	PhotoURL  string    `json:"photo_url" binding:"required"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
