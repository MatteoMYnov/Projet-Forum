package repositories

import (
	"database/sql"
	"fmt"
	"forum/models"
)

// ReactionRepository gère les opérations sur les réactions
type ReactionRepository struct {
	db *sql.DB
}

// NewReactionRepository crée une nouvelle instance du repository
func NewReactionRepository(db *sql.DB) *ReactionRepository {
	return &ReactionRepository{db: db}
}

// Create ajoute une nouvelle réaction
func (r *ReactionRepository) Create(reaction models.Reaction) (*models.Reaction, error) {
	query := `
		INSERT INTO reactions (user_id, thread_id, message_id, reaction_type, created_at) 
		VALUES (?, ?, ?, ?, NOW())
	`
	
	result, err := r.db.Exec(query, reaction.UserID, reaction.ThreadID, reaction.MessageID, reaction.ReactionType)
	if err != nil {
		return nil, fmt.Errorf("erreur création réaction: %v", err)
	}
	
	reactionID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération ID réaction: %v", err)
	}
	
	reaction.ID = int(reactionID)
	return &reaction, nil
}

// Delete supprime une réaction
func (r *ReactionRepository) Delete(userID int, threadID *int, messageID *int, reactionType string) error {
	var query string
	var args []interface{}
	
	if threadID != nil {
		query = `DELETE FROM reactions WHERE user_id = ? AND thread_id = ? AND reaction_type = ?`
		args = []interface{}{userID, *threadID, reactionType}
	} else if messageID != nil {
		query = `DELETE FROM reactions WHERE user_id = ? AND message_id = ? AND reaction_type = ?`
		args = []interface{}{userID, *messageID, reactionType}
	} else {
		return fmt.Errorf("thread_id ou message_id requis")
	}
	
	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("erreur suppression réaction: %v", err)
	}
	
	return nil
}

// GetUserReaction récupère la réaction d'un utilisateur sur un thread ou message
func (r *ReactionRepository) GetUserReaction(userID int, threadID *int, messageID *int) (*models.Reaction, error) {
	var query string
	var args []interface{}
	
	if threadID != nil {
		query = `
			SELECT id_reaction, user_id, thread_id, message_id, reaction_type, created_at
			FROM reactions 
			WHERE user_id = ? AND thread_id = ?
		`
		args = []interface{}{userID, *threadID}
	} else if messageID != nil {
		query = `
			SELECT id_reaction, user_id, thread_id, message_id, reaction_type, created_at
			FROM reactions 
			WHERE user_id = ? AND message_id = ?
		`
		args = []interface{}{userID, *messageID}
	} else {
		return nil, fmt.Errorf("thread_id ou message_id requis")
	}
	
	var reaction models.Reaction
	err := r.db.QueryRow(query, args...).Scan(
		&reaction.ID,
		&reaction.UserID,
		&reaction.ThreadID,
		&reaction.MessageID,
		&reaction.ReactionType,
		&reaction.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Pas de réaction trouvée
		}
		return nil, fmt.Errorf("erreur récupération réaction: %v", err)
	}
	
	return &reaction, nil
}

// GetThreadReactionCounts récupère les comptes de réactions pour un thread
func (r *ReactionRepository) GetThreadReactionCounts(threadID int) (map[string]int, error) {
	query := `
		SELECT reaction_type, COUNT(*) 
		FROM reactions 
		WHERE thread_id = ? 
		GROUP BY reaction_type
	`
	
	rows, err := r.db.Query(query, threadID)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération comptes réactions: %v", err)
	}
	defer rows.Close()
	
	counts := make(map[string]int)
	for rows.Next() {
		var reactionType string
		var count int
		
		err := rows.Scan(&reactionType, &count)
		if err != nil {
			return nil, fmt.Errorf("erreur scan compte réaction: %v", err)
		}
		
		counts[reactionType] = count
	}
	
	return counts, nil
}

// GetMessageReactionCounts récupère les comptes de réactions pour un message
func (r *ReactionRepository) GetMessageReactionCounts(messageID int) (map[string]int, error) {
	query := `
		SELECT reaction_type, COUNT(*) 
		FROM reactions 
		WHERE message_id = ? 
		GROUP BY reaction_type
	`
	
	rows, err := r.db.Query(query, messageID)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération comptes réactions: %v", err)
	}
	defer rows.Close()
	
	counts := make(map[string]int)
	for rows.Next() {
		var reactionType string
		var count int
		
		err := rows.Scan(&reactionType, &count)
		if err != nil {
			return nil, fmt.Errorf("erreur scan compte réaction: %v", err)
		}
		
		counts[reactionType] = count
	}
	
	return counts, nil
}

// UpdateThreadCounts met à jour les comptes de like/dislike dans la table threads
func (r *ReactionRepository) UpdateThreadCounts(threadID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("erreur démarrage transaction: %v", err)
	}
	defer tx.Rollback()
	
	// Compter les likes
	var likeCount int
	err = tx.QueryRow(`SELECT COUNT(*) FROM reactions WHERE thread_id = ? AND reaction_type = 'like'`, threadID).Scan(&likeCount)
	if err != nil {
		return fmt.Errorf("erreur compte likes: %v", err)
	}
	
		// Compter les dislikes
	var dislikeCount int
	err = tx.QueryRow(`SELECT COUNT(*) FROM reactions WHERE thread_id = ? AND reaction_type = 'dislike'`, threadID).Scan(&dislikeCount)
	if err != nil {
		return fmt.Errorf("erreur compte dislikes: %v", err)
	}

	// Compter les loves
	var loveCount int
	err = tx.QueryRow(`SELECT COUNT(*) FROM reactions WHERE thread_id = ? AND reaction_type = 'love'`, threadID).Scan(&loveCount)
	if err != nil {
		return fmt.Errorf("erreur compte loves: %v", err)
	}

	// Mettre à jour le thread
	_, err = tx.Exec(`UPDATE threads SET like_count = ?, dislike_count = ?, love_count = ? WHERE id_thread = ?`, likeCount, dislikeCount, loveCount, threadID)
	if err != nil {
		return fmt.Errorf("erreur mise à jour thread: %v", err)
	}
	
	return tx.Commit()
}

// UpdateMessageCounts met à jour les comptes de like/dislike dans la table messages
func (r *ReactionRepository) UpdateMessageCounts(messageID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("erreur démarrage transaction: %v", err)
	}
	defer tx.Rollback()
	
	// Compter les likes
	var likeCount int
	err = tx.QueryRow(`SELECT COUNT(*) FROM reactions WHERE message_id = ? AND reaction_type = 'like'`, messageID).Scan(&likeCount)
	if err != nil {
		return fmt.Errorf("erreur compte likes: %v", err)
	}
	
		// Compter les dislikes
	var dislikeCount int
	err = tx.QueryRow(`SELECT COUNT(*) FROM reactions WHERE message_id = ? AND reaction_type = 'dislike'`, messageID).Scan(&dislikeCount)
	if err != nil {
		return fmt.Errorf("erreur compte dislikes: %v", err)
	}

	// Compter les loves
	var loveCount int
	err = tx.QueryRow(`SELECT COUNT(*) FROM reactions WHERE message_id = ? AND reaction_type = 'love'`, messageID).Scan(&loveCount)
	if err != nil {
		return fmt.Errorf("erreur compte loves: %v", err)
	}

	// Mettre à jour le message
	_, err = tx.Exec(`UPDATE messages SET like_count = ?, dislike_count = ?, love_count = ? WHERE id_message = ?`, likeCount, dislikeCount, loveCount, messageID)
	if err != nil {
		return fmt.Errorf("erreur mise à jour message: %v", err)
	}
	
	return tx.Commit()
}

// GetReactionsByThread récupère toutes les réactions d'un thread avec les informations utilisateur
func (r *ReactionRepository) GetReactionsByThread(threadID int) ([]models.Reaction, error) {
	query := `
		SELECT r.id_reaction, r.user_id, r.thread_id, r.message_id, r.reaction_type, r.created_at
		FROM reactions r
		WHERE r.thread_id = ?
		ORDER BY r.created_at DESC
	`
	
	rows, err := r.db.Query(query, threadID)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération réactions: %v", err)
	}
	defer rows.Close()
	
	var reactions []models.Reaction
	for rows.Next() {
		var reaction models.Reaction
		
		err := rows.Scan(
			&reaction.ID,
			&reaction.UserID,
			&reaction.ThreadID,
			&reaction.MessageID,
			&reaction.ReactionType,
			&reaction.CreatedAt,
		)
		
		if err != nil {
			return nil, fmt.Errorf("erreur scan réaction: %v", err)
		}
		
		reactions = append(reactions, reaction)
	}
	
	return reactions, nil
}

// CheckReactionExists vérifie si une réaction existe déjà
func (r *ReactionRepository) CheckReactionExists(userID int, threadID *int, messageID *int, reactionType string) (bool, error) {
	var query string
	var args []interface{}
	
	if threadID != nil {
		query = `SELECT COUNT(*) FROM reactions WHERE user_id = ? AND thread_id = ? AND reaction_type = ?`
		args = []interface{}{userID, *threadID, reactionType}
	} else if messageID != nil {
		query = `SELECT COUNT(*) FROM reactions WHERE user_id = ? AND message_id = ? AND reaction_type = ?`
		args = []interface{}{userID, *messageID, reactionType}
	} else {
		return false, fmt.Errorf("thread_id ou message_id requis")
	}
	
	var count int
	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("erreur vérification réaction: %v", err)
	}
	
	return count > 0, nil
} 