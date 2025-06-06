package middleware

import (
	"context"
	"forum/models"
	"forum/utils"
	"log"
	"net/http"
)

// ContextKey type pour les clés de contexte
type ContextKey string

const UserContextKey ContextKey = "user"

// RequireAuth middleware qui vérifie l'authentification
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Récupérer le token depuis le cookie
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			log.Printf("🔒 Accès non autorisé - pas de cookie: %s", r.URL.Path)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// Valider le token
		token := utils.ExtractTokenFromCookie(cookie.Value)
		sessionInfo, err := utils.ValidateJWT(token)
		if err != nil {
			log.Printf("🔒 Token invalide: %v", err)
			// Supprimer le cookie invalide
			http.SetCookie(w, &http.Cookie{
				Name:   "auth_token",
				Value:  "",
				MaxAge: -1,
				Path:   "/",
			})
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// Ajouter les infos utilisateur au contexte
		ctx := context.WithValue(r.Context(), UserContextKey, sessionInfo)
		next(w, r.WithContext(ctx))
	}
}

// GetUserFromContext récupère l'utilisateur depuis le contexte
func GetUserFromContext(r *http.Request) *models.SessionInfo {
	user, ok := r.Context().Value(UserContextKey).(*models.SessionInfo)
	if !ok {
		return nil
	}
	return user
}

// OptionalAuth middleware qui récupère l'utilisateur s'il est connecté (optionnel)
func OptionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Essayer de récupérer le token
		cookie, err := r.Cookie("auth_token")
		if err == nil {
			token := utils.ExtractTokenFromCookie(cookie.Value)
			if sessionInfo, err := utils.ValidateJWT(token); err == nil {
				// Ajouter au contexte si valide
				ctx := context.WithValue(r.Context(), UserContextKey, sessionInfo)
				r = r.WithContext(ctx)
			}
		}

		next(w, r)
	}
}

// RequireRole middleware qui vérifie le rôle de l'utilisateur
func RequireRole(role string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return RequireAuth(func(w http.ResponseWriter, r *http.Request) {
			user := GetUserFromContext(r)
			if user == nil || user.Role != role {
				log.Printf("🚫 Accès refusé - rôle requis: %s, rôle utilisateur: %s", role, user.Role)
				http.Error(w, "Accès interdit", http.StatusForbidden)
				return
			}
			next(w, r)
		})
	}
}

// LogRequest middleware pour logger les requêtes
func LogRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromContext(r)
		userInfo := "anonyme"
		if user != nil {
			userInfo = user.Username
		}

		log.Printf("📝 %s %s - User: %s", r.Method, r.URL.Path, userInfo)
		next(w, r)
	}
}

// CORS middleware pour gérer les CORS
func CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}
