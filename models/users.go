package models

import (
	"time"
)

// User représente un utilisateur du forum
type User struct {
	ID             int        `json:"id" db:"id_user"`
	Username       string     `json:"username" db:"username"`
	Email          string     `json:"email" db:"email"`
	PasswordHash   string     `json:"-" db:"password_hash"` // Ne pas exposer en JSON
	ProfilePicture *string    `json:"profile_picture" db:"profile_picture"`
	Banner         *string    `json:"banner" db:"banner"`
	Bio            *string    `json:"bio" db:"bio"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	LastLogin      *time.Time `json:"last_login" db:"last_login"`
	IsVerified     bool       `json:"is_verified" db:"is_verified"`
	IsBanned       bool       `json:"is_banned" db:"is_banned"`
	Role           string     `json:"role" db:"role"` // "user" ou "admin"
	FollowerCount  int        `json:"follower_count" db:"follower_count"`
	FollowingCount int        `json:"following_count" db:"following_count"`
	ThreadCount    int        `json:"thread_count" db:"thread_count"`
}

// UserProfile représente le profil public d'un utilisateur
type UserProfile struct {
	User
	IsFollowing   bool     `json:"is_following"`
	IsFollowedBy  bool     `json:"is_followed_by"`
	RecentThreads []Thread `json:"recent_threads,omitempty"`
}
