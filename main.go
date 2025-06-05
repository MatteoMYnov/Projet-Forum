package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

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

// Sert la page d'accueil depuis /home
func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./website/template/home.html")
}

// Sert la page de profil depuis /profile
func profileHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./website/template/profile.html")
}

func main() {
	// Redirection de la racine vers /home
	http.HandleFunc("/", rootRedirectHandler)

	// Page Home
	http.HandleFunc("/profile", profileHandler)
	http.HandleFunc("/home", homeHandler)

	// Fichiers statiques
	setupFileServer("./website/styles", "/styles/")
	setupFileServer("./website/img", "/img/")
	setupFileServer("./website/js", "/js/")

	dsn := "root:@tcp(localhost:3306)/forum_y"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erreur lors de l'ouverture :", err)
	}
	defer db.Close()

	// Vérifie la connexion
	err = db.Ping()
	if err != nil {
		log.Fatal("Connexion impossible :", err)
	}

	fmt.Println("Connexion réussie à la base de données 🎉")

	// Lancement du serveur
	if err := http.ListenAndServe(":2557", nil); err != nil {
		log.Fatalf("Erreur lors du démarrage du serveur: %v", err)
	}

}
