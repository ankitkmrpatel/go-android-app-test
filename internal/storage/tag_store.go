package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/goBookMarker/internal/models"
)

type TagStore struct {
	db *sql.DB
}

func NewTagStore(db *sql.DB) *TagStore {
	return &TagStore{db: db}
}

func (s *TagStore) InitSchema() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS tags (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			color TEXT NOT NULL,
			description TEXT,
			parent_id TEXT,
			tag_order INTEGER DEFAULT 0,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			usage_stats TEXT,
			FOREIGN KEY(parent_id) REFERENCES tags(id) ON DELETE SET NULL
		)`,
		`CREATE TABLE IF NOT EXISTS tag_groups (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			tag_ids TEXT NOT NULL,
			group_order INTEGER DEFAULT 0,
			expanded BOOLEAN DEFAULT true
		)`,
		`CREATE TABLE IF NOT EXISTS bookmark_tags (
			bookmark_id TEXT NOT NULL,
			tag_id TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			PRIMARY KEY (bookmark_id, tag_id),
			FOREIGN KEY(bookmark_id) REFERENCES bookmarks(id) ON DELETE CASCADE,
			FOREIGN KEY(tag_id) REFERENCES tags(id) ON DELETE CASCADE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_tags_parent_id ON tags(parent_id)`,
		`CREATE INDEX IF NOT EXISTS idx_tags_order ON tags(tag_order)`,
		`CREATE INDEX IF NOT EXISTS idx_tag_groups_order ON tag_groups(group_order)`,
	}

	for _, query := range queries {
		if _, err := s.db.Exec(query); err != nil {
			return fmt.Errorf("failed to create schema: %w", err)
		}
	}
	return nil
}

// CRUD operations for tags
func (s *TagStore) CreateTag(tag models.Tag) error {
	statsJSON, err := json.Marshal(tag.UsageStats)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO tags (id, name, color, description, parent_id, tag_order, created_at, updated_at, usage_stats)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = s.db.Exec(query,
		tag.ID,
		tag.Name,
		tag.Color,
		tag.Description,
		tag.ParentID,
		tag.Order,
		tag.CreatedAt,
		tag.UpdatedAt,
		string(statsJSON),
	)
	return err
}

func (s *TagStore) GetTag(id string) (*models.Tag, error) {
	var tag models.Tag
	var statsJSON string

	query := `
		SELECT id, name, color, description, parent_id, tag_order, created_at, updated_at, usage_stats
		FROM tags WHERE id = ?
	`
	err := s.db.QueryRow(query, id).Scan(
		&tag.ID,
		&tag.Name,
		&tag.Color,
		&tag.Description,
		&tag.ParentID,
		&tag.Order,
		&tag.CreatedAt,
		&tag.UpdatedAt,
		&statsJSON,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(statsJSON), &tag.UsageStats); err != nil {
		return nil, err
	}

	return &tag, nil
}

func (s *TagStore) UpdateTag(tag models.Tag) error {
	statsJSON, err := json.Marshal(tag.UsageStats)
	if err != nil {
		return err
	}

	query := `
		UPDATE tags 
		SET name = ?, color = ?, description = ?, parent_id = ?, 
			tag_order = ?, updated_at = ?, usage_stats = ?
		WHERE id = ?
	`
	_, err = s.db.Exec(query,
		tag.Name,
		tag.Color,
		tag.Description,
		tag.ParentID,
		tag.Order,
		time.Now(),
		string(statsJSON),
		tag.ID,
	)
	return err
}

func (s *TagStore) DeleteTag(id string) error {
	_, err := s.db.Exec("DELETE FROM tags WHERE id = ?", id)
	return err
}

// Tag Group operations
func (s *TagStore) CreateTagGroup(group models.TagGroup) error {
	tagIDsJSON, err := json.Marshal(group.TagIDs)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO tag_groups (id, name, tag_ids, group_order, expanded)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err = s.db.Exec(query,
		group.ID,
		group.Name,
		string(tagIDsJSON),
		group.Order,
		group.Expanded,
	)
	return err
}

func (s *TagStore) GetTagGroup(id string) (*models.TagGroup, error) {
	var group models.TagGroup
	var tagIDsJSON string

	query := `SELECT id, name, tag_ids, group_order, expanded FROM tag_groups WHERE id = ?`
	err := s.db.QueryRow(query, id).Scan(
		&group.ID,
		&group.Name,
		&tagIDsJSON,
		&group.Order,
		&group.Expanded,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(tagIDsJSON), &group.TagIDs); err != nil {
		return nil, err
	}

	return &group, nil
}

func (s *TagStore) UpdateTagGroup(group models.TagGroup) error {
	tagIDsJSON, err := json.Marshal(group.TagIDs)
	if err != nil {
		return err
	}

	query := `
		UPDATE tag_groups 
		SET name = ?, tag_ids = ?, group_order = ?, expanded = ?
		WHERE id = ?
	`
	_, err = s.db.Exec(query,
		group.Name,
		string(tagIDsJSON),
		group.Order,
		group.Expanded,
		group.ID,
	)
	return err
}

func (s *TagStore) DeleteTagGroup(id string) error {
	_, err := s.db.Exec("DELETE FROM tag_groups WHERE id = ?", id)
	return err
}

// Advanced queries
func (s *TagStore) GetTagsByParent(parentID string) ([]models.Tag, error) {
	query := `
		SELECT id, name, color, description, parent_id, tag_order, created_at, updated_at, usage_stats
		FROM tags 
		WHERE parent_id = ?
		ORDER BY tag_order
	`
	rows, err := s.db.Query(query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var tag models.Tag
		var statsJSON string
		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.Color,
			&tag.Description,
			&tag.ParentID,
			&tag.Order,
			&tag.CreatedAt,
			&tag.UpdatedAt,
			&statsJSON,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(statsJSON), &tag.UsageStats); err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}
	return tags, nil
}

func (s *TagStore) GetAllTags() ([]models.Tag, error) {
	query := `
		SELECT id, name, color, description, parent_id, tag_order, created_at, updated_at, usage_stats
		FROM tags 
		ORDER BY tag_order
	`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var tag models.Tag
		var statsJSON string
		err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.Color,
			&tag.Description,
			&tag.ParentID,
			&tag.Order,
			&tag.CreatedAt,
			&tag.UpdatedAt,
			&statsJSON,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(statsJSON), &tag.UsageStats); err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}
	return tags, nil
}

func (s *TagStore) GetAllTagGroups() ([]models.TagGroup, error) {
	query := `SELECT id, name, tag_ids, group_order, expanded FROM tag_groups ORDER BY group_order`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []models.TagGroup
	for rows.Next() {
		var group models.TagGroup
		var tagIDsJSON string
		err := rows.Scan(
			&group.ID,
			&group.Name,
			&tagIDsJSON,
			&group.Order,
			&group.Expanded,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(tagIDsJSON), &group.TagIDs); err != nil {
			return nil, err
		}

		groups = append(groups, group)
	}
	return groups, nil
}

// Import/Export functionality
func (s *TagStore) ExportTags() (*models.TagExport, error) {
	tags, err := s.GetAllTags()
	if err != nil {
		return nil, err
	}

	groups, err := s.GetAllTagGroups()
	if err != nil {
		return nil, err
	}

	return &models.TagExport{
		Tags:       tags,
		TagGroups:  groups,
		ExportedAt: time.Now(),
		Version:    "1.0",
	}, nil
}

func (s *TagStore) ImportTags(export *models.TagExport) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Clear existing data
	if _, err := tx.Exec("DELETE FROM tag_groups"); err != nil {
		return err
	}
	if _, err := tx.Exec("DELETE FROM tags"); err != nil {
		return err
	}

	// Import tags
	for _, tag := range export.Tags {
		statsJSON, err := json.Marshal(tag.UsageStats)
		if err != nil {
			return err
		}

		query := `
			INSERT INTO tags (id, name, color, description, parent_id, tag_order, created_at, updated_at, usage_stats)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
		if _, err := tx.Exec(query,
			tag.ID,
			tag.Name,
			tag.Color,
			tag.Description,
			tag.ParentID,
			tag.Order,
			tag.CreatedAt,
			tag.UpdatedAt,
			string(statsJSON),
		); err != nil {
			return err
		}
	}

	// Import tag groups
	for _, group := range export.TagGroups {
		tagIDsJSON, err := json.Marshal(group.TagIDs)
		if err != nil {
			return err
		}

		query := `
			INSERT INTO tag_groups (id, name, tag_ids, group_order, expanded)
			VALUES (?, ?, ?, ?, ?)
		`
		if _, err := tx.Exec(query,
			group.ID,
			group.Name,
			string(tagIDsJSON),
			group.Order,
			group.Expanded,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}
