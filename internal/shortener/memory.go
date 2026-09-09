package shortener

import (
	"context"
	"sync"
	"time"
)

type Memory struct {
	mu       sync.RWMutex
	record   map[string]*URLShortener
	urlIndex map[string]string
}

func NewMemoryStore() *Memory {
	return &Memory{
		record:   make(map[string]*URLShortener),
		urlIndex: make(map[string]string),
	}
}
func (m *Memory) SaveAndRecord(ctx context.Context, shortener *URLShortener) (URLShortener, error) {

	m.mu.Lock()
	defer m.mu.Unlock()

	if existingKey, found := m.urlIndex[shortener.OriginalURL]; found {
		existing := m.record[existingKey]
		existing.Visits++
		m.record[existingKey] = existing
		return *existing, nil
	}

	shortener.CreatedAt = time.Now()
	m.record[shortener.ShortKey] = shortener
	m.urlIndex[shortener.OriginalURL] = shortener.ShortKey
	return *shortener, nil
}

func (m *Memory) GetRecord(ctx context.Context, shortKey string) (URLShortener, error) {

	m.mu.RLock()
	defer m.mu.RUnlock()

	found, ok := m.record[shortKey]
	if !ok {
		return URLShortener{}, ErrNotFound
	}

	return *found, nil
}
