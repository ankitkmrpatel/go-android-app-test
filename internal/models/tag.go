package models

import (
	"encoding/json"
	"time"
)

type Tag struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Color       string    `json:"color"`
	Description string    `json:"description"`
	ParentID    string    `json:"parent_id,omitempty"` // For tag hierarchies
	Order       int       `json:"order"`               // For custom ordering
	Count       int       `json:"count"`               // Number of bookmarks using this tag
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UsageStats  TagStats  `json:"usage_stats"`
}

type TagStats struct {
	LastUsed        time.Time `json:"last_used"`
	UsageCount      int       `json:"usage_count"`      // Total times this tag was used
	BookmarkCount   int       `json:"bookmark_count"`   // Current number of bookmarks
	HistoricalCount int       `json:"historical_count"` // Total bookmarks ever tagged
}

type TagGroup struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	TagIDs   []string `json:"tag_ids"`
	Order    int      `json:"order"`
	Expanded bool     `json:"expanded"`
}

// TagExport represents a collection of tags and tag groups for import/export
type TagExport struct {
	Tags       []Tag       `json:"tags"`
	TagGroups  []TagGroup  `json:"tag_groups"`
	ExportedAt time.Time   `json:"exported_at"`
	Version    string      `json:"version"`
}

func (t *Tag) MarshalJSON() ([]byte, error) {
	type Alias Tag
	return json.Marshal(&struct {
		*Alias
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}{
		Alias:     (*Alias)(t),
		CreatedAt: t.CreatedAt.Format(time.RFC3339),
		UpdatedAt: t.UpdatedAt.Format(time.RFC3339),
	})
}

func (t *Tag) UnmarshalJSON(data []byte) error {
	type Alias Tag
	aux := &struct {
		*Alias
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var err error
	t.CreatedAt, err = time.Parse(time.RFC3339, aux.CreatedAt)
	if err != nil {
		return err
	}
	t.UpdatedAt, err = time.Parse(time.RFC3339, aux.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// Helper methods for tag operations
func (t *Tag) UpdateStats(bookmarkID string, operation string) {
	now := time.Now()
	t.UpdatedAt = now
	t.UsageStats.LastUsed = now

	switch operation {
	case "add":
		t.UsageStats.UsageCount++
		t.UsageStats.BookmarkCount++
		t.UsageStats.HistoricalCount++
	case "remove":
		t.UsageStats.BookmarkCount--
	}
}

func (t *Tag) IsParent() bool {
	return t.ParentID == ""
}

func (t *Tag) IsChild() bool {
	return t.ParentID != ""
}

// Tag group operations
func (g *TagGroup) AddTag(tagID string) {
	if !g.ContainsTag(tagID) {
		g.TagIDs = append(g.TagIDs, tagID)
	}
}

func (g *TagGroup) RemoveTag(tagID string) {
	for i, id := range g.TagIDs {
		if id == tagID {
			g.TagIDs = append(g.TagIDs[:i], g.TagIDs[i+1:]...)
			break
		}
	}
}

func (g *TagGroup) ContainsTag(tagID string) bool {
	for _, id := range g.TagIDs {
		if id == tagID {
			return true
		}
	}
	return false
}
