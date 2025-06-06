package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"forum/config"
	"forum/models"
	"strings"
	"time"
)

var (
	ErrInvalidToken = errors.New("token invalide")
	ErrTokenExpired = errors.New("token expiré")
)

// JWTHeader représente l'en-tête du token JWT
type JWTHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

// GenerateJWT génère un token JWT pour un utilisateur
func GenerateJWT(userID int, username, role string) (string, error) {
	// En-tête
	header := JWTHeader{
		Alg: "HS256",
		Typ: "JWT",
	}

	// Payload (claims)
	now := time.Now()
	payload := models.SessionInfo{
		UserID:    userID,
		Username:  username,
		Role:      role,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(24 * time.Hour).Unix(), // Expire dans 24h
	}

	// Encoder l'en-tête et le payload en Base64URL
	headerJSON, _ := json.Marshal(header)
	payloadJSON, _ := json.Marshal(payload)

	headerB64 := base64URLEncode(headerJSON)
	payloadB64 := base64URLEncode(payloadJSON)

	// Créer la signature
	message := headerB64 + "." + payloadB64
	signature := sign(message, config.JWTSecret)
	signatureB64 := base64URLEncode(signature)

	// Token final
	token := message + "." + signatureB64
	return token, nil
}

// ValidateJWT valide et décode un token JWT
func ValidateJWT(token string) (*models.SessionInfo, error) {
	if token == "" {
		return nil, ErrInvalidToken
	}

	// Diviser le token en parties
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	headerB64, payloadB64, signatureB64 := parts[0], parts[1], parts[2]

	// Vérifier la signature
	message := headerB64 + "." + payloadB64
	expectedSignature := sign(message, config.JWTSecret)
	expectedSignatureB64 := base64URLEncode(expectedSignature)

	if signatureB64 != expectedSignatureB64 {
		return nil, ErrInvalidToken
	}

	// Décoder le payload
	payloadBytes, err := base64URLDecode(payloadB64)
	if err != nil {
		return nil, ErrInvalidToken
	}

	var sessionInfo models.SessionInfo
	if err := json.Unmarshal(payloadBytes, &sessionInfo); err != nil {
		return nil, ErrInvalidToken
	}

	// Vérifier l'expiration
	if time.Now().Unix() > sessionInfo.ExpiresAt {
		return nil, ErrTokenExpired
	}

	return &sessionInfo, nil
}

// ExtractTokenFromCookie extrait le token depuis un cookie
func ExtractTokenFromCookie(cookieValue string) string {
	// Si le cookie contient "Bearer token", extraire le token
	if strings.HasPrefix(cookieValue, "Bearer ") {
		return strings.TrimPrefix(cookieValue, "Bearer ")
	}
	return cookieValue
}

// sign crée une signature HMAC-SHA256
func sign(message, secret string) []byte {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return h.Sum(nil)
}

// base64URLEncode encode en Base64URL (sans padding)
func base64URLEncode(data []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(data), "=")
}

// base64URLDecode décode depuis Base64URL
func base64URLDecode(data string) ([]byte, error) {
	// Ajouter le padding si nécessaire
	switch len(data) % 4 {
	case 2:
		data += "=="
	case 3:
		data += "="
	}
	return base64.URLEncoding.DecodeString(data)
}

// CreateSecureCookie crée un cookie sécurisé pour le token
func CreateSecureCookie(name, value string, maxAge int) string {
	return fmt.Sprintf("%s=%s; HttpOnly; Secure; SameSite=Strict; Max-Age=%d; Path=/",
		name, value, maxAge)
}
