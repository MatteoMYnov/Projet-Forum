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
		limit = 10 // Limite par défaut
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
		limit = 10
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
		limit = 10 // Limite par défaut
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

// GetThreadsStatistics récupère les statistiques des threads
func (s *ThreadService) GetThreadsStatistics() (map[string]int, error) {
	stats := make(map[string]int)

	// Total des threads
	total, err := s.threadRepo.GetTotalCount()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération total threads: %v", err)
	}
	stats["total"] = total

	// Threads créés aujourd'hui
	today, err := s.threadRepo.GetTodayThreadsCount()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération threads aujourd'hui: %v", err)
	}
	stats["today"] = today

	// Threads créés cette semaine
	week, err := s.threadRepo.GetWeekThreadsCount()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération threads cette semaine: %v", err)
	}
	stats["week"] = week

	return stats, nil
}

// GetTrendingThreads récupère les threads trending triés par likes/reactions
func (s *ThreadService) GetTrendingThreads(limit int) ([]models.Thread, error) {
	if limit < 1 || limit > 10 {
		limit = 5 // Limite par défaut pour trending
	}

	return s.threadRepo.GetTrendingThreads(limit)
}

// ChangeThreadStatus change l'état d'un thread (open, closed, archived)
func (s *ThreadService) ChangeThreadStatus(threadID int, newStatus string, userID int, isAdmin bool) error {
	// Valider le statut
	validStatuses := map[string]bool{
		"open":     true,
		"closed":   true,
		"archived": true,
	}
	
	if !validStatuses[newStatus] {
		return fmt.Errorf("statut invalide: %s. Statuts valides: open, closed, archived", newStatus)
	}

	// Récupérer le thread pour vérifier l'auteur
	thread, err := s.threadRepo.GetByID(threadID)
	if err != nil {
		return err
	}

	// Vérifier les permissions (seul l'auteur ou un admin peut modifier l'état)
	if thread.AuthorID != userID && !isAdmin {
		return fmt.Errorf("vous n'avez pas l'autorisation de modifier l'état de ce thread")
	}

	// Mettre à jour le statut
	err = s.threadRepo.UpdateStatus(threadID, newStatus)
	if err != nil {
		return fmt.Errorf("erreur lors de la mise à jour du statut: %v", err)
	}

	log.Printf("✅ Statut du thread %d changé vers '%s' par utilisateur %d", threadID, newStatus, userID)
	return nil
}

// CloseThread ferme un thread (reste visible mais plus de nouveaux messages)
func (s *ThreadService) CloseThread(threadID, userID int, isAdmin bool) error {
	return s.ChangeThreadStatus(threadID, "closed", userID, isAdmin)
}

// ArchiveThread archive un thread (plus visible dans les listes)
func (s *ThreadService) ArchiveThread(threadID, userID int, isAdmin bool) error {
	return s.ChangeThreadStatus(threadID, "archived", userID, isAdmin)
}

// ReopenThread réouvre un thread fermé ou archivé
func (s *ThreadService) ReopenThread(threadID, userID int, isAdmin bool) error {
	return s.ChangeThreadStatus(threadID, "open", userID, isAdmin)
}

// GetVisibleThreadsWithPagination récupère les threads visibles (non archivés) avec pagination
func (s *ThreadService) GetVisibleThreadsWithPagination(page, limit int) ([]models.Thread, *models.Meta, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10 // Limite par défaut
	}

	// Récupérer le total des threads visibles (non archivés)
	totalCount, err := s.threadRepo.GetVisibleThreadsCount()
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

	// Récupérer les threads visibles
	offset := (page - 1) * limit
	threads, err := s.threadRepo.GetVisibleThreads(limit, offset)
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

// CanPostMessage vérifie si un utilisateur peut poster un message dans un thread
func (s *ThreadService) CanPostMessage(threadID int) (bool, error) {
	thread, err := s.threadRepo.GetByID(threadID)
	if err != nil {
		return false, err
	}

	// Un message peut être posté seulement si le thread est ouvert
	return thread.Status == "open", nil
}

// GetThreadsByStatus récupère les threads filtrés par statut avec pagination
func (s *ThreadService) GetThreadsByStatus(status string, page, limit int) ([]models.Thread, *models.Meta, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Valider le statut
	validStatuses := map[string]bool{
		"open":     true,
		"closed":   true,
		"archived": true,
	}
	
	if !validStatuses[status] {
		return nil, nil, fmt.Errorf("statut invalide: %s", status)
	}

	// Récupérer le total pour ce statut
	totalCount, err := s.threadRepo.GetCountByStatus(status)
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
	threads, err := s.threadRepo.GetByStatus(status, limit, offset)
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

// UpdateThreadTitle met à jour le titre d'un thread (pour le créateur uniquement)
func (s *ThreadService) UpdateThreadTitle(threadID int, newTitle string, userID int) error {
	// Validation du titre
	newTitle = strings.TrimSpace(newTitle)
	if newTitle == "" {
		return fmt.Errorf("le titre ne peut pas être vide")
	}
	if len(newTitle) < 3 {
		return fmt.Errorf("le titre doit contenir au moins 3 caractères")
	}
	if len(newTitle) > 200 {
		return fmt.Errorf("le titre ne peut pas dépasser 200 caractères")
	}

	// Vérifier que le thread existe et que l'utilisateur en est le créateur
	thread, err := s.threadRepo.GetByID(threadID)
	if err != nil {
		return fmt.Errorf("thread non trouvé: %w", err)
	}

	if thread.AuthorID != userID {
		return fmt.Errorf("seul le créateur peut modifier le titre du thread")
	}

	// Mettre à jour le titre
	err = s.threadRepo.UpdateTitle(threadID, newTitle)
	if err != nil {
		return fmt.Errorf("erreur lors de la mise à jour: %w", err)
	}

	log.Printf("✅ Titre du thread %d mis à jour par l'utilisateur %d", threadID, userID)
	return nil
} 