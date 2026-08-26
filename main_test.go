package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortenAndRedirectFlow(t *testing.T) {
	store := NewMemoryStore()
	router := SetupRouter(store)

	// 1. Test POST /shorten_url/ with valid accessible URL
	body := map[string]string{"url": "https://www.google.com"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/shorten_url/", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp ShortenResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Key)

	// 2. Test GET /shorten_url/<key> Redirection
	wRedir := httptest.NewRecorder()
	reqRedir, _ := http.NewRequest("GET", "/shorten_url/"+resp.Key, nil)
	router.ServeHTTP(wRedir, reqRedir)

	assert.Equal(t, http.StatusFound, wRedir.Code)
	assert.Equal(t, "https://www.google.com", wRedir.Header().Get("Location"))

	// 3. Test GET /analytics/<key> (Click count verification)
	wAnalytics := httptest.NewRecorder()
	reqAnalytics, _ := http.NewRequest("GET", "/analytics/"+resp.Key, nil)
	router.ServeHTTP(wAnalytics, reqAnalytics)

	assert.Equal(t, http.StatusOK, wAnalytics.Code)
	var analyticsRecord URLRecord
	_ = json.Unmarshal(wAnalytics.Body.Bytes(), &analyticsRecord)
	assert.Equal(t, int64(1), analyticsRecord.ClickCount)
}

func TestInvalidAndInaccessibleURLs(t *testing.T) {
	store := NewMemoryStore()
	router := SetupRouter(store)

	// Inaccessible / invalid domain
	body := map[string]string{"url": "https://invalid-non-existent-domain-999.com"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/shorten_url/", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNotFoundKey(t *testing.T) {
	store := NewMemoryStore()
	router := SetupRouter(store)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/shorten_url/nonexistentkey", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
