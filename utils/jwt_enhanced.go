package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"forum/models"
	"time"
)

// RefreshToken structure pour les tokens de renouvellement
type RefreshToken struct {
	Token     string    `json:"token"`
	UserID    int       `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// GenerateRefreshToken génère un token de renouvellement
func GenerateRefreshToken(userID int) (*RefreshToken, error) {
	// Générer un token aléatoirement sécurisé
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, fmt.Errorf("erreur génération token: %v", err)
	}

	token := hex.EncodeToString(tokenBytes)

	return &RefreshToken{
		Token:     token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 jours
		CreatedAt: time.Now(),
	}, nil
}

// GenerateTokenPair génère une paire access/refresh token
func GenerateTokenPair(userID int, username, role string) (accessToken, refreshToken string, err error) {
	// Token d'accès (courte durée - 15 minutes)
	accessToken, err = GenerateJWTWithDuration(userID, username, role, 15*time.Minute)
	if err != nil {
		return "", "", fmt.Errorf("erreur génération access token: %v", err)
	}

	// Token de renouvellement (longue durée - 7 jours)
	refreshTokenObj, err := GenerateRefreshToken(userID)
	if err != nil {
		return "", "", fmt.Errorf("erreur génération refresh token: %v", err)
	}

	return accessToken, refreshTokenObj.Token, nil
}

// GenerateJWTWithDuration génère un JWT avec une durée personnalisée
func GenerateJWTWithDuration(userID int, username, role string, duration time.Duration) (string, error) {
	now := time.Now()
	payload := models.SessionInfo{
		UserID:    userID,
		Username:  username,
		Role:      role,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(duration).Unix(),
	}

	return generateJWTFromPayload(payload)
}

// ValidateAndRefreshToken valide un token et le renouvelle si nécessaire
func ValidateAndRefreshToken(token string) (*models.SessionInfo, string, error) {
	sessionInfo, err := ValidateJWT(token)
	if err == nil {
		// Token valide, vérifier s'il doit être renouvelé (si expire dans moins de 5 min)
		timeToExpire := time.Unix(sessionInfo.ExpiresAt, 0).Sub(time.Now())
		if timeToExpire < 5*time.Minute {
			// Renouveler le token
			newToken, err := GenerateJWT(sessionInfo.UserID, sessionInfo.Username, sessionInfo.Role)
			if err != nil {
				return sessionInfo, "", nil // Retourner l'ancien token en cas d'erreur
			}
			return sessionInfo, newToken, nil
		}
		return sessionInfo, "", nil
	}

	if err == ErrTokenExpired {
		// Token expiré, impossible de renouveler automatiquement
		return nil, "", ErrTokenExpired
	}

	return nil, "", err
}

// BlacklistToken ajoute un token à la liste noire (pour déconnexion sécurisée)
var tokenBlacklist = make(map[string]time.Time)

func BlacklistToken(token string, expiry time.Time) {
	tokenBlacklist[token] = expiry
	
	// Nettoyage périodique des tokens expirés
	go cleanupBlacklist()
}

func IsTokenBlacklisted(token string) bool {
	expiry, exists := tokenBlacklist[token]
	if !exists {
		return false
	}
	
	// Si le token a expiré, le retirer de la blacklist
	if time.Now().After(expiry) {
		delete(tokenBlacklist, token)
		return false
	}
	
	return true
}

func cleanupBlacklist() {
	now := time.Now()
	for token, expiry := range tokenBlacklist {
		if now.After(expiry) {
			delete(tokenBlacklist, token)
		}
	}
}

// Fonctions helper internes
func generateJWTFromPayload(payload models.SessionInfo) (string, error) {
	// Réutiliser la logique existante de GenerateJWT
	return GenerateJWT(payload.UserID, payload.Username, payload.Role)
} 