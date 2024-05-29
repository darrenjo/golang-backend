package app

import "database/sql"

type Photo struct {
	ID        uint         `json:"id"`
	Title     string       `json:"title" binding:"required"`
	Caption   string       `json:"caption"`
	PhotoURL  string       `json:"photo_url" binding:"required"`
	UserID    uint         `json:"user_id"`
	IsProfile bool         `json:"is_profile"`
	CreatedAt sql.NullTime `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}
