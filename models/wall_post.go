package models

import "time"

// WallPost représente une publication sur le mur d'un utilisateur
type WallPost struct {
	ID         int       `json:"id" db:"id"`
	UserID     int       `json:"user_id" db:"user_id"`
	AuthorID   int       `json:"author_id" db:"author_id"`
	Content    string    `json:"content" db:"content"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// CreateWallPostRequest représente la requête de création d'un post sur le mur
type CreateWallPostRequest struct {
	UserID  int    `json:"user_id" validate:"required"`
	Content string `json:"content" validate:"required,min=1,max=1000"`
}

// WallPostWithAuthor représente un post du mur avec les informations de l'auteur
type WallPostWithAuthor struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	AuthorID    int       `json:"author_id" db:"author_id"`
	Content     string    `json:"content" db:"content"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	AuthorName  string    `json:"author_name" db:"author_name"`
	AuthorEmail string    `json:"author_email" db:"author_email"`
	AvatarPath  string    `json:"avatar_path" db:"avatar_path"`
} 