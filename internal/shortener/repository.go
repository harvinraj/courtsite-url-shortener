package shortener

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("record not found")
)

// Repository defines the data access contract for URL records.
type Repository interface {
	SaveAndRecord(ctx context.Context, shortener *URLShortener) (URLShortener, error)
	GetRecord(ctx context.Context, shortKey string) (URLShortener, error)
}
