package shortener

import (
	"context"
	"sync"
	"time"
)

type MemoryStore struct {
	mu      sync.RWMutex
	records map[string]URLShortener
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		records: make(map[string]URLShortener),
	}
}

func (memory *MemoryStore) RecordVisitandSave(ctx context.Context, shortener URLShortener) (URLShortener, error) {
	memory.mu.Lock()
	defer memory.mu.Unlock()

	existing, exists := memory.records[shortener.ShortKey]
	if exists {
		shortener.Visits = existing.Visits + 1
		shortener.CreatedAt = time.Now()
	} else {
		shortener.CreatedAt = existing.CreatedAt
	}

	memory.records[shortener.ShortKey] = shortener
	return shortener, nil

}

func (memory *MemoryStore) FindByShortKey(ctx context.Context, shortKey string) (URLShortener, error) {
	memory.mu.RLock()
	defer memory.mu.RUnlock()

	shortener, exists := memory.records[shortKey]
	if !exists {
		return URLShortener{}, ErrNotFound
	}
	return shortener, nil
}

func (memory *MemoryStore) GetRecordVisits(ctx context.Context, key string) (int64, error) {
	memory.mu.RLock()
	defer memory.mu.RUnlock()

	shortener, exists := memory.records[key]
	if !exists {
		return 0, ErrNotFound
	}
	return shortener.Visits, nil
}
