package app

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/goBookMarker/internal/models"
)

type AppState struct {
	mu          sync.RWMutex
	bookmarks   []models.Bookmark
	currentUser *models.User
	searchQuery string
	currentPage string
	tags        []models.Tag
	tagGroups   []models.TagGroup
}

func NewAppState() *AppState {
	return &AppState{
		bookmarks: make([]models.Bookmark, 0),
		tags:      make([]models.Tag, 0),
		tagGroups: make([]models.TagGroup, 0),
	}
}

func (s *AppState) GetBookmarks() []models.Bookmark {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.bookmarks
}

func (s *AppState) Search(query string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.searchQuery = query
}

func (s *AppState) ShowAddBookmark() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentPage = "add_bookmark"
}

func (s *AppState) CurrentUser() *models.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentUser
}

func (s *AppState) SaveUser(user *models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentUser = user
	return nil
}

func (s *AppState) SetCurrentPage(page string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentPage = page
}

func (s *AppState) GetCurrentPage() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentPage
}

func (s *AppState) SaveBookmark(bookmark *models.Bookmark) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, b := range s.bookmarks {
		if b.ID == bookmark.ID {
			s.bookmarks[i] = *bookmark
			return nil
		}
	}
	s.bookmarks = append(s.bookmarks, *bookmark)
	return nil
}

func (s *AppState) EditBookmark(bookmark *models.Bookmark) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentPage = "edit_bookmark"
}

func (s *AppState) DeleteBookmark(bookmark *models.Bookmark) {
	s.mu.Lock()
	defer s.mu.Unlock()
	newBookmarks := make([]models.Bookmark, 0, len(s.bookmarks)-1)
	for _, b := range s.bookmarks {
		if b.ID != bookmark.ID {
			newBookmarks = append(newBookmarks, b)
		}
	}
	s.bookmarks = newBookmarks
}

func (s *AppState) ShareBookmark(bookmark *models.Bookmark) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Implementation will be added later
}

func (s *AppState) GetTags() []models.Tag {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tags
}

func (s *AppState) GetTagsWithLocking() []models.Tag {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.tags
}

func (s *AppState) SaveTag(tag *models.Tag) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, t := range s.tags {
		if t.ID == tag.ID {
			s.tags[i] = *tag
			return nil
		}
	}
	s.tags = append(s.tags, *tag)
	return nil
}

func (s *AppState) DeleteTag(tag *models.Tag) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	newTags := make([]models.Tag, 0, len(s.tags)-1)
	for _, t := range s.tags {
		if t.ID != tag.ID {
			newTags = append(newTags, t)
		}
	}
	s.tags = newTags
	return nil
}

func (s *AppState) DeleteTags(tagIDs []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create a map for quick lookup of tag IDs to delete
	toDelete := make(map[string]bool)
	for _, id := range tagIDs {
		toDelete[id] = true
	}

	// Filter out the tags to be deleted
	newTags := make([]models.Tag, 0, len(s.tags)-len(tagIDs))
	for _, tag := range s.tags {
		if !toDelete[tag.ID] {
			newTags = append(newTags, tag)
		}
	}

	// Update tag groups to remove deleted tags
	for i, group := range s.tagGroups {
		newTagIDs := make([]string, 0, len(group.TagIDs))
		for _, tagID := range group.TagIDs {
			if !toDelete[tagID] {
				newTagIDs = append(newTagIDs, tagID)
			}
		}
		s.tagGroups[i].TagIDs = newTagIDs
	}

	s.tags = newTags
	return nil
}

func (s *AppState) GetTagGroups() []models.TagGroup {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tagGroups
}

func (s *AppState) SaveTagGroup(group *models.TagGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, existingGroup := range s.tagGroups {
		if existingGroup.ID == group.ID {
			s.tagGroups[i] = *group
			return nil
		}
	}
	s.tagGroups = append(s.tagGroups, *group)
	return nil
}

func (s *AppState) DeleteTagGroup(group *models.TagGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, existingGroup := range s.tagGroups {
		if existingGroup.ID == group.ID {
			s.tagGroups = append(s.tagGroups[:i], s.tagGroups[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("tag group not found")
}

// SearchTags searches for tags based on a query string
func (s *AppState) SearchTags(query string) []models.Tag {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if query == "" {
		return s.tags
	}

	var filteredTags []models.Tag
	for _, tag := range s.tags {
		// Case-insensitive search across tag name, description
		if strings.Contains(strings.ToLower(tag.Name), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(tag.Description), strings.ToLower(query)) {
			filteredTags = append(filteredTags, tag)
		}
	}

	return filteredTags
}

func (s *AppState) Logout() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentUser = nil
	// Optionally reset other state if needed
	s.bookmarks = make([]models.Bookmark, 0)
	s.tags = make([]models.Tag, 0)
	s.tagGroups = make([]models.TagGroup, 0)
	s.currentPage = ""
	s.searchQuery = ""
}

func (s *AppState) GroupTags(tagIDs []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create a new tag group
	group := &models.TagGroup{
		ID:       generateID(), // You'll need to implement this
		Name:     "New Group",  // Default name, can be changed later
		TagIDs:   tagIDs,
		Order:    len(s.tagGroups), // Add to end of list
		Expanded: true,
	}

	return s.SaveTagGroup(group)
}

func (s *AppState) ExportTags(tagIDs []string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Create a map for quick lookup of tag IDs to export
	toExport := make(map[string]bool)
	for _, id := range tagIDs {
		toExport[id] = true
	}

	// Collect tags to export
	tagsToExport := make([]models.Tag, 0, len(tagIDs))
	for _, tag := range s.tags {
		if toExport[tag.ID] {
			tagsToExport = append(tagsToExport, tag)
		}
	}

	// Collect tag groups that contain any of the exported tags
	groupsToExport := make([]models.TagGroup, 0)
	for _, group := range s.tagGroups {
		hasExportedTag := false
		for _, tagID := range group.TagIDs {
			if toExport[tagID] {
				hasExportedTag = true
				break
			}
		}
		if hasExportedTag {
			groupsToExport = append(groupsToExport, group)
		}
	}

	// Create export data
	export := &models.TagExport{
		Tags:       tagsToExport,
		TagGroups:  groupsToExport,
		ExportedAt: time.Now(),
		Version:    "1.0",
	}

	// Convert to JSON
	data, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal export data: %w", err)
	}

	// Write to file
	exportPath := fmt.Sprintf("tag_export_%s.json", time.Now().Format("2006-01-02_15-04-05"))
	err = os.WriteFile(exportPath, data, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write export file: %w", err)
	}

	return exportPath, nil
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
