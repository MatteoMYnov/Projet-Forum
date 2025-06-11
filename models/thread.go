package models

import (
	"time"
)

// Thread représente un fil de discussion (équivalent à un Tweet)
type Thread struct {
	ID            int         `json:"id" db:"id_thread"`
	Title         string      `json:"title" db:"title"`
	Content       string      `json:"content" db:"content"`
	AuthorID      int         `json:"author_id" db:"author_id"`
	CategoryID    *int        `json:"category_id" db:"category_id"`
	Status        string      `json:"status" db:"status"` // "open", "closed", "archived"
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" db:"updated_at"`
	IsPinned      bool        `json:"is_pinned" db:"is_pinned"`
	ViewCount     int         `json:"view_count" db:"view_count"`
	LikeCount     int         `json:"like_count" db:"like_count"`
	DislikeCount  int         `json:"dislike_count" db:"dislike_count"`
	LoveCount     int         `json:"love_count" db:"love_count"`
	MessageCount  int         `json:"message_count" db:"message_count"`
	LastActivity  time.Time   `json:"last_activity" db:"last_activity"`
	
	// Relations (chargées si nécessaire)
	Author       *User       `json:"author,omitempty"`
	Category     *Category   `json:"category,omitempty"`
	Messages     []Message   `json:"messages,omitempty"`
	Hashtags     []Hashtag   `json:"hashtags,omitempty"`
	UserReaction *string     `json:"user_reaction,omitempty"` // "like", "dislike", null
}

// ThreadCreateRequest représente une demande de création de thread
type ThreadCreateRequest struct {
	Title      string   `json:"title" validate:"required,min=1,max=280"`
	Content    string   `json:"content" validate:"required,min=1,max=5000"`
	CategoryID *int     `json:"category_id"`
	Hashtags   []string `json:"hashtags,omitempty"`
}
