package services

import (
	"database/sql"
	"fmt"
	"forum/models"
	"forum/repositories"
	"log"
	"strings"
)

// MessageService gère la logique métier des messages
type MessageService struct {
	messageRepo *repositories.MessageRepository
	threadRepo  *repositories.ThreadRepository
}

// NewMessageService crée une nouvelle instance du service
func NewMessageService(db *sql.DB) *MessageService {
	return &MessageService{
		messageRepo: repositories.NewMessageRepository(db),
		threadRepo:  repositories.NewThreadRepository(db),
	}
}

// CreateMessage crée un nouveau message
func (s *MessageService) CreateMessage(request models.MessageCreateRequest, authorID int) (*models.Message, error) {
	// Validation
	if err := s.validateMessageRequest(request); err != nil {
		return nil, err
	}

	// Vérifier que le thread existe
	_, err := s.threadRepo.GetByID(request.ThreadID)
	if err != nil {
		return nil, fmt.Errorf("thread non trouvé: %v", err)
	}

	// Créer le message
	message := &models.Message{
		ThreadID: request.ThreadID,
		AuthorID: authorID,
		Content:  strings.TrimSpace(request.Content),
	}

	// Ajouter le parent si spécifié
	if request.ParentMessageID != nil && *request.ParentMessageID > 0 {
		message.ParentMessageID = request.ParentMessageID
	}

	// Créer le message dans la base de données
	createdMessage, err := s.messageRepo.Create(message)
	if err != nil {
		return nil, err
	}

	// Mettre à jour le nombre de messages du thread
	err = s.threadRepo.UpdateMessageCount(request.ThreadID)
	if err != nil {
		log.Printf("⚠️ Erreur mise à jour nombre messages thread %d: %v", request.ThreadID, err)
	}

	log.Printf("✅ Message créé par utilisateur %d dans thread %d", authorID, request.ThreadID)
	return createdMessage, nil
}

// GetMessagesByThread récupère tous les messages d'un thread
func (s *MessageService) GetMessagesByThread(threadID int) ([]models.Message, error) {
	return s.messageRepo.GetByThreadID(threadID)
}

// GetMessage récupère un message par son ID
func (s *MessageService) GetMessage(messageID int) (*models.Message, error) {
	return s.messageRepo.GetByID(messageID)
}

// DeleteMessage supprime un message (seul l'auteur ou un admin peut le faire)
func (s *MessageService) DeleteMessage(messageID, userID int, isAdmin bool) error {
	// Récupérer le message pour vérifier l'auteur
	message, err := s.messageRepo.GetByID(messageID)
	if err != nil {
		return err
	}

	// Vérifier les permissions
	if message.AuthorID != userID && !isAdmin {
		return fmt.Errorf("vous n'avez pas l'autorisation de supprimer ce message")
	}

	err = s.messageRepo.Delete(messageID)
	if err != nil {
		return err
	}

	// Mettre à jour le nombre de messages du thread
	err = s.threadRepo.UpdateMessageCount(message.ThreadID)
	if err != nil {
		log.Printf("⚠️ Erreur mise à jour nombre messages thread %d: %v", message.ThreadID, err)
	}

	return nil
}

// validateMessageRequest valide les données d'une demande de création de message
func (s *MessageService) validateMessageRequest(request models.MessageCreateRequest) error {
	// Validation du contenu
	content := strings.TrimSpace(request.Content)
	if content == "" {
		return fmt.Errorf("le contenu est requis")
	}
	if len(content) > 2000 {
		return fmt.Errorf("le contenu ne peut pas dépasser 2000 caractères")
	}

	// Validation du thread ID
	if request.ThreadID <= 0 {
		return fmt.Errorf("ID de thread invalide")
	}

	return nil
}

// DeleteMessageByThreadOwner permet au créateur d'un thread de supprimer n'importe quel message de son thread
func (s *MessageService) DeleteMessageByThreadOwner(messageID, threadID, threadOwnerID int) error {
	// Vérifier que le message existe et appartient au thread
	message, err := s.messageRepo.GetByID(messageID)
	if err != nil {
		return fmt.Errorf("message non trouvé: %w", err)
	}

	if message.ThreadID != threadID {
		return fmt.Errorf("le message n'appartient pas à ce thread")
	}

	// Supprimer le message
	err = s.messageRepo.Delete(messageID)
	if err != nil {
		return fmt.Errorf("erreur lors de la suppression: %w", err)
	}

	// Mettre à jour le nombre de messages du thread
	err = s.threadRepo.UpdateMessageCount(threadID)
	if err != nil {
		log.Printf("⚠️ Erreur mise à jour nombre messages thread %d: %v", threadID, err)
	}

	log.Printf("✅ Message %d supprimé par le créateur du thread %d", messageID, threadID)
	return nil
}

// DeleteMultipleMessagesByThreadOwner supprime plusieurs messages en lot pour le créateur d'un thread
func (s *MessageService) DeleteMultipleMessagesByThreadOwner(messageIDs []int, threadID, threadOwnerID int) error {
	if len(messageIDs) == 0 {
		return fmt.Errorf("aucun message à supprimer")
	}

	deletedCount := 0
	for _, messageID := range messageIDs {
		err := s.DeleteMessageByThreadOwner(messageID, threadID, threadOwnerID)
		if err != nil {
			log.Printf("⚠️ Erreur suppression message %d: %v", messageID, err)
			continue
		}
		deletedCount++
	}

	if deletedCount == 0 {
		return fmt.Errorf("aucun message n'a pu être supprimé")
	}

	log.Printf("✅ %d/%d messages supprimés par le créateur du thread %d", deletedCount, len(messageIDs), threadID)
	return nil
} 