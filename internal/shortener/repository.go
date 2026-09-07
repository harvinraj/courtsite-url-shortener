package shortener

import (
	"context"
	"errors"
	"net/http"
)

var (
	ErrNotFound = errors.New("record not found")
)

// Repository defines the data access contract for URL records.
type Repository interface {
	FindByShortKey(ctx context.Context, shortKey string) (URLShortener, error)
	RecordVisitandSave(ctx context.Context, shortener URLShortener) (URLShortener, error)
	GetRecordVisits(ctx context.Context, key string) (int64, error)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
