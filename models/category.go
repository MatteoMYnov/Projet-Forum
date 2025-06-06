package models

import (
	"time"
)

// Category représente une catégorie/tag (comme les hashtags Twitter)
type Category struct {
	ID          int       `json:"id" db:"id_category"`
	Name        string    `json:"name" db:"name"`
	Color       string    `json:"color" db:"color"`
	Description *string   `json:"description" db:"description"`
	ThreadCount int       `json:"thread_count" db:"thread_count"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	IsActive    bool      `json:"is_active" db:"is_active"`
}