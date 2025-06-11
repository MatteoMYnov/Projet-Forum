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