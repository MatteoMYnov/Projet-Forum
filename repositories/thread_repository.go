package repositories

import (
	"database/sql"
	"fmt"
	"forum/models"
	"log"
	"strings"
	"time"
)

// ThreadRepository gère les opérations de base de données pour les threads
type ThreadRepository struct {
	db *sql.DB
}

// NewThreadRepository crée une nouvelle instance du repository
func NewThreadRepository(db *sql.DB) *ThreadRepository {
	return &ThreadRepository{db: db}
}

// Create crée un nouveau thread
func (r *ThreadRepository) Create(thread *models.Thread) (*models.Thread, error) {
	query := `
		INSERT INTO threads (title, content, author_id, category_id, status, created_at, updated_at, last_activity)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := r.db.Exec(
		query,
		thread.Title,
		thread.Content,
		thread.AuthorID,
		thread.CategoryID,
		"open", // Status par défaut
		now,
		now,
		now,
	)

	if err != nil {
		log.Printf("❌ Erreur création thread: %v", err)
		return nil, fmt.Errorf("erreur lors de la création du thread: %v", err)
	}

	// Récupérer l'ID du thread créé
	threadID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération ID thread: %v", err)
	}

	// Retourner le thread avec son ID
	thread.ID = int(threadID)
	thread.Status = "open"
	thread.CreatedAt = now
	thread.UpdatedAt = now
	thread.LastActivity = now

	log.Printf("✅ Thread créé avec succès: ID=%d, Titre=%s", thread.ID, thread.Title)
	return thread, nil
}

// GetByID récupère un thread par son ID
func (r *ThreadRepository) GetByID(threadID int) (*models.Thread, error) {
	query := `
		SELECT t.id_thread, t.title, t.content, t.author_id, t.category_id, t.status,
		       t.created_at, t.updated_at, t.is_pinned, t.view_count, t.like_count,
		       t.dislike_count, t.love_count, t.message_count, t.last_activity,
		       u.username, u.email, u.profile_picture
		FROM threads t
		JOIN users u ON t.author_id = u.id_user
		WHERE t.id_thread = ?
	`

	var thread models.Thread
	var author models.User

	err := r.db.QueryRow(query, threadID).Scan(
		&thread.ID,
		&thread.Title,
		&thread.Content,
		&thread.AuthorID,
		&thread.CategoryID,
		&thread.Status,
		&thread.CreatedAt,
		&thread.UpdatedAt,
		&thread.IsPinned,
		&thread.ViewCount,
		&thread.LikeCount,
		&thread.DislikeCount,
		&thread.LoveCount,
		&thread.MessageCount,
		&thread.LastActivity,
		&author.Username,
		&author.Email,
		&author.ProfilePicture,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("thread non trouvé")
		}
		return nil, fmt.Errorf("erreur récupération thread: %v", err)
	}

	// Attacher l'auteur
	author.ID = thread.AuthorID
	thread.Author = &author

	return &thread, nil
}

// GetAll récupère tous les threads avec pagination
func (r *ThreadRepository) GetAll(limit, offset int) ([]models.Thread, error) {
	query := `
		SELECT t.id_thread, t.title, t.content, t.author_id, t.category_id, t.status,
		       t.created_at, t.updated_at, t.is_pinned, t.view_count, t.like_count,
		       t.dislike_count, t.love_count, t.message_count, t.last_activity,
		       u.username, u.email, u.profile_picture
		FROM threads t
		JOIN users u ON t.author_id = u.id_user
		ORDER BY t.is_pinned DESC, t.last_activity DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération threads: %v", err)
	}
	defer rows.Close()

	var threads []models.Thread
	for rows.Next() {
		var thread models.Thread
		var author models.User

		err := rows.Scan(
			&thread.ID,
			&thread.Title,
			&thread.Content,
			&thread.AuthorID,
			&thread.CategoryID,
			&thread.Status,
			&thread.CreatedAt,
			&thread.UpdatedAt,
			&thread.IsPinned,
			&thread.ViewCount,
			&thread.LikeCount,
			&thread.DislikeCount,
			&thread.LoveCount,
			&thread.MessageCount,
			&thread.LastActivity,
			&author.Username,
			&author.Email,
			&author.ProfilePicture,
		)

		if err != nil {
			return nil, fmt.Errorf("erreur scan thread: %v", err)
		}

		// Attacher l'auteur
		author.ID = thread.AuthorID
		thread.Author = &author

		threads = append(threads, thread)
	}

	return threads, nil
}

// GetByUserID récupère les threads d'un utilisateur
func (r *ThreadRepository) GetByUserID(userID int, limit, offset int) ([]models.Thread, error) {
	query := `
		SELECT t.id_thread, t.title, t.content, t.author_id, t.category_id, t.status,
		       t.created_at, t.updated_at, t.is_pinned, t.view_count, t.like_count,
		       t.dislike_count, t.love_count, t.message_count, t.last_activity,
		       u.username, u.email, u.profile_picture
		FROM threads t
		JOIN users u ON t.author_id = u.id_user
		WHERE t.author_id = ?
		ORDER BY t.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération threads utilisateur: %v", err)
	}
	defer rows.Close()

	var threads []models.Thread
	for rows.Next() {
		var thread models.Thread
		var author models.User

		err := rows.Scan(
			&thread.ID,
			&thread.Title,
			&thread.Content,
			&thread.AuthorID,
			&thread.CategoryID,
			&thread.Status,
			&thread.CreatedAt,
			&thread.UpdatedAt,
			&thread.IsPinned,
			&thread.ViewCount,
			&thread.LikeCount,
			&thread.DislikeCount,
			&thread.LoveCount,
			&thread.MessageCount,
			&thread.LastActivity,
			&author.Username,
			&author.Email,
			&author.ProfilePicture,
		)

		if err != nil {
			return nil, fmt.Errorf("erreur scan thread: %v", err)
		}

		// Attacher l'auteur
		author.ID = thread.AuthorID
		thread.Author = &author

		threads = append(threads, thread)
	}

	return threads, nil
}

// UpdateViewCount incrémente le nombre de vues d'un thread
func (r *ThreadRepository) UpdateViewCount(threadID int) error {
	query := `UPDATE threads SET view_count = view_count + 1 WHERE id_thread = ?`
	_, err := r.db.Exec(query, threadID)
	if err != nil {
		return fmt.Errorf("erreur mise à jour vue: %v", err)
	}
	return nil
}

// UpdateMessageCount met à jour le nombre de messages d'un thread
func (r *ThreadRepository) UpdateMessageCount(threadID int) error {
	query := `
		UPDATE threads t 
		SET message_count = (
			SELECT COUNT(*) 
			FROM messages m 
			WHERE m.thread_id = t.id_thread
		), 
		last_activity = NOW()
		WHERE t.id_thread = ?
	`
	_, err := r.db.Exec(query, threadID)
	if err != nil {
		return fmt.Errorf("erreur mise à jour nombre messages: %v", err)
	}
	return nil
}

// Delete supprime un thread (soft delete en changeant le status)
func (r *ThreadRepository) Delete(threadID int) error {
	query := `UPDATE threads SET status = 'deleted' WHERE id_thread = ?`
	_, err := r.db.Exec(query, threadID)
	if err != nil {
		return fmt.Errorf("erreur suppression thread: %v", err)
	}
	return nil
}

// GetCategories récupère toutes les catégories actives
func (r *ThreadRepository) GetCategories() ([]models.Category, error) {
	query := `
		SELECT id_category, name, color, description, thread_count, created_at, is_active
		FROM categories
		WHERE is_active = true
		ORDER BY name ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération catégories: %v", err)
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Color,
			&category.Description,
			&category.ThreadCount,
			&category.CreatedAt,
			&category.IsActive,
		)

		if err != nil {
			return nil, fmt.Errorf("erreur scan catégorie: %v", err)
		}

		categories = append(categories, category)
	}

	return categories, nil
}

// ProcessHashtags extrait et traite les hashtags
func ProcessHashtags(content string) []string {
	words := strings.Fields(content)
	var hashtags []string
	
	for _, word := range words {
		if strings.HasPrefix(word, "#") && len(word) > 1 {
			// Nettoyer le hashtag (enlever la ponctuation à la fin)
			hashtag := strings.TrimRight(word[1:], ".,!?;:")
			if hashtag != "" {
				hashtags = append(hashtags, hashtag)
			}
		}
	}
	
	return hashtags
}

// GetTotalCount récupère le nombre total de threads (tous statuts)
func (r *ThreadRepository) GetTotalCount() (int, error) {
	query := `SELECT COUNT(*) FROM threads WHERE status != 'deleted'`
	
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("erreur compte threads: %v", err)
	}
	
	return count, nil
}

// UpdateStatus met à jour le statut d'un thread
func (r *ThreadRepository) UpdateStatus(threadID int, status string) error {
	query := `UPDATE threads SET status = ?, updated_at = NOW() WHERE id_thread = ?`
	
	result, err := r.db.Exec(query, status, threadID)
	if err != nil {
		return fmt.Errorf("erreur mise à jour statut: %v", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erreur vérification mise à jour: %v", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("thread non trouvé")
	}
	
	return nil
}

// GetVisibleThreadsCount récupère le nombre de threads visibles (non archivés)
func (r *ThreadRepository) GetVisibleThreadsCount() (int, error) {
	query := `SELECT COUNT(*) FROM threads WHERE status IN ('open', 'closed')`
	
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("erreur compte threads visibles: %v", err)
	}
	
	return count, nil
}

// GetVisibleThreads récupère les threads visibles (non archivés) avec pagination
func (r *ThreadRepository) GetVisibleThreads(limit, offset int) ([]models.Thread, error) {
	query := `
		SELECT t.id_thread, t.title, t.content, t.author_id, t.category_id, t.status,
		       t.created_at, t.updated_at, t.is_pinned, t.view_count, t.like_count,
		       t.dislike_count, t.love_count, t.message_count, t.last_activity,
		       u.username, u.email, u.profile_picture
		FROM threads t
		JOIN users u ON t.author_id = u.id_user
		WHERE t.status IN ('open', 'closed')
		ORDER BY t.is_pinned DESC, t.last_activity DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération threads visibles: %v", err)
	}
	defer rows.Close()

	var threads []models.Thread
	for rows.Next() {
		var thread models.Thread
		var author models.User

		err := rows.Scan(
			&thread.ID,
			&thread.Title,
			&thread.Content,
			&thread.AuthorID,
			&thread.CategoryID,
			&thread.Status,
			&thread.CreatedAt,
			&thread.UpdatedAt,
			&thread.IsPinned,
			&thread.ViewCount,
			&thread.LikeCount,
			&thread.DislikeCount,
			&thread.LoveCount,
			&thread.MessageCount,
			&thread.LastActivity,
			&author.Username,
			&author.Email,
			&author.ProfilePicture,
		)

		if err != nil {
			return nil, fmt.Errorf("erreur scan thread: %v", err)
		}

		// Attacher l'auteur
		author.ID = thread.AuthorID
		thread.Author = &author

		threads = append(threads, thread)
	}

	return threads, nil
}

// GetCountByStatus récupère le nombre de threads pour un statut donné
func (r *ThreadRepository) GetCountByStatus(status string) (int, error) {
	query := `SELECT COUNT(*) FROM threads WHERE status = ?`
	
	var count int
	err := r.db.QueryRow(query, status).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("erreur compte threads par statut: %v", err)
	}
	
	return count, nil
}

// GetByStatus récupère les threads filtrés par statut avec pagination
func (r *ThreadRepository) GetByStatus(status string, limit, offset int) ([]models.Thread, error) {
	query := `
		SELECT t.id_thread, t.title, t.content, t.author_id, t.category_id, t.status,
		       t.created_at, t.updated_at, t.is_pinned, t.view_count, t.like_count,
		       t.dislike_count, t.love_count, t.message_count, t.last_activity,
		       u.username, u.email, u.profile_picture
		FROM threads t
		JOIN users u ON t.author_id = u.id_user
		WHERE t.status = ?
		ORDER BY t.is_pinned DESC, t.last_activity DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération threads par statut: %v", err)
	}
	defer rows.Close()

	var threads []models.Thread
	for rows.Next() {
		var thread models.Thread
		var author models.User

		err := rows.Scan(
			&thread.ID,
			&thread.Title,
			&thread.Content,
			&thread.AuthorID,
			&thread.CategoryID,
			&thread.Status,
			&thread.CreatedAt,
			&thread.UpdatedAt,
			&thread.IsPinned,
			&thread.ViewCount,
			&thread.LikeCount,
			&thread.DislikeCount,
			&thread.LoveCount,
			&thread.MessageCount,
			&thread.LastActivity,
			&author.Username,
			&author.Email,
			&author.ProfilePicture,
		)

		if err != nil {
			return nil, fmt.Errorf("erreur scan thread: %v", err)
		}

		// Attacher l'auteur
		author.ID = thread.AuthorID
		thread.Author = &author

		threads = append(threads, thread)
	}

	return threads, nil
}

// UpdateTitle met à jour le titre d'un thread
func (r *ThreadRepository) UpdateTitle(threadID int, newTitle string) error {
	query := `UPDATE threads SET title = ?, updated_at = NOW() WHERE id_thread = ?`
	
	result, err := r.db.Exec(query, newTitle, threadID)
	if err != nil {
		return fmt.Errorf("erreur lors de la mise à jour du titre: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erreur vérification mise à jour: %v", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("thread non trouvé")
	}
	
	return nil
} 