package services

import (
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"forum/models"
	"forum/repositories"
	"forum/utils"
	"log"
	"regexp"
	"strings"
	"time"
)

var (
	ErrWeakPassword      = errors.New("mot de passe trop faible")
	ErrInvalidEmail      = errors.New("email invalide")
	ErrInvalidUsername   = errors.New("nom d'utilisateur invalide")
	ErrUserAlreadyExists = errors.New("utilisateur déjà existant")
	ErrInvalidLogin      = errors.New("identifiants invalides")
	ErrUserBanned        = errors.New("utilisateur banni")
)

// AuthService gère l'authentification
type AuthService struct {
	userRepo *repositories.UserRepository
}

// NewAuthService crée une nouvelle instance du service
func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{
		userRepo: repositories.NewUserRepository(db),
	}
}

// Register inscrit un nouvel utilisateur
func (s *AuthService) Register(req models.RegisterRequest) (*models.User, error) {
	log.Printf("📝 Tentative d'inscription: %s (%s)", req.Username, req.Email)

	// 1. Validation des données
	if err := s.validateRegistration(req); err != nil {
		return nil, err
	}

	// 2. Vérifier que l'utilisateur n'existe pas déjà
	if exists, _ := s.userRepo.UsernameExists(req.Username); exists {
		return nil, fmt.Errorf("nom d'utilisateur '%s' déjà utilisé", req.Username)
	}

	if exists, _ := s.userRepo.EmailExists(req.Email); exists {
		return nil, fmt.Errorf("email '%s' déjà utilisé", req.Email)
	}

	// 3. Hasher le mot de passe
	hashedPassword := s.hashPassword(req.Password)

	// 4. Créer l'utilisateur
	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
		Role:         "user",
		IsVerified:   false,
		IsBanned:     false,
	}

	createdUser, err := s.userRepo.CreateUser(user)
	if err != nil {
		log.Printf("❌ Erreur création utilisateur: %v", err)
		return nil, fmt.Errorf("erreur lors de l'inscription: %v", err)
	}

	log.Printf("✅ Inscription réussie: %s (ID: %d)", createdUser.Username, createdUser.ID)
	return createdUser, nil
}

// Login connecte un utilisateur
func (s *AuthService) Login(req models.LoginRequest) (*models.User, string, error) {
	log.Printf("🔑 Tentative de connexion: %s", req.Identifier)

	// 1. Validation basique
	if strings.TrimSpace(req.Identifier) == "" || strings.TrimSpace(req.Password) == "" {
		return nil, "", ErrInvalidLogin
	}

	// 2. Récupérer l'utilisateur
	user, err := s.userRepo.GetUserByIdentifier(req.Identifier)
	if err != nil {
		if err == repositories.ErrUserNotFound {
			log.Printf("⚠️ Utilisateur non trouvé: %s", req.Identifier)
			return nil, "", ErrInvalidLogin
		}
		return nil, "", fmt.Errorf("erreur récupération utilisateur: %v", err)
	}

	// 3. Vérifier si l'utilisateur est banni
	if user.IsBanned {
		log.Printf("🚫 Tentative de connexion d'un utilisateur banni: %s", user.Username)
		return nil, "", ErrUserBanned
	}

	// 4. Vérifier le mot de passe
	hashedPassword := s.hashPassword(req.Password)
	if hashedPassword != user.PasswordHash {
		log.Printf("❌ Mot de passe incorrect pour: %s", req.Identifier)
		return nil, "", ErrInvalidLogin
	}

	// 5. Mettre à jour la dernière connexion
	if err := s.userRepo.UpdateLastLogin(user.ID); err != nil {
		log.Printf("⚠️ Erreur mise à jour last_login: %v", err)
	}

	// 6. Générer un token JWT
	token, err := utils.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, "", fmt.Errorf("erreur génération token: %v", err)
	}

	log.Printf("✅ Connexion réussie: %s (ID: %d)", user.Username, user.ID)
	return user, token, nil
}

// GetUserByID récupère un utilisateur par son ID
func (s *AuthService) GetUserByID(id int) (*models.User, error) {
	return s.userRepo.GetUserByID(id)
}

// validateRegistration valide les données d'inscription
func (s *AuthService) validateRegistration(req models.RegisterRequest) error {
	// Validation du nom d'utilisateur
	if err := s.validateUsername(req.Username); err != nil {
		return err
	}

	// Validation de l'email
	if err := s.validateEmail(req.Email); err != nil {
		return err
	}

	// Validation du mot de passe
	if err := s.validatePassword(req.Password); err != nil {
		return err
	}

	return nil
}

// validateUsername valide le nom d'utilisateur
func (s *AuthService) validateUsername(username string) error {
	username = strings.TrimSpace(username)

	if len(username) < 3 {
		return fmt.Errorf("nom d'utilisateur trop court (minimum 3 caractères)")
	}

	if len(username) > 50 {
		return fmt.Errorf("nom d'utilisateur trop long (maximum 50 caractères)")
	}

	// Caractères alphanumériques + underscore uniquement
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	if !matched {
		return fmt.Errorf("nom d'utilisateur invalide (lettres, chiffres et _ uniquement)")
	}

	return nil
}

// validateEmail valide l'adresse email
func (s *AuthService) validateEmail(email string) error {
	email = strings.TrimSpace(email)

	if len(email) == 0 {
		return fmt.Errorf("email requis")
	}

	// Regex basique pour l'email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}

	return nil
}

// validatePassword valide le mot de passe selon les critères de sécurité
func (s *AuthService) validatePassword(password string) error {
	if len(password) < 12 {
		return fmt.Errorf("mot de passe trop court (minimum 12 caractères)")
	}

	// Vérifier la présence d'au moins une majuscule
	hasUpper, _ := regexp.MatchString(`[A-Z]`, password)
	if !hasUpper {
		return fmt.Errorf("mot de passe doit contenir au moins une majuscule")
	}

	// Vérifier la présence d'au moins un caractère spécial
	hasSpecial, _ := regexp.MatchString(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`, password)
	if !hasSpecial {
		return fmt.Errorf("mot de passe doit contenir au moins un caractère spécial")
	}

	return nil
}

// hashPassword hash un mot de passe avec SHA512
func (s *AuthService) hashPassword(password string) string {
	hasher := sha512.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}
