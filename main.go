package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// Sert les fichiers statiques d'un dossier donné à une route précise
func setupFileServer(path, route string) {
	fs := http.FileServer(http.Dir(path))
	http.Handle(route, http.StripPrefix(route, fs))
}

// Redirige "/" vers "/home"
func rootRedirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", http.StatusFound)
}

// Sert les pages HTML
func serveHTML(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	}
}

func main() {
	// Connexion à la base de données
	dsn := "root:@tcp(localhost:3306)/forum_y"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erreur lors de l'ouverture :", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Connexion impossible :", err)
	}
	fmt.Println("Connexion réussie à la base de données 🎉")

	// Redirection de la racine vers /home
	http.HandleFunc("/", rootRedirectHandler)

	// Pages HTML
	http.HandleFunc("/home", serveHTML("./website/template/home.html"))
	http.HandleFunc("/profile", serveHTML("./website/template/profile.html"))
	http.HandleFunc("/register", serveHTML("./website/template/register.html"))
	http.HandleFunc("/login", serveHTML("./website/template/login.html"))

	// Fichiers statiques
	setupFileServer("./website/styles", "/styles/")
	setupFileServer("./website/img", "/img/")
	setupFileServer("./website/js", "/js/")

	// Routes API
	http.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}
		pseudo := r.FormValue("pseudo")
		email := r.FormValue("email")
		password := r.FormValue("password")

		if pseudo == "" || email == "" || password == "" {
			http.Error(w, "Champs requis manquants", http.StatusBadRequest)
			return
		}

		_, err := db.Exec("INSERT INTO users (pseudo, email, password) VALUES (?, ?, ?)", pseudo, email, password)
		if err != nil {
			http.Error(w, "Erreur lors de l'inscription", http.StatusInternalServerError)
			log.Println("Erreur INSERT:", err)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})

	http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}
		identifiant := r.FormValue("identifiant")
		password := r.FormValue("password")

		if identifiant == "" || password == "" {
			http.Error(w, "Identifiants manquants", http.StatusBadRequest)
			return
		}

		var dbPassword string
		query := "SELECT password FROM users WHERE pseudo = ? OR email = ? LIMIT 1"
		err := db.QueryRow(query, identifiant, identifiant).Scan(&dbPassword)
		if err != nil {
			http.Error(w, "Utilisateur non trouvé ou mot de passe incorrect", http.StatusUnauthorized)
			log.Println("Erreur SELECT:", err)
			return
		}

		if strings.TrimSpace(password) != strings.TrimSpace(dbPassword) {
			http.Error(w, "Mot de passe incorrect", http.StatusUnauthorized)
			return
		}

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	})

	// Lancement du serveur
	if err := http.ListenAndServe(":2557", nil); err != nil {
		log.Fatalf("Erreur lors du démarrage du serveur: %v", err)
	}
}
	