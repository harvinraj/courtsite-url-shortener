package shortener

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
)

var ErrInvalidURL = errors.New("invalid url")

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func generateKey() string {
	bytes := make([]byte, 6)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)[:6]
}

func (s *Service) ShortenUrl(ctx context.Context, originalURL string) (URLShortener, error) {

	if len(originalURL) == 0 {
		return URLShortener{}, ErrInvalidURL
	}

	shortKey := generateKey()

	shorten := &URLShortener{
		OriginalURL: originalURL,
		ShortKey:    shortKey,
		Visits:      1,
	}

	saved, err := s.repo.SaveAndRecord(ctx, shorten)
	if err != nil {
		fmt.Println("failed to save url : " + originalURL)
		return URLShortener{}, err
	}

	return saved, nil

}

func (s *Service) GetShortenURL(ctx context.Context, shortKey string) (URLShortener, error) {

	found, err := s.repo.GetRecord(ctx, shortKey)
	if err != nil {
		fmt.Println("failed to retrieve url : " + shortKey)
		return URLShortener{}, err
	}

	return found, nil

}
