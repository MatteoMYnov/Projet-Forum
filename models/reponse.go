package models

// APIResponse représente une réponse API standardisée
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta contient des métadonnées pour la pagination
type Meta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalPages int `json:"total_pages"`
	TotalCount int `json:"total_count"`
}

// PaginationRequest représente une demande de pagination
type PaginationRequest struct {
	Page    int    `json:"page" validate:"min=1"`
	PerPage int    `json:"per_page" validate:"min=1,max=100"`
	SortBy  string `json:"sort_by,omitempty"`
	Order   string `json:"order,omitempty" validate:"omitempty,oneof=asc desc"`
}

// LoginRequest représente une demande de connexion
type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"` // username ou email
	Password   string `json:"password" validate:"required,min=8"`
}

// RegisterRequest représente une demande d'inscription
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=12"`
}

// SessionInfo contient les informations de session utilisateur
type SessionInfo struct {
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	IssuedAt  int64  `json:"issued_at"`
	ExpiresAt int64  `json:"expires_at"`
}
