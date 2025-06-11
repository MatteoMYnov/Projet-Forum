package services

import (
	"database/sql"
	"fmt"
	"forum/models"
	"forum/repositories"
	"log"
	"strings"
)

// ThreadService gère la logique métier des threads
type ThreadService struct {
	threadRepo *repositories.ThreadRepository
}

// NewThreadService crée une nouvelle instance du service
func NewThreadService(db *sql.DB) *ThreadService {
	return &ThreadService{
		threadRepo: repositories.NewThreadRepository(db),
	}
}

// CreateThread crée un nouveau thread
func (s *ThreadService) CreateThread(request models.ThreadCreateRequest, authorID int) (*models.Thread, error) {
	// Validation
	if err := s.validateThreadRequest(request); err != nil {
		return nil, err
	}

	// Créer le thread
	thread := &models.Thread{
		Title:    strings.TrimSpace(request.Title),
		Content:  strings.TrimSpace(request.Content),
		AuthorID: authorID,
	}

	// Ajouter la catégorie si spécifiée
	if request.CategoryID != nil && *request.CategoryID > 0 {
		thread.CategoryID = request.CategoryID
	}

	// Créer le thread dans la base de données
	createdThread, err := s.threadRepo.Create(thread)
	if err != nil {
		return nil, err
	}

	log.Printf("✅ Thread créé par utilisateur %d: %s", authorID, createdThread.Title)
	return createdThread, nil
}

// GetThread récupère un thread par son ID et incrémente les vues
func (s *ThreadService) GetThread(threadID int) (*models.Thread, error) {
	thread, err := s.threadRepo.GetByID(threadID)
	if err != nil {
		return nil, err
	}

	// Incrémenter le nombre de vues
	err = s.threadRepo.UpdateViewCount(threadID)
	if err != nil {
		log.Printf("⚠️ Erreur mise à jour vues pour thread %d: %v", threadID, err)
		// Ne pas faire échouer la requête pour ça
	}

	return thread, nil
}

// GetAllThreads récupère tous les threads avec pagination
func (s *ThreadService) GetAllThreads(page, limit int) ([]models.Thread, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20 // Limite par défaut
	}

	offset := (page - 1) * limit
	return s.threadRepo.GetAll(limit, offset)
}

// GetUserThreads récupère les threads d'un utilisateur
func (s *ThreadService) GetUserThreads(userID, page, limit int) ([]models.Thread, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	return s.threadRepo.GetByUserID(userID, limit, offset)
}

// GetCategories récupère toutes les catégories disponibles
func (s *ThreadService) GetCategories() ([]models.Category, error) {
	return s.threadRepo.GetCategories()
}

// DeleteThread supprime un thread (seul l'auteur ou un admin peut le faire)
func (s *ThreadService) DeleteThread(threadID, userID int, isAdmin bool) error {
	// Récupérer le thread pour vérifier l'auteur
	thread, err := s.threadRepo.GetByID(threadID)
	if err != nil {
		return err
	}

	// Vérifier les permissions
	if thread.AuthorID != userID && !isAdmin {
		return fmt.Errorf("vous n'avez pas l'autorisation de supprimer ce thread")
	}

	return s.threadRepo.Delete(threadID)
}

// validateThreadRequest valide les données d'une demande de création de thread
func (s *ThreadService) validateThreadRequest(request models.ThreadCreateRequest) error {
	// Validation du titre
	title := strings.TrimSpace(request.Title)
	if title == "" {
		return fmt.Errorf("le titre est requis")
	}
	if len(title) > 280 {
		return fmt.Errorf("le titre ne peut pas dépasser 280 caractères")
	}

	// Validation du contenu
	content := strings.TrimSpace(request.Content)
	if content == "" {
		return fmt.Errorf("le contenu est requis")
	}
	if len(content) > 5000 {
		return fmt.Errorf("le contenu ne peut pas dépasser 5000 caractères")
	}

	return nil
}

// ProcessHashtagsFromRequest traite les hashtags depuis la requête
func (s *ThreadService) ProcessHashtagsFromRequest(hashtagsInput string) []string {
	if hashtagsInput == "" {
		return []string{}
	}

	// Traiter la chaîne d'hashtags
	return repositories.ProcessHashtags(hashtagsInput)
}

// GetThreadsWithPagination récupère les threads avec métadonnées de pagination
func (s *ThreadService) GetThreadsWithPagination(page, limit int) ([]models.Thread, *models.Meta, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20 // Limite par défaut
	}

	// Récupérer le total
	totalCount, err := s.threadRepo.GetTotalCount()
	if err != nil {
		return nil, nil, err
	}

	// Calculer le nombre total de pages
	totalPages := (totalCount + limit - 1) / limit
	if totalPages < 1 {
		totalPages = 1
	}

	// Vérifier que la page demandée n'est pas au-delà du total
	if page > totalPages {
		page = totalPages
	}

	// Récupérer les threads
	offset := (page - 1) * limit
	threads, err := s.threadRepo.GetAll(limit, offset)
	if err != nil {
		return nil, nil, err
	}

	// Créer les métadonnées
	meta := &models.Meta{
		Page:       page,
		PerPage:    limit,
		TotalPages: totalPages,
		TotalCount: totalCount,
	}

	return threads, meta, nil
} 