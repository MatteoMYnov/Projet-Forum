package services

import (
	"database/sql"
	"fmt"
	"forum/models"
	"forum/repositories"
	"log"
	"time"
)

// ReactionService gère la logique métier pour les réactions
type ReactionService struct {
	reactionRepo *repositories.ReactionRepository
	threadRepo   *repositories.ThreadRepository
}

// NewReactionService crée une nouvelle instance du service
func NewReactionService(db *sql.DB) *ReactionService {
	return &ReactionService{
		reactionRepo: repositories.NewReactionRepository(db),
		threadRepo:   repositories.NewThreadRepository(db),
	}
}

// ToggleThreadReaction ajoute ou supprime une réaction sur un thread
func (s *ReactionService) ToggleThreadReaction(userID, threadID int, reactionType string) (*models.Reaction, error) {
	log.Printf("🔄 ToggleThreadReaction - User %d, Thread %d, Type %s", userID, threadID, reactionType)

	// Valider le type de réaction
	if !isValidReactionType(reactionType) {
		return nil, fmt.Errorf("type de réaction invalide: %s", reactionType)
	}

	// Vérifier si le thread existe
	_, err := s.threadRepo.GetByID(threadID)
	if err != nil {
		return nil, fmt.Errorf("thread non trouvé: %v", err)
	}

	threadIDPtr := &threadID

	// Vérifier si l'utilisateur a déjà une réaction sur ce thread
	existingReaction, err := s.reactionRepo.GetUserReaction(userID, threadIDPtr, nil)
	if err != nil {
		return nil, fmt.Errorf("erreur vérification réaction existante: %v", err)
	}

	// Si l'utilisateur a déjà la même réaction, on la supprime (toggle off)
	if existingReaction != nil && existingReaction.ReactionType == reactionType {
		err = s.reactionRepo.Delete(userID, threadIDPtr, nil, reactionType)
		if err != nil {
			return nil, fmt.Errorf("erreur suppression réaction: %v", err)
		}

		// Mettre à jour les comptes
		err = s.reactionRepo.UpdateThreadCounts(threadID)
		if err != nil {
			log.Printf("⚠️ Erreur mise à jour comptes thread: %v", err)
		}

		log.Printf("✅ Réaction supprimée - User %d, Thread %d, Type %s", userID, threadID, reactionType)
		return nil, nil
	}

	// Si l'utilisateur a une réaction différente, on la supprime d'abord
	if existingReaction != nil && existingReaction.ReactionType != reactionType {
		err = s.reactionRepo.Delete(userID, threadIDPtr, nil, existingReaction.ReactionType)
		if err != nil {
			return nil, fmt.Errorf("erreur suppression ancienne réaction: %v", err)
		}
	}

	// Créer la nouvelle réaction
	newReaction := models.Reaction{
		UserID:       userID,
		ThreadID:     threadIDPtr,
		MessageID:    nil,
		ReactionType: reactionType,
		CreatedAt:    time.Now(),
	}

	createdReaction, err := s.reactionRepo.Create(newReaction)
	if err != nil {
		return nil, fmt.Errorf("erreur création réaction: %v", err)
	}

	// Mettre à jour les comptes
	err = s.reactionRepo.UpdateThreadCounts(threadID)
	if err != nil {
		log.Printf("⚠️ Erreur mise à jour comptes thread: %v", err)
	}

	log.Printf("✅ Réaction ajoutée - ID %d, User %d, Thread %d, Type %s", 
		createdReaction.ID, userID, threadID, reactionType)

	return createdReaction, nil
}

// ToggleMessageReaction ajoute ou supprime une réaction sur un message
func (s *ReactionService) ToggleMessageReaction(userID, messageID int, reactionType string) (*models.Reaction, error) {
	log.Printf("🔄 ToggleMessageReaction - User %d, Message %d, Type %s", userID, messageID, reactionType)

	// Valider le type de réaction
	if !isValidReactionType(reactionType) {
		return nil, fmt.Errorf("type de réaction invalide: %s", reactionType)
	}

	messageIDPtr := &messageID

	// Vérifier si l'utilisateur a déjà une réaction sur ce message
	existingReaction, err := s.reactionRepo.GetUserReaction(userID, nil, messageIDPtr)
	if err != nil {
		return nil, fmt.Errorf("erreur vérification réaction existante: %v", err)
	}

	// Si l'utilisateur a déjà la même réaction, on la supprime (toggle off)
	if existingReaction != nil && existingReaction.ReactionType == reactionType {
		err = s.reactionRepo.Delete(userID, nil, messageIDPtr, reactionType)
		if err != nil {
			return nil, fmt.Errorf("erreur suppression réaction: %v", err)
		}

		// Mettre à jour les comptes
		err = s.reactionRepo.UpdateMessageCounts(messageID)
		if err != nil {
			log.Printf("⚠️ Erreur mise à jour comptes message: %v", err)
		}

		log.Printf("✅ Réaction supprimée - User %d, Message %d, Type %s", userID, messageID, reactionType)
		return nil, nil
	}

	// Si l'utilisateur a une réaction différente, on la supprime d'abord
	if existingReaction != nil && existingReaction.ReactionType != reactionType {
		err = s.reactionRepo.Delete(userID, nil, messageIDPtr, existingReaction.ReactionType)
		if err != nil {
			return nil, fmt.Errorf("erreur suppression ancienne réaction: %v", err)
		}
	}

	// Créer la nouvelle réaction
	newReaction := models.Reaction{
		UserID:       userID,
		ThreadID:     nil,
		MessageID:    messageIDPtr,
		ReactionType: reactionType,
		CreatedAt:    time.Now(),
	}

	createdReaction, err := s.reactionRepo.Create(newReaction)
	if err != nil {
		return nil, fmt.Errorf("erreur création réaction: %v", err)
	}

	// Mettre à jour les comptes
	err = s.reactionRepo.UpdateMessageCounts(messageID)
	if err != nil {
		log.Printf("⚠️ Erreur mise à jour comptes message: %v", err)
	}

	log.Printf("✅ Réaction ajoutée - ID %d, User %d, Message %d, Type %s", 
		createdReaction.ID, userID, messageID, reactionType)

	return createdReaction, nil
}

// GetThreadReactionCounts récupère les comptes de réactions pour un thread
func (s *ReactionService) GetThreadReactionCounts(threadID int) (map[string]int, error) {
	return s.reactionRepo.GetThreadReactionCounts(threadID)
}

// GetMessageReactionCounts récupère les comptes de réactions pour un message
func (s *ReactionService) GetMessageReactionCounts(messageID int) (map[string]int, error) {
	return s.reactionRepo.GetMessageReactionCounts(messageID)
}

// GetUserThreadReaction récupère la réaction d'un utilisateur sur un thread
func (s *ReactionService) GetUserThreadReaction(userID, threadID int) (*models.Reaction, error) {
	threadIDPtr := &threadID
	return s.reactionRepo.GetUserReaction(userID, threadIDPtr, nil)
}

// GetUserMessageReaction récupère la réaction d'un utilisateur sur un message
func (s *ReactionService) GetUserMessageReaction(userID, messageID int) (*models.Reaction, error) {
	messageIDPtr := &messageID
	return s.reactionRepo.GetUserReaction(userID, nil, messageIDPtr)
}

// GetThreadReactions récupère toutes les réactions d'un thread
func (s *ReactionService) GetThreadReactions(threadID int) ([]models.Reaction, error) {
	return s.reactionRepo.GetReactionsByThread(threadID)
}

// ValidateReactionRequest valide une demande de réaction
func (s *ReactionService) ValidateReactionRequest(request models.ReactionRequest) error {
	// Vérifier le type de cible
	if request.TargetType != "thread" && request.TargetType != "message" {
		return fmt.Errorf("target_type doit être 'thread' ou 'message'")
	}

	// Vérifier l'ID de la cible
	if request.TargetID <= 0 {
		return fmt.Errorf("target_id doit être un entier positif")
	}

	// Vérifier le type de réaction
	if !isValidReactionType(request.ReactionType) {
		return fmt.Errorf("reaction_type invalide: %s", request.ReactionType)
	}

	return nil
}

// ProcessReaction traite une demande de réaction
func (s *ReactionService) ProcessReaction(userID int, request models.ReactionRequest) (*models.Reaction, error) {
	// Valider la demande
	err := s.ValidateReactionRequest(request)
	if err != nil {
		return nil, err
	}

	// Traiter selon le type de cible
	switch request.TargetType {
	case "thread":
		return s.ToggleThreadReaction(userID, request.TargetID, request.ReactionType)
	case "message":
		return s.ToggleMessageReaction(userID, request.TargetID, request.ReactionType)
	default:
		return nil, fmt.Errorf("type de cible non supporté: %s", request.TargetType)
	}
}

// GetReactionSummary récupère un résumé des réactions pour un thread avec la réaction de l'utilisateur
func (s *ReactionService) GetReactionSummary(userID, threadID int) (map[string]interface{}, error) {
	// Récupérer les comptes
	counts, err := s.GetThreadReactionCounts(threadID)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération comptes: %v", err)
	}

	// Récupérer la réaction de l'utilisateur
	var userReaction *string
	if userID > 0 {
		reaction, err := s.GetUserThreadReaction(userID, threadID)
		if err != nil {
			log.Printf("⚠️ Erreur récupération réaction utilisateur: %v", err)
		} else if reaction != nil {
			userReaction = &reaction.ReactionType
		}
	}

	// Construire le résumé
	summary := map[string]interface{}{
		"counts":        counts,
		"user_reaction": userReaction,
		"total":         0,
	}

	// Calculer le total
	total := 0
	for _, count := range counts {
		total += count
	}
	summary["total"] = total

	return summary, nil
}

// isValidReactionType vérifie si un type de réaction est valide
func isValidReactionType(reactionType string) bool {
	validTypes := []string{"like", "dislike", "love", "laugh", "wow", "sad", "angry", "repost"}
	for _, validType := range validTypes {
		if reactionType == validType {
			return true
		}
	}
	return false
}

// GetReactionEmoji retourne l'emoji correspondant au type de réaction
func GetReactionEmoji(reactionType string) string {
	switch reactionType {
	case "like":
		return "👍"
	case "dislike":
		return "👎"
	case "love":
		return "❤️"
	case "laugh":
		return "😂"
	case "wow":
		return "😮"
	case "sad":
		return "😢"
	case "angry":
		return "😠"
	case "repost":
		return "🔄"
	default:
		return "👍"
	}
} 