package models

import (
	"time"
)

// Hashtag représente un hashtag (comme sur Twitter)
type Hashtag struct {
	ID            int       `json:"id" db:"id_hashtag"`
	Name          string    `json:"name" db:"name"`
	UsageCount    int       `json:"usage_count" db:"usage_count"`
	TrendingScore float64   `json:"trending_score" db:"trending_score"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	LastUsed      time.Time `json:"last_used" db:"last_used"`
	IsTrending    bool      `json:"is_trending" db:"is_trending"`
}