package shortener

import (
	"context"
	"courtsite-url-shortener/internal/domain"
	"sync"
	"time"
)

type Memory struct {
	mu       sync.RWMutex
	url      map[int64]*domain.UrlShortener
	key      map[string]*domain.UrlShortenerKey
	urlIndex map[string]string
	nextID   int64
}

func NewMemoryStore() *Memory {
	return &Memory{
		url:      make(map[int64]*domain.UrlShortener),
		key:      make(map[string]*domain.UrlShortenerKey),
		urlIndex: make(map[string]string),
		nextID:   1,
	}
}
func (m *Memory) SaveAndRecord(ctx context.Context, originalURL string, key string) (domain.UrlShortener, domain.UrlShortenerKey, error) {

	m.mu.Lock()
	defer m.mu.Unlock()

	m.urlIndex[originalURL] = key

	url := &domain.UrlShortener{
		ID:          m.nextID,
		OriginalURL: originalURL,
		Visits:      1,
		CreatedAt:   time.Now(),
	}

	shortKey := &domain.UrlShortenerKey{
		ID:             m.nextID,
		URLShortenerID: url.ID,
		ShortKey:       key,
		CreatedAt:      time.Now(),
		ExpiredAt:      time.Now().Add(5 * time.Second),
	}

	m.url[url.ID] = url
	m.key[key] = shortKey
	m.nextID += 1

	return *url, *shortKey, nil
}

func (m *Memory) GetRecordByShortKey(ctx context.Context, key string) (domain.UrlShortener, error) {

	m.mu.RLock()
	defer m.mu.RUnlock()

	shortKey, ok := m.key[key]
	if !ok {
		return domain.UrlShortener{}, domain.ErrNotFound
	}

	record, ok := m.url[shortKey.URLShortenerID]
	if !ok {
		return domain.UrlShortener{}, domain.ErrNotFound
	}

	return *record, nil
}

func (m *Memory) GetActiveKey(ctx context.Context, originalUrl string) (domain.UrlShortenerKey, bool) {

	m.mu.RLock()
	defer m.mu.RUnlock()

	key, ok := m.urlIndex[originalUrl]
	if !ok {
		return domain.UrlShortenerKey{}, false
	}

	foundKey, ok := m.key[key]
	if !ok {
		return domain.UrlShortenerKey{}, false
	}

	return *foundKey, true

}

func (m *Memory) UpdateRecordWithKey(ctx context.Context, originalUrl string) (domain.UrlShortener, error) {

	m.mu.Lock()
	defer m.mu.Unlock()

	key, ok := m.urlIndex[originalUrl]
	if !ok {
		return domain.UrlShortener{}, domain.ErrNotFound
	}

	foundKey, ok := m.key[key]
	if !ok {
		return domain.UrlShortener{}, domain.ErrNotFound
	}

	record, ok := m.url[foundKey.URLShortenerID]
	record.Visits += 1
	return *record, nil

}
