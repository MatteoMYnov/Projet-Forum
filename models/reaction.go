package models

import (
	"time"
)

// Reaction représente une réaction (like/dislike/repost)
type Reaction struct {
	ID           int       `json:"id" db:"id_reaction"`
	UserID       int       `json:"user_id" db:"user_id"`
	ThreadID     *int      `json:"thread_id" db:"thread_id"`
	MessageID    *int      `json:"message_id" db:"message_id"`
	ReactionType string    `json:"reaction_type" db:"reaction_type"` // "like", "dislike", "love", "repost"
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// ReactionRequest représente une demande de réaction
type ReactionRequest struct {
	TargetType   string `json:"target_type" validate:"required,oneof=thread message"`
	TargetID     int    `json:"target_id" validate:"required"`
	ReactionType string `json:"reaction_type" validate:"required,oneof=like dislike love repost"`
}