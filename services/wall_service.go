package services

import (
	"database/sql"
	"fmt"
	"forum/models"
	"log"
	"strings"
	"time"
)

type WallService struct {
	db *sql.DB
}

func NewWallService(db *sql.DB) *WallService {
	return &WallService{db: db}
}

// CreateWallPost crée une nouvelle publication sur le mur
func (s *WallService) CreateWallPost(userID, authorID int, content string) (*models.WallPost, error) {
	log.Printf("📝 Création d'un post sur le mur de l'utilisateur %d par l'auteur %d", userID, authorID)
	
	// Valider le contenu
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, fmt.Errorf("le contenu ne peut pas être vide")
	}
	if len(content) > 1000 {
		return nil, fmt.Errorf("le contenu ne peut pas dépasser 1000 caractères")
	}
	
	query := `
		INSERT INTO wall_posts (user_id, author_id, content, created_at, updated_at)
		VALUES (?, ?, ?, NOW(), NOW())
	`
	
	result, err := s.db.Exec(query, userID, authorID, content)
	if err != nil {
		log.Printf("❌ Erreur lors de la création du post: %v", err)
		return nil, fmt.Errorf("erreur lors de la création du post: %v", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération de l'ID: %v", err)
	}
	
	wallPost := &models.WallPost{
		ID:        int(id),
		UserID:    userID,
		AuthorID:  authorID,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	log.Printf("✅ Post créé avec succès (ID: %d)", id)
	return wallPost, nil
}

// GetWallPosts récupère tous les posts du mur d'un utilisateur avec les informations des auteurs
func (s *WallService) GetWallPosts(userID int) ([]models.WallPostWithAuthor, error) {
	log.Printf("📖 Récupération des posts du mur pour l'utilisateur %d", userID)
	
	query := `
		SELECT 
			wp.id, wp.user_id, wp.author_id, wp.content, wp.created_at,
			u.username as author_name, u.email as author_email,
			COALESCE(u.profile_picture, '/img/avatars/default-avatar.png') as avatar_path
		FROM wall_posts wp
		JOIN users u ON wp.author_id = u.id_user
		WHERE wp.user_id = ?
		ORDER BY wp.created_at DESC
		LIMIT 50
	`
	
	rows, err := s.db.Query(query, userID)
	if err != nil {
		log.Printf("❌ Erreur lors de la récupération des posts: %v", err)
		return nil, fmt.Errorf("erreur lors de la récupération des posts: %v", err)
	}
	defer rows.Close()
	
	var posts []models.WallPostWithAuthor
	for rows.Next() {
		var post models.WallPostWithAuthor
		err := rows.Scan(
			&post.ID, &post.UserID, &post.AuthorID, &post.Content, &post.CreatedAt,
			&post.AuthorName, &post.AuthorEmail, &post.AvatarPath,
		)
		if err != nil {
			log.Printf("❌ Erreur lors du scan du post: %v", err)
			continue
		}
		posts = append(posts, post)
	}
	
	log.Printf("✅ %d posts récupérés pour l'utilisateur %d", len(posts), userID)
	return posts, nil
}

// DeleteWallPost supprime un post du mur
func (s *WallService) DeleteWallPost(postID, authorID int) error {
	log.Printf("🗑️ Suppression du post %d par l'auteur %d", postID, authorID)
	
	query := `DELETE FROM wall_posts WHERE id = ? AND author_id = ?`
	result, err := s.db.Exec(query, postID, authorID)
	if err != nil {
		log.Printf("❌ Erreur lors de la suppression: %v", err)
		return fmt.Errorf("erreur lors de la suppression: %v", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erreur lors de la vérification: %v", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("post non trouvé ou vous n'êtes pas autorisé à le supprimer")
	}
	
	log.Printf("✅ Post %d supprimé avec succès", postID)
	return nil
}

// GetWallPostsCount retourne le nombre de posts sur le mur d'un utilisateur
func (s *WallService) GetWallPostsCount(userID int) (int, error) {
	query := `SELECT COUNT(*) FROM wall_posts WHERE user_id = ?`
	var count int
	err := s.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("erreur lors du comptage des posts: %v", err)
	}
	return count, nil
} 