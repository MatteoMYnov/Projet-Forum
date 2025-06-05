package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// LoadEnv charge les variables d'environnement
func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Erreur chargement fichier .env: %v", err)
	}
}

// GetEnvWithDefault récupère une variable avec valeur par défaut
func GetEnvWithDefault(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}
	return value
}

// GetEnvAsInt récupère une variable d'environnement en tant qu'entier
func GetEnvAsInt(key string, defaultValue int) int {
	valueStr := GetEnvWithDefault(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Erreur conversion %s en entier, utilisation valeur par défaut: %d", key, defaultValue)
		return defaultValue
	}
	return value
}

// ValidateRequiredEnvVars vérifie que les variables obligatoires sont présentes
func ValidateRequiredEnvVars() {
	required := []string{"DB_HOST", "DB_USER", "DB_NAME"}

	for _, key := range required {
		if GetEnvWithDefault(key, "") == "" {
			log.Fatalf("Variable d'environnement manquante: %s", key)
		}
	}
}
