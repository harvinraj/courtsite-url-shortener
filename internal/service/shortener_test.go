package service

import (
	"context"
	"courtsite-url-shortener/internal/domain"
	"errors"
	"testing"
	"time"
)

type MockRepository struct {
	SaveAndRecordFunc       func(ctx context.Context, originalUrl string, key string) (domain.UrlShortener, domain.UrlShortenerKey, error)
	GetRecordByShortKeyFunc func(ctx context.Context, key string) (domain.UrlShortener, error)
	GetActiveKeyFunc        func(ctx context.Context, originalUrl string) (domain.UrlShortenerKey, bool)
	UpdateRecordWithKeyFunc func(ctx context.Context, originalUrl string) (domain.UrlShortener, error)
}

func (m *MockRepository) SaveAndRecord(ctx context.Context, originalUrl string, key string) (domain.UrlShortener, domain.UrlShortenerKey, error) {
	if m.SaveAndRecordFunc != nil {
		return m.SaveAndRecordFunc(ctx, originalUrl, key)
	}
	return domain.UrlShortener{}, domain.UrlShortenerKey{}, nil

}

func (m *MockRepository) GetRecordByShortKey(ctx context.Context, key string) (domain.UrlShortener, error) {
	if m.GetRecordByShortKeyFunc != nil {
		return m.GetRecordByShortKeyFunc(ctx, key)
	}

	return domain.UrlShortener{}, nil
}

func (m *MockRepository) GetActiveKey(ctx context.Context, originalUrl string) (domain.UrlShortenerKey, bool) {
	if m.GetActiveKeyFunc != nil {
		return m.GetActiveKeyFunc(ctx, originalUrl)
	}

	return domain.UrlShortenerKey{}, false
}

func (m *MockRepository) UpdateRecordWithKey(ctx context.Context, originalUrl string) (domain.UrlShortener, error) {
	if m.UpdateRecordWithKeyFunc != nil {
		return m.UpdateRecordWithKeyFunc(ctx, originalUrl)
	}

	return domain.UrlShortener{}, nil
}

func TestShortenUrl(t *testing.T) {

	ctx := context.Background()
	testUrl := "www.google.com"

	tests := []struct {
		name         string
		originalUrl  string
		setupMock    func() *MockRepository
		wantKeyEmpty bool
		wantRecordID int64
		wantErr      error
	}{
		{
			name:        "test empty url",
			originalUrl: "",
			setupMock: func() *MockRepository {
				return &MockRepository{}
			},
			wantErr: domain.ErrEmptyUrl,
		},
		{
			name:        "found - active key NOT EXPIRED, reusing key",
			originalUrl: testUrl,
			setupMock: func() *MockRepository {
				return &MockRepository{
					GetActiveKeyFunc: func(ctx context.Context, originalUrl string) (domain.UrlShortenerKey, bool) {
						return domain.UrlShortenerKey{
							ShortKey:  "activekey",
							ExpiredAt: time.Now().Add(10 * time.Minute),
						}, true
					},

					UpdateRecordWithKeyFunc: func(ctx context.Context, originalUrl string) (domain.UrlShortener, error) {
						return domain.UrlShortener{
							ID:          1,
							OriginalURL: testUrl,
							Visits:      2,
						}, nil
					},
				}
			},
			wantErr:      nil,
			wantKeyEmpty: true,
			wantRecordID: 1,
		},
		{
			name:        "found - active key EXPIRED, create NEW key",
			originalUrl: testUrl,
			setupMock: func() *MockRepository {
				return &MockRepository{
					GetActiveKeyFunc: func(ctx context.Context, originalUrl string) (domain.UrlShortenerKey, bool) {
						return domain.UrlShortenerKey{
							ShortKey:  "expiredKey",
							ExpiredAt: time.Now().Add(-10 * time.Minute),
						}, true
					},

					SaveAndRecordFunc: func(ctx context.Context, originalUrl, key string) (domain.UrlShortener, domain.UrlShortenerKey, error) {
						return domain.UrlShortener{
								ID:          1,
								OriginalURL: testUrl,
							},

							domain.UrlShortenerKey{
								ID:             2,
								URLShortenerID: 1,
								ShortKey:       key,
							}, nil
					},
				}
			},
			wantErr:      nil,
			wantKeyEmpty: false,
			wantRecordID: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			mockRepo := tc.setupMock()
			svc := NewService(mockRepo)

			record, key, err := svc.ShortenUrl(ctx, tc.originalUrl)

			// 1. Check for expected errors
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("expected error: %v, got: %v", tc.wantErr, err)
				}
				return // Expected error occurred, test case passes, move to next
			}

			// 2. Check for unexpected errors
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// 3. Assert values (using Errorf so it checks all fields even if one fails)
			if record.ID != tc.wantRecordID {
				t.Errorf("expected Record ID %d, got %d", tc.wantRecordID, record.ID)
			}

			if tc.wantKeyEmpty && key.ShortKey != "" {
				t.Errorf("expected empty ShortKey, got %q", key.ShortKey)
			}

			if !tc.wantKeyEmpty && key.ShortKey == "" {
				t.Errorf("expected a generated ShortKey, got empty string")
			}

		})
	}
}
