package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	log.Println("🔄 Script de debug des avatars...")

	// Charger les variables d'environnement
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Fichier .env non trouvé, utilisation des variables d'environnement système")
	}

	// Récupérer les variables de connexion à la base de données
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "3306")
	dbUser := getEnvOrDefault("DB_USER", "root")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "")
	dbName := getEnvOrDefault("DB_NAME", "forum_y")

	// Construire la chaîne de connexion
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Se connecter à la base de données
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ Erreur connexion base de données: %v", err)
	}
	defer db.Close()

	// Tester la connexion
	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Erreur ping base de données: %v", err)
	}

	log.Println("✅ Connexion à la base de données réussie")

	// Afficher la structure de la table users
	log.Println("🔍 Structure de la table users:")
	rows, err := db.Query("DESCRIBE users")
	if err != nil {
		log.Fatalf("❌ Erreur DESCRIBE: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var field, fieldType, null, key, defaultValue, extra string
		err := rows.Scan(&field, &fieldType, &null, &key, &defaultValue, &extra)
		if err != nil {
			log.Printf("❌ Erreur scan structure: %v", err)
			continue
		}
		log.Printf("   %s (%s)", field, fieldType)
	}

	// Afficher les valeurs des profile_picture
	log.Println("🔍 Valeurs actuelles des profile_picture:")
	rows2, err := db.Query("SELECT id_user, username, profile_picture FROM users")
	if err != nil {
		log.Fatalf("❌ Erreur requête: %v", err)
	}
	defer rows2.Close()

	for rows2.Next() {
		var idUser int
		var username string
		var profilePicture *string
		
		err := rows2.Scan(&idUser, &username, &profilePicture)
		if err != nil {
			log.Printf("❌ Erreur scan: %v", err)
			continue
		}
		
		if profilePicture != nil {
			log.Printf("   ID:%d, Username:%s, ProfilePicture:%s", idUser, username, *profilePicture)
		} else {
			log.Printf("   ID:%d, Username:%s, ProfilePicture:NULL", idUser, username)
		}
	}
}

// getEnvOrDefault récupère une variable d'environnement ou une valeur par défaut
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 