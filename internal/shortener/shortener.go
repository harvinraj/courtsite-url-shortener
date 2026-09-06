package shortener

import "time"

type URLShortener struct {
	ShortKey    string    `json:"short_key"`
	OriginalURL string    `json:"original_url"`
	Visits      int64     `json:"visits"`
	CreatedAt   time.Time `json:"created_at"`
}
