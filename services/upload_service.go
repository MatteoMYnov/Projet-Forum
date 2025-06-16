package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
	"crypto/rand"
	"encoding/hex"
)

// UploadService gère les uploads de fichiers
type UploadService struct {
	uploadDir string
	maxSize   int64
}

// NewUploadService crée une nouvelle instance du service
func NewUploadService(uploadDir string, maxSize int64) *UploadService {
	// Créer le dossier s'il n'existe pas
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(fmt.Sprintf("Impossible de créer le dossier d'upload: %v", err))
	}
	
	return &UploadService{
		uploadDir: uploadDir,
		maxSize:   maxSize,
	}
}

// UploadProfilePicture télécharge et sauvegarde une image de profil
func (s *UploadService) UploadProfilePicture(file multipart.File, header *multipart.FileHeader) (string, error) {
	// Validation de la taille
	if header.Size > s.maxSize {
		return "", fmt.Errorf("fichier trop volumineux (max %d MB)", s.maxSize/(1024*1024))
	}
	
	// Validation du type de fichier
	if !s.isValidImageType(header.Filename) {
		return "", fmt.Errorf("type de fichier non supporté")
	}
	
	// Générer un nom de fichier unique
	filename, err := s.generateUniqueFilename(header.Filename)
	if err != nil {
		return "", fmt.Errorf("erreur génération nom de fichier: %v", err)
	}
	
	// Créer le chemin complet
	fullPath := filepath.Join(s.uploadDir, filename)
	
	// Créer le fichier de destination
	destFile, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("erreur création fichier: %v", err)
	}
	defer destFile.Close()
	
	// Copier le contenu
	_, err = io.Copy(destFile, file)
	if err != nil {
		os.Remove(fullPath) // Nettoyer en cas d'erreur
		return "", fmt.Errorf("erreur sauvegarde fichier: %v", err)
	}
	
	// Retourner le chemin relatif pour la base de données
	return "/img/avatars/" + filename, nil
}

// isValidImageType vérifie si le fichier est une image valide
func (s *UploadService) isValidImageType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	
	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}
	return false
}

// generateUniqueFilename génère un nom de fichier unique
func (s *UploadService) generateUniqueFilename(originalFilename string) (string, error) {
	// Extraire l'extension
	ext := filepath.Ext(originalFilename)
	
	// Générer un identifiant unique
	randBytes := make([]byte, 16)
	_, err := rand.Read(randBytes)
	if err != nil {
		return "", err
	}
	
	uniqueID := hex.EncodeToString(randBytes)
	timestamp := time.Now().Unix()
	
	// Construire le nom de fichier: timestamp_uniqueID.ext
	filename := fmt.Sprintf("%d_%s%s", timestamp, uniqueID, ext)
	
	return filename, nil
}

// DeleteProfilePicture supprime une image de profil
func (s *UploadService) DeleteProfilePicture(profilePicturePath string) error {
	if profilePicturePath == "" || profilePicturePath == "/img/avatars/default-avatar.png" {
		return nil // Ne pas supprimer l'image par défaut
	}
	
	// Construire le chemin complet
	filename := filepath.Base(profilePicturePath)
	fullPath := filepath.Join(s.uploadDir, filename)
	
	// Supprimer le fichier s'il existe
	if _, err := os.Stat(fullPath); err == nil {
		return os.Remove(fullPath)
	}
	
	return nil // Fichier n'existe pas, pas d'erreur
}

// GetDefaultAvatarPath retourne le chemin vers l'avatar par défaut
func (s *UploadService) GetDefaultAvatarPath() string {
	return "/img/avatars/default-avatar.png"
}

// UploadBanner télécharge et sauvegarde une bannière de profil
func (s *UploadService) UploadBanner(file multipart.File, header *multipart.FileHeader) (string, error) {
	// Validation de la taille
	if header.Size > s.maxSize {
		return "", fmt.Errorf("fichier trop volumineux (max %d MB)", s.maxSize/(1024*1024))
	}
	
	// Validation du type de fichier
	if !s.isValidImageType(header.Filename) {
		return "", fmt.Errorf("type de fichier non supporté")
	}
	
	// Générer un nom de fichier unique
	filename, err := s.generateUniqueFilename(header.Filename)
	if err != nil {
		return "", fmt.Errorf("erreur génération nom de fichier: %v", err)
	}
	
	// Créer le dossier banners s'il n'existe pas
	bannersDir := filepath.Join(filepath.Dir(s.uploadDir), "banners")
	if err := os.MkdirAll(bannersDir, 0755); err != nil {
		return "", fmt.Errorf("impossible de créer le dossier banners: %v", err)
	}
	
	// Créer le chemin complet
	fullPath := filepath.Join(bannersDir, filename)
	
	// Créer le fichier de destination
	destFile, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("erreur création fichier: %v", err)
	}
	defer destFile.Close()
	
	// Copier le contenu
	_, err = io.Copy(destFile, file)
	if err != nil {
		os.Remove(fullPath) // Nettoyer en cas d'erreur
		return "", fmt.Errorf("erreur sauvegarde fichier: %v", err)
	}
	
	// Retourner le chemin relatif pour la base de données
	return "/img/banners/" + filename, nil
}

// DeleteBanner supprime une bannière
func (s *UploadService) DeleteBanner(bannerPath string) error {
	if bannerPath == "" || bannerPath == "/img/banners/default-avatar.png" {
		return nil // Ne pas supprimer la bannière par défaut
	}
	
	// Construire le chemin complet
	filename := filepath.Base(bannerPath)
	bannersDir := filepath.Join(filepath.Dir(s.uploadDir), "banners")
	fullPath := filepath.Join(bannersDir, filename)
	
	// Supprimer le fichier s'il existe
	if _, err := os.Stat(fullPath); err == nil {
		return os.Remove(fullPath)
	}
	
	return nil // Fichier n'existe pas, pas d'erreur
}

// GetDefaultBannerPath retourne le chemin vers la bannière par défaut
func (s *UploadService) GetDefaultBannerPath() string {
	return "/img/banners/default-avatar.png"
}

// SaveImage sauvegarde une image dans le dossier spécifié (générique pour avatars et bannières)
func (s *UploadService) SaveImage(file multipart.File, header *multipart.FileHeader, subfolder string) (string, error) {
	// Validation de la taille
	if header.Size > s.maxSize {
		return "", fmt.Errorf("fichier trop volumineux (max %d MB)", s.maxSize/(1024*1024))
	}
	
	// Validation du type de fichier
	if !s.isValidImageType(header.Filename) {
		return "", fmt.Errorf("type de fichier non supporté")
	}
	
	// Générer un nom de fichier unique
	filename, err := s.generateUniqueFilename(header.Filename)
	if err != nil {
		return "", fmt.Errorf("erreur génération nom de fichier: %v", err)
	}
	
	// Créer le dossier s'il n'existe pas
	targetDir := filepath.Join("./website/img", subfolder)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("impossible de créer le dossier %s: %v", subfolder, err)
	}
	
	// Créer le chemin complet
	fullPath := filepath.Join(targetDir, filename)
	
	// Créer le fichier de destination
	destFile, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("erreur création fichier: %v", err)
	}
	defer destFile.Close()
	
	// Copier le contenu
	_, err = io.Copy(destFile, file)
	if err != nil {
		os.Remove(fullPath) // Nettoyer en cas d'erreur
		return "", fmt.Errorf("erreur sauvegarde fichier: %v", err)
	}
	
	// Retourner le chemin relatif pour la base de données
	return "/img/" + subfolder + "/" + filename, nil
}