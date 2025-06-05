package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// DbContext variable globale pour la connexion
var DbContext *sql.DB

// InitDB initialise la connexion à la base de données
func InitDB() {
	// Récupération des paramètres
	host := GetEnvWithDefault("DB_HOST", "localhost")
	port := GetEnvAsInt("DB_PORT", 3306)
	user := GetEnvWithDefault("DB_USER", "root")
	password := GetEnvWithDefault("DB_PASSWORD", "")
	dbname := GetEnvWithDefault("DB_NAME", "forum_y")

	// Construction de la chaîne de connexion
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	var err error
	DbContext, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Erreur ouverture connexion BDD: %v", err)
	}

	// Test de la connexion
	err = DbContext.Ping()
	if err != nil {
		DbContext.Close()
		log.Fatalf("Impossible de se connecter à la BDD: %v", err)
	}

	// Configuration du pool de connexions
	DbContext.SetMaxOpenConns(25)
	DbContext.SetMaxIdleConns(5)

	fmt.Println("✅ Connexion à la base de données réussie!")
}

// CloseDB ferme la connexion à la base de données
func CloseDB() {
	if DbContext != nil {
		DbContext.Close()
		fmt.Println("🔒 Connexion BDD fermée")
	}
}
