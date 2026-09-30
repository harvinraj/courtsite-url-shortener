package service

import (
	"context"
	"courtsite-url-shortener/internal/domain"
	"crypto/rand"
	"encoding/base64"
	"time"
)

type ShortenerService struct {
	repo domain.ShortenerRepository
}

func NewService(repo domain.ShortenerRepository) *ShortenerService {
	return &ShortenerService{
		repo: repo,
	}
}

func generateKey() string {
	bytes := make([]byte, 6)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)[:6]
}

func (s *ShortenerService) ShortenUrl(ctx context.Context, originalUrl string) (domain.UrlShortener, domain.UrlShortenerKey, error) {

	if originalUrl == "" {
		return domain.UrlShortener{}, domain.UrlShortenerKey{}, domain.ErrEmptyUrl
	}

	key, found := s.repo.GetActiveKey(ctx, originalUrl)
	if found {
		if key.ExpiredAt.After(time.Now()) {

			url, err := s.repo.UpdateRecordWithKey(ctx, originalUrl)
			if err != nil {
				return domain.UrlShortener{}, domain.UrlShortenerKey{}, err
			}
			return url, domain.UrlShortenerKey{
				ShortKey: key.ShortKey,
			}, nil
		}

	}

	newKey := generateKey()

	savedRecord, savedKey, err := s.repo.SaveAndRecord(ctx, originalUrl, newKey)
	if err != nil {
		return domain.UrlShortener{}, domain.UrlShortenerKey{}, err
	}
	return savedRecord, savedKey, nil
}

func (s *ShortenerService) GetShortenUrl(ctx context.Context, shortKey string) (domain.UrlShortener, error) {

	if shortKey == "" {
		return domain.UrlShortener{}, domain.ErrEmptyShortKey
	}

	found, err := s.repo.GetRecordByShortKey(ctx, shortKey)
	if err != nil {
		return domain.UrlShortener{}, err
	}

	return found, nil
}
