package main

import (
	"crypto/rand"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// URLRecord stores the original URL and basic business analytics.
type URLRecord struct {
	OriginalURL string    `json:"original_url"`
	ShortKey    string    `json:"short_key"`
	CreatedAt   time.Time `json:"created_at"`
	ClickCount  int64     `json:"click_count"`
	LastVisited time.Time `json:"last_visited,omitempty"`
}

// MemoryStore provides a thread-safe in-memory key-value store.
type MemoryStore struct {
	mu    sync.RWMutex
	store map[string]*URLRecord
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		store: make(map[string]*URLRecord),
	}
}

func (s *MemoryStore) Save(key, originalURL string) *URLRecord {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := &URLRecord{
		OriginalURL: originalURL,
		ShortKey:    key,
		CreatedAt:   time.Now(),
		ClickCount:  0,
	}
	s.store[key] = record
	return record
}

func (s *MemoryStore) GetAndRecordVisit(key string) (*URLRecord, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, exists := s.store[key]
	if !exists {
		return nil, false
	}

	record.ClickCount++
	record.LastVisited = time.Now()
	return record, true
}

func (s *MemoryStore) GetAnalytics(key string) (*URLRecord, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	record, exists := s.store[key]
	if !exists {
		return nil, false
	}
	return record, true
}

// GenerateBase62Key creates a clean, 7-character random key.
func GenerateBase62Key(length int) (string, error) {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[num.Int64()]
	}
	return string(b), nil
}

// ValidateAndCheckReachability validates URL syntax and executes a lightweight HEAD check.
func ValidateAndCheckReachability(rawURL string) bool {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return false
	}

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	// Helper to build a request with a realistic User-Agent header
	createRequest := func(method string) (*http.Request, error) {
		req, err := http.NewRequest(method, rawURL, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", "CourtsiteURLChecker/1.0")
		return req, nil
	}

	// 1. Try HEAD request
	if req, err := createRequest("HEAD"); err == nil {
		if resp, err := client.Do(req); err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode < 400 {
				return true
			}
		}
	}

	// 2. Fallback to GET request if HEAD is blocked/rejected
	if req, err := createRequest("GET"); err == nil {
		if resp, err := client.Do(req); err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode < 400 {
				return true
			}
		}
	}

	return false
}

type ShortenRequest struct {
	URL string `json:"url" binding:"required"`
}

type ShortenResponse struct {
	Key      string `json:"key"`
	ShortURL string `json:"short_url,omitempty"`
}

func SetupRouter(store *MemoryStore) *gin.Engine {
	r := gin.Default()

	// POST /shorten_url/
	r.POST("/shorten_url/", func(c *gin.Context) {
		var req ShortenRequest
		if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.URL) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload. 'url' is required."})
			return
		}

		// Validate URL format and accessibility
		if !ValidateAndCheckReachability(req.URL) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Provided URL is invalid or inaccessible."})
			return
		}

		key, err := GenerateBase62Key(7)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate short key."})
			return
		}

		store.Save(key, req.URL)

		c.JSON(http.StatusOK, ShortenResponse{
			Key: key,
		})
	})

	// GET /shorten_url/<key>
	r.GET("/shorten_url/:key", func(c *gin.Context) {
		key := c.Param("key")

		record, exists := store.GetAndRecordVisit(key)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Short key not found."})
			return
		}

		// 302 Found redirect to original URL
		c.Redirect(http.StatusFound, record.OriginalURL)
	})

	// GET /analytics/<key> (Stretch Goal: BI tracking)
	r.GET("/analytics/:key", func(c *gin.Context) {
		key := c.Param("key")

		record, exists := store.GetAnalytics(key)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Short key not found."})
			return
		}

		c.JSON(http.StatusOK, record)
	})

	return r
}

func main() {
	store := NewMemoryStore()
	r := SetupRouter(store)
	_ = r.Run(":8080")
}
