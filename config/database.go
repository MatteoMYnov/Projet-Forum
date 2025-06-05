package config

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DbContext variable globale pour la connexion
var DbContext *sql.DB

// InitDB initialise la connexion à la base de données
func InitDB() {
	log.Println("🔌 Initialisation de la connexion à la base de données...")

	// Construction de la chaîne de connexion (DSN)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=30s",
		DBUser, DBPassword, DBHost, DBPort, DBName)

	var err error
	DbContext, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ Erreur ouverture connexion BDD: %v", err)
	}

	// Configuration du pool de connexions
	DbContext.SetMaxOpenConns(25)                 // Max 25 connexions ouvertes
	DbContext.SetMaxIdleConns(5)                  // Max 5 connexions inactives
	DbContext.SetConnMaxLifetime(5 * time.Minute) // Durée de vie max: 5min

	// Test de la connexion
	if err = pingDatabase(3); err != nil {
		DbContext.Close()
		log.Fatalf("❌ Impossible de se connecter à la BDD: %v", err)
	}

	log.Printf("✅ Connexion à la base de données '%s' réussie!", DBName)

	// Vérifier que les tables existent
	if err = checkTables(); err != nil {
		log.Printf("⚠️ Problème avec les tables: %v", err)
	}
}

// pingDatabase teste la connexion avec retry
func pingDatabase(retries int) error {
	for i := 0; i < retries; i++ {
		err := DbContext.Ping()
		if err == nil {
			return nil
		}

		log.Printf("⚠️ Tentative de connexion %d/%d échouée: %v", i+1, retries, err)
		if i < retries-1 {
			time.Sleep(2 * time.Second)
		}
	}
	return fmt.Errorf("échec de connexion après %d tentatives", retries)
}

// checkTables vérifie que les tables principales existent
func checkTables() error {
	tables := []string{"users", "threads", "messages", "categories", "reactions"}

	for _, table := range tables {
		var exists int
		query := `SELECT COUNT(*) FROM information_schema.tables 
				  WHERE table_schema = ? AND table_name = ?`

		err := DbContext.QueryRow(query, DBName, table).Scan(&exists)
		if err != nil {
			return fmt.Errorf("erreur vérification table %s: %v", table, err)
		}

		if exists == 0 {
			log.Printf("⚠️ Table '%s' n'existe pas", table)
		} else {
			log.Printf("✅ Table '%s' trouvée", table)
		}
	}
	return nil
}

// GetDB retourne l'instance de la base de données
func GetDB() *sql.DB {
	if DbContext == nil {
		log.Fatal("❌ Base de données non initialisée! Appelez InitDB() d'abord.")
	}
	return DbContext
}

// CloseDB ferme proprement la connexion
func CloseDB() {
	if DbContext != nil {
		err := DbContext.Close()
		if err != nil {
			log.Printf("❌ Erreur fermeture BDD: %v", err)
		} else {
			log.Println("🔒 Connexion BDD fermée proprement")
		}
		DbContext = nil
	}
}

// TestConnection teste la connexion et affiche les infos
func TestConnection() error {
	log.Println("🧪 Test de la connexion...")

	// Test simple
	var version string
	err := DbContext.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		return fmt.Errorf("erreur test connexion: %v", err)
	}

	log.Printf("✅ MySQL Version: %s", version)

	// Statistiques de connexion
	stats := DbContext.Stats()
	log.Printf("📊 Stats connexions - Ouvertes: %d, En cours: %d, Inactives: %d",
		stats.OpenConnections, stats.InUse, stats.Idle)

	return nil
}
