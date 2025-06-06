package controllers

import (
	"database/sql"
	"encoding/json"
	"forum/middleware"
	"forum/models"
	"forum/services"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// UserRepository gère les opérations sur les utilisateurs
type UserControllers struct {
	authService *services.AuthService
}

// NewUserControllers crée une nouvelle instance du controller
func NewUserControllers(db *sql.DB) *UserControllers {
	return &UserControllers{
		authService: services.NewAuthService(db),
	}
}

func (c *UserControllers) UserRouter(r *http.ServeMux) {
	// Routes pour l'authentification
	r.HandleFunc("/register", c.RegisterPage)
	r.HandleFunc("/login", c.LoginPage)
	r.HandleFunc("/home", c.HomePage)
	r.HandleFunc("/profile", middleware.RequireAuth(c.ProfilePage))

	// Handlers pour les actions
	r.HandleFunc("/api/register", c.RegisterHandler)
	r.HandleFunc("/api/login", c.LoginHandler)
	r.HandleFunc("/api/logout", c.LogoutHandler)
	r.HandleFunc("/api/profile", middleware.RequireAuth(c.ProfileAPI))
}

// RegisterPage affiche la page d'inscription
func (c *UserControllers) RegisterPage(w http.ResponseWriter, r *http.Request) {
	// Si déjà connecté, rediriger vers le profil
	if user := middleware.GetUserFromContext(r); user != nil {
		http.Redirect(w, r, "/profile", http.StatusFound)
		return
	}

	http.ServeFile(w, r, "./website/template/register.html")
}

// LoginPage affiche la page de connexion
func (c *UserControllers) LoginPage(w http.ResponseWriter, r *http.Request) {
	// Si déjà connecté, rediriger vers le profil
	if user := middleware.GetUserFromContext(r); user != nil {
		http.Redirect(w, r, "/profile", http.StatusFound)
		return
	}

	http.ServeFile(w, r, "./website/template/login.html")
}

// HomePage affiche la page d'accueil
func (c *UserControllers) HomePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./website/template/home.html")
}

// ProfilePage affiche la page de profil (nécessite authentification)
func (c *UserControllers) ProfilePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./website/template/profile.html")
}

// RegisterHandler gère l'inscription des utilisateurs
func (c *UserControllers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer les données du formulaire
	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	log.Printf("📝 Tentative d'inscription: %s (%s)", username, email)

	// Validation basique
	if username == "" || email == "" || password == "" {
		log.Printf("❌ Champs manquants pour l'inscription")
		showErrorPage(w, r, "Tous les champs sont requis", "/register")
		return
	}

	// Créer la requête d'inscription
	registerReq := models.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	}

	// Appeler le service d'inscription
	user, err := c.authService.Register(registerReq)
	if err != nil {
		log.Printf("❌ Erreur inscription: %v", err)
		showErrorPage(w, r, err.Error(), "/register")
		return
	}

	log.Printf("✅ Inscription réussie: %s (ID: %d)", user.Username, user.ID)

	// Redirection avec message de succès
	http.Redirect(w, r, "/login?message=inscription_reussie", http.StatusSeeOther)
}

// LoginHandler gère la connexion des utilisateurs
func (c *UserControllers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer les données du formulaire
	identifier := strings.TrimSpace(r.FormValue("identifiant"))
	password := r.FormValue("password")

	log.Printf("🔑 Tentative de connexion: %s", identifier)

	// Validation basique
	if identifier == "" || password == "" {
		log.Printf("❌ Identifiants manquants")
		showErrorPage(w, r, "Tous les champs sont requis", "/login")
		return
	}

	// Créer la requête de connexion
	loginReq := models.LoginRequest{
		Identifier: identifier,
		Password:   password,
	}

	// Appeler le service de connexion
	user, token, err := c.authService.Login(loginReq)
	if err != nil {
		log.Printf("❌ Erreur connexion: %v", err)
		showErrorPage(w, r, "Identifiants incorrects", "/login")
		return
	}

	// Créer le cookie d'authentification
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true,
		Secure:   false, // Mettre à true en production avec HTTPS
		SameSite: http.SameSiteStrictMode,
		MaxAge:   24 * 60 * 60, // 24 heures
		Path:     "/",
	})

	log.Printf("✅ Connexion réussie: %s (ID: %d)", user.Username, user.ID)

	// Redirection vers le profil
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

// LogoutHandler gère la déconnexion
func (c *UserControllers) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Supprimer le cookie d'authentification
	http.SetCookie(w, &http.Cookie{
		Name:   "auth_token",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	})

	log.Printf("👋 Déconnexion effectuée")

	// Redirection vers la page d'accueil
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

// ProfileAPI retourne les données du profil en JSON
func (c *UserControllers) ProfileAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer l'utilisateur depuis le contexte
	sessionInfo := middleware.GetUserFromContext(r)
	if sessionInfo == nil {
		http.Error(w, "Non authentifié", http.StatusUnauthorized)
		return
	}

	// Récupérer les données complètes de l'utilisateur
	user, err := c.authService.GetUserByID(sessionInfo.UserID)
	if err != nil {
		log.Printf("❌ Erreur récupération profil: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// Retourner en JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.APIResponse{
		Success: true,
		Data:    user,
	})
}

// Helper functions

// showErrorPage affiche une page d'erreur simple ou redirige avec un message
func showErrorPage(w http.ResponseWriter, r *http.Request, message string, redirectTo string) {
	// Pour l'instant, on fait une redirection simple
	// En production, vous pourriez créer une vraie page d'erreur
	http.Redirect(w, r, redirectTo+"?error="+message, http.StatusSeeOther)
}

// GetCurrentUser helper pour récupérer l'utilisateur actuel
func (c *UserControllers) GetCurrentUser(r *http.Request) (*models.User, error) {
	sessionInfo := middleware.GetUserFromContext(r)
	if sessionInfo == nil {
		return nil, nil
	}

	return c.authService.GetUserByID(sessionInfo.UserID)
}

// RequireJSON middleware pour s'assurer que le content-type est JSON
func RequireJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Content-Type doit être application/json", http.StatusBadRequest)
			return
		}
		next(w, r)
	}
}

// ParseJSONBody helper pour parser le body JSON
func ParseJSONBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// WriteJSONResponse helper pour écrire une réponse JSON
func WriteJSONResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// WriteErrorResponse helper pour écrire une réponse d'erreur JSON
func WriteErrorResponse(w http.ResponseWriter, message string, status int) {
	WriteJSONResponse(w, models.APIResponse{
		Success: false,
		Error:   message,
	}, status)
}

// ParseIntParam helper pour parser un paramètre entier depuis l'URL
func ParseIntParam(r *http.Request, param string) (int, error) {
	value := r.URL.Query().Get(param)
	if value == "" {
		return 0, nil
	}
	return strconv.Atoi(value)
}
