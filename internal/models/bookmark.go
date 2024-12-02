package models

import (
	"time"
)

type Bookmark struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	FaviconURL  string    `json:"favicon_url"`
	IsFavorite  bool      `json:"is_favorite"`
	Tags        []string  `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
