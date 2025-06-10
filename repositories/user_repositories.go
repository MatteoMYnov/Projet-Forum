package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/models"
	"log"
	"time"
)

var (
	ErrUserNotFound       = errors.New("utilisateur non trouvé")
	ErrUserAlreadyExists  = errors.New("utilisateur déjà existant")
	ErrInvalidCredentials = errors.New("identifiants invalides")
)

// UserRepository gère les opérations sur les utilisateurs
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository crée une nouvelle instance du repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// CreateUser crée un nouvel utilisateur
func (r *UserRepository) CreateUser(user models.User) (*models.User, error) {
	query := `
		INSERT INTO users (username, email, password_hash, profile_picture, created_at, role) 
		VALUES (?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.ProfilePicture,
		time.Now(),
		"user")

	if err != nil {
		log.Printf("❌ Erreur création utilisateur: %v", err)
		// Vérifier si c'est une erreur de doublon
		if isDuplicateError(err) {
			return nil, ErrUserAlreadyExists
		}
		return nil, fmt.Errorf("erreur création utilisateur: %v", err)
	}

	// Récupérer l'ID généré
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération ID: %v", err)
	}

	user.ID = int(id)
	user.CreatedAt = time.Now()
	user.Role = "user"

	log.Printf("✅ Utilisateur créé: %s (ID: %d)", user.Username, user.ID)
	return &user, nil
}

// GetUserByID récupère un utilisateur par son ID
func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	query := `
		SELECT id_user, username, email, password_hash, profile_picture, bio, 
			   created_at, last_login, is_verified, is_banned, role, 
			   follower_count, following_count, thread_count
		FROM users 
		WHERE id_user = ?`

	var user models.User
	row := r.db.QueryRow(query, id)

	err := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.ProfilePicture, &user.Bio, &user.CreatedAt, &user.LastLogin,
		&user.IsVerified, &user.IsBanned, &user.Role,
		&user.FollowerCount, &user.FollowingCount, &user.ThreadCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("erreur récupération utilisateur: %v", err)
	}

	return &user, nil
}

// GetUserByUsername récupère un utilisateur par son nom d'utilisateur
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	query := `
		SELECT id_user, username, email, password_hash, profile_picture, bio, 
			   created_at, last_login, is_verified, is_banned, role, 
			   follower_count, following_count, thread_count
		FROM users 
		WHERE username = ?`

	var user models.User
	row := r.db.QueryRow(query, username)

	err := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.ProfilePicture, &user.Bio, &user.CreatedAt, &user.LastLogin,
		&user.IsVerified, &user.IsBanned, &user.Role,
		&user.FollowerCount, &user.FollowingCount, &user.ThreadCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("erreur récupération utilisateur: %v", err)
	}

	return &user, nil
}

// GetUserByEmail récupère un utilisateur par son email
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT id_user, username, email, password_hash, profile_picture, bio, 
			   created_at, last_login, is_verified, is_banned, role, 
			   follower_count, following_count, thread_count
		FROM users 
		WHERE email = ?`

	var user models.User
	row := r.db.QueryRow(query, email)

	err := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.ProfilePicture, &user.Bio, &user.CreatedAt, &user.LastLogin,
		&user.IsVerified, &user.IsBanned, &user.Role,
		&user.FollowerCount, &user.FollowingCount, &user.ThreadCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("erreur récupération utilisateur: %v", err)
	}

	return &user, nil
}

// GetUserByIdentifier récupère un utilisateur par nom d'utilisateur ou email
func (r *UserRepository) GetUserByIdentifier(identifier string) (*models.User, error) {
	query := `
		SELECT id_user, username, email, password_hash, profile_picture, bio, 
			   created_at, last_login, is_verified, is_banned, role, 
			   follower_count, following_count, thread_count
		FROM users 
		WHERE username = ? OR email = ?`

	var user models.User
	row := r.db.QueryRow(query, identifier, identifier)

	err := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.ProfilePicture, &user.Bio, &user.CreatedAt, &user.LastLogin,
		&user.IsVerified, &user.IsBanned, &user.Role,
		&user.FollowerCount, &user.FollowingCount, &user.ThreadCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("erreur récupération utilisateur: %v", err)
	}

	return &user, nil
}

// UpdateLastLogin met à jour la dernière connexion
func (r *UserRepository) UpdateLastLogin(userID int) error {
	query := `UPDATE users SET last_login = ? WHERE id_user = ?`

	_, err := r.db.Exec(query, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("erreur mise à jour last_login: %v", err)
	}

	return nil
}

// UsernameExists vérifie si un nom d'utilisateur existe déjà
func (r *UserRepository) UsernameExists(username string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE username = ?`

	var count int
	err := r.db.QueryRow(query, username).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("erreur vérification username: %v", err)
	}

	return count > 0, nil
}

// EmailExists vérifie si un email existe déjà
func (r *UserRepository) EmailExists(email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = ?`

	var count int
	err := r.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("erreur vérification email: %v", err)
	}

	return count > 0, nil
}

// isDuplicateError vérifie si l'erreur est due à une contrainte d'unicité
func isDuplicateError(err error) bool {
	// Pour MySQL, erreur 1062 = Duplicate entry
	return err != nil && (
	// Vérifier les messages d'erreur MySQL
	fmt.Sprintf("%v", err) == "Error 1062" ||
		// Ou vérifier le contenu du message
		fmt.Sprintf("%v", err) == "Duplicate entry" ||
		// Vérifier avec UNIQUE constraint
		fmt.Sprintf("%v", err) == "UNIQUE constraint")
}
