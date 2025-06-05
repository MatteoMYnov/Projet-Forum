package config

import (
	"log"
	"os"
	"strconv"
)

// Variables globales pour un accès facile
var (
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
	JWTSecret  string
	Debug      bool
	UploadPath string
)

// LoadEnv charge les variables d'environnement
func LoadEnv() {
	log.Println("🔧 Chargement des variables d'environnement...")

	// Base de données
	DBHost = getEnvWithDefault("DB_HOST", "localhost")
	DBPort = getEnvAsInt("DB_PORT", 3306)
	DBUser = getEnvWithDefault("DB_USER", "root")
	DBPassword = getEnvWithDefault("DB_PASSWORD", "")
	DBName = getEnvWithDefault("DB_NAME", "forum_y")

	// Serveur
	ServerPort = getEnvWithDefault("PORT", "2557")

	// JWT
	JWTSecret = getEnvWithDefault("JWT_SECRET", "default_secret_change_me")

	// Autres
	Debug = getEnvAsBool("DEBUG", true)
	UploadPath = getEnvWithDefault("UPLOAD_PATH", "./uploads")

	// Validation des variables critiques
	validateRequiredVars()

	log.Println("✅ Variables d'environnement chargées avec succès")
}

// getEnvWithDefault récupère une variable avec valeur par défaut
func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("⚠️ Variable %s non définie, utilisation de la valeur par défaut: %s", key, defaultValue)
		return defaultValue
	}
	return value
}

// getEnvAsInt récupère une variable d'environnement en tant qu'entier
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		log.Printf("⚠️ Variable %s non définie, utilisation de la valeur par défaut: %d", key, defaultValue)
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("❌ Erreur conversion %s en entier: %v, utilisation valeur par défaut: %d", key, err, defaultValue)
		return defaultValue
	}
	return value
}

// getEnvAsBool récupère une variable d'environnement en tant que booléen
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		log.Printf("❌ Erreur conversion %s en booléen: %v, utilisation valeur par défaut: %t", key, err, defaultValue)
		return defaultValue
	}
	return value
}

// validateRequiredVars vérifie que les variables obligatoires sont présentes
func validateRequiredVars() {
	errors := []string{}

	if DBHost == "" {
		errors = append(errors, "DB_HOST")
	}
	if DBUser == "" {
		errors = append(errors, "DB_USER")
	}
	if DBName == "" {
		errors = append(errors, "DB_NAME")
	}
	if JWTSecret == "default_secret_change_me" {
		log.Println("⚠️ ATTENTION: Vous utilisez le secret JWT par défaut, changez-le en production!")
	}

	if len(errors) > 0 {
		log.Fatalf("❌ Variables d'environnement manquantes: %v", errors)
	}
}
