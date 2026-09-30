package shortener

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_ShortenURL(t *testing.T) {

	tests := []struct {
		name       string
		preset     []string
		inputUrl   string
		wantErr    error
		wantVisits int64
	}{
		{
			name:       "valid url test",
			preset:     nil,
			inputUrl:   "www.google.com",
			wantErr:    nil,
			wantVisits: 1,
		},

		{
			name:       "duplicate url test",
			preset:     []string{"www.google.com"},
			inputUrl:   "www.google.com",
			wantErr:    nil,
			wantVisits: 2,
		},

		{
			name:       "ErrInvalidURL : empty url test",
			preset:     nil,
			inputUrl:   "",
			wantErr:    ErrEmptyURL,
			wantVisits: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			store := NewMemoryStore()
			svc := NewService(store)
			ctx := context.Background()

			for _, p := range tc.preset {
				_, _ = svc.ShortenUrl(ctx, p)
			}

			got, err := svc.ShortenUrl(ctx, tc.inputUrl)

			assert.Equal(t, tc.inputUrl, got.OriginalURL)
			assert.Equal(t, tc.wantVisits, got.Visits)
			assert.NotNil(t, got.ShortKey)

			if err != nil {
				assert.Error(t, tc.wantErr, err)
			}
		})
	}

}

func TestService_GetShortenURL(t *testing.T) {

	tests := []struct {
		name        string
		originalUrl string
		shortKey    bool
		preset      URLShortener
		wantErr     error
	}{
		{
			name:        "valid shortkey test",
			originalUrl: "wwww.google.com",
			shortKey:    true,
			wantErr:     nil,
		},

		{
			name:        "invalid shortkey test",
			originalUrl: "wwww.google.com",
			shortKey:    true,
			wantErr:     ErrNotFound,
		},

		{
			name:        "ErrEmptyShortKey: shortkey not present",
			originalUrl: "wwww.google.com",
			shortKey:    false,
			wantErr:     ErrEmptyShortKey,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			store := NewMemoryStore()
			svc := NewService(store)
			ctx := context.Background()

			saved, err := svc.ShortenUrl(ctx, tc.originalUrl)
			require.NoError(t, err) // Stop test if setup fails

			targetKey := saved.ShortKey

			if !tc.shortKey {
				targetKey = ""
			}

			got, err := svc.GetShortenURL(ctx, targetKey)

			assert.Equal(t, targetKey, got.ShortKey)
			assert.NotNil(t, tc.originalUrl, got.OriginalURL)

			if err != nil {
				assert.Error(t, tc.wantErr, err)
			}

		})
	}
}
