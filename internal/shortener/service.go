package shortener

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

type Service struct {
	repo       Repository
	httpClient HTTPClient
}

func NewService(repo Repository, client HTTPClient) *Service {
	return &Service{
		repo:       repo,
		httpClient: client,
	}
}

// GenerateBase62Key generates a random short key for the URL
func generateKey() string {
	bytes := make([]byte, 6)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)[:6]
}

func (s *Service) CheckReachability(ctx context.Context, targetURL string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, targetURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "CourtsiteURLShortener/2.0")

	resp, err := s.httpClient.Do(req)
	if err == nil {
		resp.Body.Close()
		if resp.StatusCode < 500 {
			return nil
		}
	}

	// Fallback to GET if HEAD fails or returns a server error
	reqGet, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return err
	}
	reqGet.Header.Set("User-Agent", "CourtsiteURLShortener/1.0")

	respGet, err := s.httpClient.Do(reqGet)
	if err != nil {
		return fmt.Errorf("reachability check failed via GET: %w", err)
	}
	defer respGet.Body.Close()

	if respGet.StatusCode >= 400 {
		return fmt.Errorf("target url returned status code: %d", respGet.StatusCode)
	}

	return nil
}

func (s *Service) Shorten(ctx context.Context, originalURL string) (URLShortener, error) {

	key := generateKey()
	record := &URLShortener{
		ShortKey:    key,
		OriginalURL: originalURL,
		Visits:      0,
		CreatedAt:   time.Now(),
	}

	saved, err := s.repo.RecordVisitandSave(ctx, *record)
	if err != nil {
		fmt.Errorf("failed to save shortener: %s : (key = %s) %w", originalURL, key, err)
		return URLShortener{}, err
	}

	return saved, nil
}

func (s *Service) Retrieve(ctx context.Context, key string) (int64, error) {

	retrieved, err := s.repo.GetRecordVisits(ctx, key)
	if err != nil {
		fmt.Errorf("no short key found: %s : ", key)
		return 0, nil
	}

	return retrieved, nil
}
