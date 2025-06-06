package models

import (
	"time"
)

// Message représente une réponse à un thread (comme une réponse Twitter)
type Message struct {
	ID              int       `json:"id" db:"id_message"`
	ThreadID        int       `json:"thread_id" db:"thread_id"`
	AuthorID        int       `json:"author_id" db:"author_id"`
	Content         string    `json:"content" db:"content"`
	ParentMessageID *int      `json:"parent_message_id" db:"parent_message_id"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	LikeCount       int       `json:"like_count" db:"like_count"`
	DislikeCount    int       `json:"dislike_count" db:"dislike_count"`
	IsEdited        bool      `json:"is_edited" db:"is_edited"`

	// Relations
	Author        *User     `json:"author,omitempty"`
	ParentMessage *Message  `json:"parent_message,omitempty"`
	Replies       []Message `json:"replies,omitempty"`
	UserReaction  *string   `json:"user_reaction,omitempty"`
}

// MessageCreateRequest représente une demande de création de message
type MessageCreateRequest struct {
	ThreadID        int    `json:"thread_id" validate:"required"`
	Content         string `json:"content" validate:"required,min=1,max=2000"`
	ParentMessageID *int   `json:"parent_message_id,omitempty"`
}
