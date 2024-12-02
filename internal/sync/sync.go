package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/goBookMarker/internal/models"
	"golang.org/x/oauth2"
)

type CloudSync interface {
	Upload(data []byte) error
	Download() ([]byte, error)
	LastSync() time.Time
}

type SyncManager struct {
	provider CloudSync
	interval time.Duration
	lastSync time.Time
	running  bool
}

func NewSyncManager(provider CloudSync, intervalMinutes int) *SyncManager {
	return &SyncManager{
		provider: provider,
		interval: time.Duration(intervalMinutes) * time.Minute,
	}
}

func (sm *SyncManager) StartSync(bookmarks []models.Bookmark) {
	if sm.running {
		return
	}

	sm.running = true
	go func() {
		for sm.running {
			if err := sm.sync(bookmarks); err != nil {
				// Handle error (maybe through a channel)
				fmt.Printf("Sync error: %v\n", err)
			}
			time.Sleep(sm.interval)
		}
	}()
}

func (sm *SyncManager) StopSync() {
	sm.running = false
}

func (sm *SyncManager) sync(bookmarks []models.Bookmark) error {
	// Convert bookmarks to JSON
	data, err := json.Marshal(bookmarks)
	if err != nil {
		return fmt.Errorf("failed to marshal bookmarks: %v", err)
	}

	// Upload to cloud
	if err := sm.provider.Upload(data); err != nil {
		return fmt.Errorf("failed to upload: %v", err)
	}

	sm.lastSync = time.Now()
	return nil
}

// Google Drive implementation
type GoogleDriveSync struct {
	client *http.Client
	token  *oauth2.Token
}

func NewGoogleDriveSync(token *oauth2.Token) *GoogleDriveSync {
	return &GoogleDriveSync{
		token:  token,
		client: oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token)),
	}
}

func (s *GoogleDriveSync) Upload(data []byte) error {
	// TODO: Implement Google Drive upload
	return fmt.Errorf("not implemented")
}

func (s *GoogleDriveSync) Download() ([]byte, error) {
	// TODO: Implement Google Drive download
	return nil, fmt.Errorf("not implemented")
}

func (s *GoogleDriveSync) LastSync() time.Time {
	return time.Now() // TODO: Implement proper last sync tracking
}

// OneDrive implementation
type OneDriveSync struct {
	client *http.Client
	token  *oauth2.Token
}

func NewOneDriveSync(token *oauth2.Token) *OneDriveSync {
	return &OneDriveSync{
		token:  token,
		client: oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token)),
	}
}

func (s *OneDriveSync) Upload(data []byte) error {
	// TODO: Implement OneDrive upload
	return fmt.Errorf("not implemented")
}

func (s *OneDriveSync) Download() ([]byte, error) {
	// TODO: Implement OneDrive download
	return nil, fmt.Errorf("not implemented")
}

func (s *OneDriveSync) LastSync() time.Time {
	return time.Now() // TODO: Implement proper last sync tracking
}
