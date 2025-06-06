package models

import (
	"time"
)

// Follow représente un abonnement (comme sur Twitter)
type Follow struct {
	ID          int       `json:"id" db:"id_follow"`
	FollowerID  int       `json:"follower_id" db:"follower_id"`
	FollowingID int       `json:"following_id" db:"following_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	
	// Relations
	Follower  *User `json:"follower,omitempty"`
	Following *User `json:"following,omitempty"`
}