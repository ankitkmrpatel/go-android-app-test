package share

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/goBookMarker/internal/models"
)

type ShareHandler struct {
	// Channel for receiving shared content
	SharedContent chan *SharedItem
}

type SharedItem struct {
	Type        string // url, image
	Content     string // URL or base64 image data
	Title       string
	Description string
}

func NewShareHandler() *ShareHandler {
	return &ShareHandler{
		SharedContent: make(chan *SharedItem, 10),
	}
}

func (h *ShareHandler) HandleSharedContent(contentType, content string) error {
	item := &SharedItem{}

	// Determine content type and process accordingly
	if isURL(content) {
		item.Type = "url"
		if err := h.processURL(item, content); err != nil {
			return fmt.Errorf("failed to process URL: %v", err)
		}
	} else if isImage(content) {
		item.Type = "image"
		if err := h.processImage(item, content); err != nil {
			return fmt.Errorf("failed to process image: %v", err)
		}
	} else {
		return fmt.Errorf("unsupported content type")
	}

	// Send to channel for processing
	h.SharedContent <- item
	return nil
}

func (h *ShareHandler) processURL(item *SharedItem, urlStr string) error {
	item.Content = urlStr

	// Fetch metadata
	resp, err := http.Get(urlStr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the first 1MB of the response
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		return err
	}

	// Extract title and description from HTML
	bodyStr := string(body)
	item.Title = extractMetaTag(bodyStr, "title")
	item.Description = extractMetaTag(bodyStr, "description")

	return nil
}

func (h *ShareHandler) processImage(item *SharedItem, imageData string) error {
	item.Content = imageData
	item.Title = "Shared Image"
	item.Description = fmt.Sprintf("Image shared on %s", time.Now().Format("Jan 2, 2006"))
	return nil
}

func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https")
}

func isImage(str string) bool {
	return strings.HasPrefix(str, "data:image/") ||
		strings.HasSuffix(strings.ToLower(str), ".jpg") ||
		strings.HasSuffix(strings.ToLower(str), ".jpeg") ||
		strings.HasSuffix(strings.ToLower(str), ".png") ||
		strings.HasSuffix(strings.ToLower(str), ".gif")
}

func extractMetaTag(html, tag string) string {
	tagStart := strings.Index(html, fmt.Sprintf(`<meta name="%s"`, tag))
	if tagStart == -1 {
		tagStart = strings.Index(html, fmt.Sprintf(`<meta property="%s"`, tag))
	}
	if tagStart == -1 {
		return ""
	}

	contentStart := strings.Index(html[tagStart:], `content="`)
	if contentStart == -1 {
		return ""
	}
	contentStart += tagStart + 9

	contentEnd := strings.Index(html[contentStart:], `"`)
	if contentEnd == -1 {
		return ""
	}
	contentEnd += contentStart

	return html[contentStart:contentEnd]
}

// Convert SharedItem to Bookmark
func (item *SharedItem) ToBookmark() *models.Bookmark {
	return &models.Bookmark{
		URL:         item.Content,
		Title:       item.Title,
		Description: item.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
