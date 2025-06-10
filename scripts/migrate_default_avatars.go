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
	log.Println("🔄 Migration des avatars par défaut...")

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

	// Debug: afficher les valeurs actuelles
	rows, err := db.Query("SELECT id_user, username, profile_picture FROM users")
	if err != nil {
		log.Fatalf("❌ Erreur requête debug: %v", err)
	}
	defer rows.Close()

	log.Println("🔍 Valeurs actuelles des profile_picture:")
	for rows.Next() {
		var idUsers int
		var username string
		var profilePicture *string
		
		err := rows.Scan(&idUsers, &username, &profilePicture)
		if err != nil {
			log.Printf("❌ Erreur scan: %v", err)
			continue
		}
		
		if profilePicture != nil {
			log.Printf("   ID:%d, Username:%s, ProfilePicture:%s", idUsers, username, *profilePicture)
		} else {
			log.Printf("   ID:%d, Username:%s, ProfilePicture:NULL", idUsers, username)
		}
	}

	// Compter les utilisateurs sans photo de profil
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE profile_picture IS NULL OR profile_picture = '' OR profile_picture = '/img/default-avatar.png'").Scan(&count)
	if err != nil {
		log.Fatalf("❌ Erreur comptage utilisateurs: %v", err)
	}

	log.Printf("📊 %d utilisateurs avec chemin avatar incorrect trouvés", count)

	if count == 0 {
		log.Println("✅ Aucune migration nécessaire")
		return
	}

	// Mettre à jour les utilisateurs sans photo de profil OU avec le mauvais chemin
	defaultAvatar := "/img/avatars/default-avatar.png"
	result, err := db.Exec("UPDATE users SET profile_picture = ? WHERE profile_picture IS NULL OR profile_picture = '' OR profile_picture = '/img/default-avatar.png'", defaultAvatar)
	if err != nil {
		log.Fatalf("❌ Erreur mise à jour: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("❌ Erreur récupération nombre de lignes: %v", err)
	}

	log.Printf("✅ Migration terminée: %d utilisateurs mis à jour avec l'avatar par défaut", rowsAffected)
}

// getEnvOrDefault récupère une variable d'environnement ou une valeur par défaut
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 