package repositories

import (
	"database/sql"
	"fmt"
	"forum/models"
	"time"
)

// MessageRepository gère les opérations sur les messages
type MessageRepository struct {
	db *sql.DB
}

// NewMessageRepository crée une nouvelle instance du repository
func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

// Create crée un nouveau message
func (r *MessageRepository) Create(message *models.Message) (*models.Message, error) {
	query := `
		INSERT INTO messages (thread_id, author_id, content, parent_message_id, created_at, updated_at) 
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`
	
	result, err := r.db.Exec(query, message.ThreadID, message.AuthorID, message.Content, message.ParentMessageID)
	if err != nil {
		return nil, fmt.Errorf("erreur création message: %v", err)
	}
	
	messageID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération ID message: %v", err)
	}
	
	message.ID = int(messageID)
	message.CreatedAt = time.Now()
	message.UpdatedAt = time.Now()
	
	return message, nil
}

// GetByID récupère un message par son ID
func (r *MessageRepository) GetByID(messageID int) (*models.Message, error) {
	query := `
		SELECT m.id_message, m.thread_id, m.author_id, m.content, m.parent_message_id, 
		       m.created_at, m.updated_at, m.like_count, m.dislike_count, m.is_edited,
		       u.username, u.profile_picture
		FROM messages m
		JOIN users u ON m.author_id = u.id_user
		WHERE m.id_message = ?
	`
	
	var message models.Message
	var author models.User
	
	err := r.db.QueryRow(query, messageID).Scan(
		&message.ID,
		&message.ThreadID,
		&message.AuthorID,
		&message.Content,
		&message.ParentMessageID,
		&message.CreatedAt,
		&message.UpdatedAt,
		&message.LikeCount,
		&message.DislikeCount,
		&message.IsEdited,
		&author.Username,
		&author.ProfilePicture,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("message non trouvé")
		}
		return nil, fmt.Errorf("erreur récupération message: %v", err)
	}
	
	author.ID = message.AuthorID
	message.Author = &author
	
	return &message, nil
}

// GetByThreadID récupère tous les messages d'un thread
func (r *MessageRepository) GetByThreadID(threadID int) ([]models.Message, error) {
	query := `
		SELECT m.id_message, m.thread_id, m.author_id, m.content, m.parent_message_id, 
		       m.created_at, m.updated_at, m.like_count, m.dislike_count, m.is_edited,
		       u.username, u.profile_picture, u.email
		FROM messages m
		JOIN users u ON m.author_id = u.id_user
		WHERE m.thread_id = ?
		ORDER BY m.created_at ASC
	`
	
	rows, err := r.db.Query(query, threadID)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération messages: %v", err)
	}
	defer rows.Close()
	
	var messages []models.Message
	for rows.Next() {
		var message models.Message
		var author models.User
		
		err := rows.Scan(
			&message.ID,
			&message.ThreadID,
			&message.AuthorID,
			&message.Content,
			&message.ParentMessageID,
			&message.CreatedAt,
			&message.UpdatedAt,
			&message.LikeCount,
			&message.DislikeCount,
			&message.IsEdited,
			&author.Username,
			&author.ProfilePicture,
			&author.Email,
		)
		
		if err != nil {
			return nil, fmt.Errorf("erreur scan message: %v", err)
		}
		
		author.ID = message.AuthorID
		message.Author = &author
		messages = append(messages, message)
	}
	
	return messages, nil
}

// Delete supprime un message
func (r *MessageRepository) Delete(messageID int) error {
	query := `DELETE FROM messages WHERE id_message = ?`
	
	result, err := r.db.Exec(query, messageID)
	if err != nil {
		return fmt.Errorf("erreur suppression message: %v", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erreur vérification suppression: %v", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("message non trouvé")
	}
	
	return nil
}

// Update met à jour un message
func (r *MessageRepository) Update(messageID int, content string) error {
	query := `
		UPDATE messages 
		SET content = ?, updated_at = NOW(), is_edited = TRUE 
		WHERE id_message = ?
	`
	
	result, err := r.db.Exec(query, content, messageID)
	if err != nil {
		return fmt.Errorf("erreur mise à jour message: %v", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erreur vérification mise à jour: %v", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("message non trouvé")
	}
	
	return nil
}

// GetByUserID récupère tous les messages d'un utilisateur
func (r *MessageRepository) GetByUserID(userID int, limit, offset int) ([]models.Message, error) {
	query := `
		SELECT m.id_message, m.thread_id, m.author_id, m.content, m.parent_message_id, 
		       m.created_at, m.updated_at, m.like_count, m.dislike_count, m.is_edited,
		       t.title as thread_title
		FROM messages m
		JOIN threads t ON m.thread_id = t.id_thread
		WHERE m.author_id = ?
		ORDER BY m.created_at DESC
		LIMIT ? OFFSET ?
	`
	
	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération messages utilisateur: %v", err)
	}
	defer rows.Close()
	
	var messages []models.Message
	for rows.Next() {
		var message models.Message
		var threadTitle string
		
		err := rows.Scan(
			&message.ID,
			&message.ThreadID,
			&message.AuthorID,
			&message.Content,
			&message.ParentMessageID,
			&message.CreatedAt,
			&message.UpdatedAt,
			&message.LikeCount,
			&message.DislikeCount,
			&message.IsEdited,
			&threadTitle,
		)
		
		if err != nil {
			return nil, fmt.Errorf("erreur scan message utilisateur: %v", err)
		}
		
		messages = append(messages, message)
	}
	
	return messages, nil
} 