package models

type User struct {
	ID           string   `json:"id"`
	Email        string   `json:"email"`
	Name         string   `json:"name"`
	NavPosition  string   `json:"nav_position"`
	NavItems     []string `json:"nav_items"`
	Theme        string   `json:"theme"`
	SyncEnabled  bool     `json:"sync_enabled"`
	LastSync     string   `json:"last_sync"`
	CreatedAt    string   `json:"created_at"`
	LastSyncTime int64    `json:"last_sync_time"`
}

type UserPreferences struct {
	NavPosition string   `json:"nav_position"`
	NavItems    []string `json:"nav_items"`
	Theme       string   `json:"theme"`
	SyncEnabled bool     `json:"sync_enabled"`
}

type Theme struct {
	Dark     bool
	Primary  string
	Accent   string
	TextSize int
}
