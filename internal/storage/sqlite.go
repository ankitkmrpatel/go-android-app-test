package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"

	"github.com/goBookMarker/internal/models"
)

type SQLiteDB struct {
	db *sql.DB
}

const (
	createUsersTable = `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE,
		name TEXT,
		nav_position TEXT DEFAULT 'bottom',
		nav_items TEXT,
		theme TEXT DEFAULT 'system',
		sync_enabled BOOLEAN DEFAULT false,
		last_sync TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	createBookmarksTable = `
	CREATE TABLE IF NOT EXISTS bookmarks (
		id TEXT PRIMARY KEY,
		user_id TEXT,
		url TEXT NOT NULL,
		title TEXT,
		description TEXT,
		image_url TEXT,
		favicon_url TEXT,
		is_favorite BOOLEAN DEFAULT false,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id)
	)`

	createTagsTable = `
	CREATE TABLE IF NOT EXISTS tags (
		id TEXT PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		color TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	createBookmarkTagsTable = `
	CREATE TABLE IF NOT EXISTS bookmark_tags (
		bookmark_id TEXT,
		tag_id TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY(bookmark_id, tag_id),
		FOREIGN KEY(bookmark_id) REFERENCES bookmarks(id) ON DELETE CASCADE,
		FOREIGN KEY(tag_id) REFERENCES tags(id) ON DELETE CASCADE
	)`
)

func NewSQLiteDB() (*SQLiteDB, error) {
	db, err := sql.Open("sqlite", "bookmarker.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	sqlite := &SQLiteDB{db: db}
	if err := sqlite.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %v", err)
	}

	return sqlite, nil
}

func (s *SQLiteDB) Close() error {
	return s.db.Close()
}

func (s *SQLiteDB) initSchema() error {
	tables := []string{
		createUsersTable,
		createBookmarksTable,
		createTagsTable,
		createBookmarkTagsTable,
	}

	for _, table := range tables {
		if _, err := s.db.Exec(table); err != nil {
			return fmt.Errorf("failed to create table: %v", err)
		}
	}
	return nil
}

func (s *SQLiteDB) GetCurrentUser() (*models.User, error) {
	var user models.User
	var navItemsJSON string

	err := s.db.QueryRow(`
		SELECT id, email, name, nav_position, nav_items, theme, sync_enabled, last_sync
		FROM users LIMIT 1
	`).Scan(&user.ID, &user.Email, &user.Name, &user.NavPosition, &navItemsJSON, &user.Theme, &user.SyncEnabled, &user.LastSync)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(navItemsJSON), &user.NavItems); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *SQLiteDB) SaveUser(user *models.User) error {
	navItemsJSON, err := json.Marshal(user.NavItems)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		INSERT INTO users (id, email, name, nav_position, nav_items, theme, sync_enabled, last_sync)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			email = excluded.email,
			name = excluded.name,
			nav_position = excluded.nav_position,
			nav_items = excluded.nav_items,
			theme = excluded.theme,
			sync_enabled = excluded.sync_enabled,
			last_sync = excluded.last_sync
	`, user.ID, user.Email, user.Name, user.NavPosition, navItemsJSON, user.Theme, user.SyncEnabled, user.LastSync)

	return err
}

func (s *SQLiteDB) GetRecentBookmarks(limit int) ([]models.Bookmark, error) {
	rows, err := s.db.Query(`
		SELECT b.id, b.url, b.title, b.description, b.image_url, b.favicon_url, b.is_favorite,
			   b.created_at, b.updated_at, GROUP_CONCAT(t.name) as tags
		FROM bookmarks b
		LEFT JOIN bookmark_tags bt ON b.id = bt.bookmark_id
		LEFT JOIN tags t ON bt.tag_id = t.id
		GROUP BY b.id
		ORDER BY b.created_at DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookmarks []models.Bookmark
	for rows.Next() {
		var b models.Bookmark
		var tags sql.NullString
		err := rows.Scan(&b.ID, &b.URL, &b.Title, &b.Description, &b.ImageURL,
			&b.FaviconURL, &b.IsFavorite, &b.CreatedAt, &b.UpdatedAt, &tags)
		if err != nil {
			return nil, err
		}

		if tags.Valid {
			b.Tags = splitTags(tags.String)
		}
		bookmarks = append(bookmarks, b)
	}

	return bookmarks, nil
}

func (s *SQLiteDB) SaveBookmark(b models.Bookmark) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert or update bookmark
	_, err = tx.Exec(`
		INSERT INTO bookmarks (id, user_id, url, title, description, image_url, favicon_url, is_favorite, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(id) DO UPDATE SET
			url = excluded.url,
			title = excluded.title,
			description = excluded.description,
			image_url = excluded.image_url,
			favicon_url = excluded.favicon_url,
			is_favorite = excluded.is_favorite,
			updated_at = CURRENT_TIMESTAMP
	`, b.ID, b.UserID, b.URL, b.Title, b.Description, b.ImageURL, b.FaviconURL, b.IsFavorite)
	if err != nil {
		return err
	}

	// Delete existing tags
	_, err = tx.Exec("DELETE FROM bookmark_tags WHERE bookmark_id = ?", b.ID)
	if err != nil {
		return err
	}

	// Insert tags and bookmark-tag relationships
	for _, tag := range b.Tags {
		var tagID string
		err := tx.QueryRow("SELECT id FROM tags WHERE name = ?", tag).Scan(&tagID)
		if err == sql.ErrNoRows {
			// Create new tag
			tagID = generateID()
			_, err = tx.Exec("INSERT INTO tags (id, name) VALUES (?, ?)", tagID, tag)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		// Create bookmark-tag relationship
		_, err = tx.Exec("INSERT INTO bookmark_tags (bookmark_id, tag_id) VALUES (?, ?)", b.ID, tagID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *SQLiteDB) SearchBookmarks(query string) ([]models.Bookmark, error) {
	rows, err := s.db.Query(`
		SELECT DISTINCT b.id, b.url, b.title, b.description, b.image_url, b.favicon_url,
			   b.is_favorite, b.created_at, b.updated_at, GROUP_CONCAT(t.name) as tags
		FROM bookmarks b
		LEFT JOIN bookmark_tags bt ON b.id = bt.bookmark_id
		LEFT JOIN tags t ON bt.tag_id = t.id
		WHERE b.title LIKE ? OR b.description LIKE ? OR b.url LIKE ? OR t.name LIKE ?
		GROUP BY b.id
		ORDER BY b.updated_at DESC
	`, "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookmarks []models.Bookmark
	for rows.Next() {
		var b models.Bookmark
		var tags sql.NullString
		err := rows.Scan(&b.ID, &b.URL, &b.Title, &b.Description, &b.ImageURL,
			&b.FaviconURL, &b.IsFavorite, &b.CreatedAt, &b.UpdatedAt, &tags)
		if err != nil {
			return nil, err
		}

		if tags.Valid {
			b.Tags = splitTags(tags.String)
		}
		bookmarks = append(bookmarks, b)
	}

	return bookmarks, nil
}

func (s *SQLiteDB) GetAllTags() ([]models.Tag, error) {
	rows, err := s.db.Query(`
		SELECT id, name, color, created_at,
			   (SELECT COUNT(*) FROM bookmark_tags WHERE tag_id = tags.id) as count
		FROM tags
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var t models.Tag
		err := rows.Scan(&t.ID, &t.Name, &t.Color, &t.CreatedAt, &t.Count)
		if err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}

	return tags, nil
}

func (s *SQLiteDB) UpdateUser(user *models.User) error {
	navItemsJSON, err := json.Marshal(user.NavItems)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		UPDATE users SET
			nav_position = ?,
			nav_items = ?,
			theme = ?,
			sync_enabled = ?,
			last_sync = ?
		WHERE id = ?
	`, user.NavPosition, navItemsJSON, user.Theme, user.SyncEnabled, time.Now().Format(time.RFC3339), user.ID)

	return err
}

func (s *SQLiteDB) CreateTag(tag models.Tag) error {
	query := `INSERT INTO tags (id, name, color, description, parent_id, order_num, count, created_at, updated_at, usage_stats)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	statsJSON, err := json.Marshal(tag.UsageStats)
	if err != nil {
		return fmt.Errorf("failed to marshal usage stats: %w", err)
	}

	_, err = s.db.Exec(query,
		tag.ID,
		tag.Name,
		tag.Color,
		tag.Description,
		tag.ParentID,
		tag.Order,
		tag.Count,
		tag.CreatedAt,
		tag.UpdatedAt,
		statsJSON,
	)
	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}
	return nil
}

func (s *SQLiteDB) UpdateTag(tag models.Tag) error {
	query := `UPDATE tags SET name = ?, color = ?, description = ?, parent_id = ?, 
			  order_num = ?, count = ?, updated_at = ?, usage_stats = ?
			  WHERE id = ?`
	
	statsJSON, err := json.Marshal(tag.UsageStats)
	if err != nil {
		return fmt.Errorf("failed to marshal usage stats: %w", err)
	}

	_, err = s.db.Exec(query,
		tag.Name,
		tag.Color,
		tag.Description,
		tag.ParentID,
		tag.Order,
		tag.Count,
		tag.UpdatedAt,
		statsJSON,
		tag.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update tag: %w", err)
	}
	return nil
}

func (s *SQLiteDB) DeleteTag(tagID string) error {
	query := `DELETE FROM tags WHERE id = ?`
	_, err := s.db.Exec(query, tagID)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}
	return nil
}

func (s *SQLiteDB) GetTagsByBookmark(bookmarkID string) ([]models.Tag, error) {
	query := `SELECT t.* FROM tags t
			  JOIN bookmark_tags bt ON bt.tag_id = t.id
			  WHERE bt.bookmark_id = ?`
	
	rows, err := s.db.Query(query, bookmarkID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var tag models.Tag
		var statsJSON []byte
		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.Color,
			&tag.Description,
			&tag.ParentID,
			&tag.Order,
			&tag.Count,
			&tag.CreatedAt,
			&tag.UpdatedAt,
			&statsJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}

		if err := json.Unmarshal(statsJSON, &tag.UsageStats); err != nil {
			return nil, fmt.Errorf("failed to unmarshal usage stats: %w", err)
		}

		tags = append(tags, tag)
	}
	return tags, nil
}

// Helper functions
func splitTags(tags string) []string {
	if tags == "" {
		return nil
	}
	return strings.Split(tags, ",")
}

func generateID() string {
	return uuid.New().String()
}
