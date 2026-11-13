package models

import "time"

type CreateURLRequest struct {
	URL string `json:"url"`
}

type URL struct {
	ID          int
	OriginalURL string
	ShortCode   string
	CreatedAt   time.Time
	ExpiresAt	time.Time
}


type CreateURLResponse struct{
    ShortURL    string    `json:"short_url"`
    OriginalURL string    `json:"original_url"`
    ExpiresAt   time.Time `json:"expires_at"`
}

