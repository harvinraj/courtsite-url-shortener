package domain

import (
	"context"
	"errors"
	"time"
)

var ErrEmptyUrl = errors.New("url value is empty")
var ErrEmptyShortKey = errors.New("short_key value is empty")
var ErrNotFound = errors.New("record not found")

type UrlShortener struct {
	ID          int64     `json:"id"`
	OriginalURL string    `json:"original_id"`
	Visits      int64     `json:"visits"`
	CreatedAt   time.Time `json:"created_at"`
}

type UrlShortenerKey struct {
	ID             int64     `json:"id"`
	URLShortenerID int64     `json:"url_shortener_id"`
	ShortKey       string    `json:"short_key"`
	CreatedAt      time.Time `json:"created_at"`
	ExpiredAt      time.Time `json:"expired_at"`
}

type ShortenerRepository interface {
	SaveAndRecord(ctx context.Context, originalUrl string, key string) (UrlShortener, UrlShortenerKey, error)
	GetRecordByShortKey(ctx context.Context, key string) (UrlShortener, error)
	GetActiveKey(ctx context.Context, originalUrl string) (UrlShortenerKey, bool)
	UpdateRecordWithKey(ctx context.Context, originalUrl string) (UrlShortener, error)
}
