package models

import "time"

type URL struct {
    ID          int        `json:"id"`
    ShortCode   string     `json:"short_code"`
    OriginalURL string     `json:"original_url"`
    UserID      *int       `json:"user_id,omitempty"` 
    Clicks      int64      `json:"clicks"`
    CreatedAt   time.Time  `json:"created_at"`
    ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

type ShortenRequest struct {
    URL          string `json:"url" binding:"required,url"`
    CustomCode   string `json:"custom_code,omitempty"`
    ExpiresInHrs int    `json:"expires_in_hrs,omitempty"`
}

type ShortenResponse struct {
    ShortCode   string     `json:"short_code"`
    ShortURL    string     `json:"short_url"`
    OriginalURL string     `json:"original_url"`
    ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

type URLStatsResponse struct {
    URL
    ShortURL string `json:"short_url"`
}
