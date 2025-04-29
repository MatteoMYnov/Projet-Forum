package main

import (
	"log"
	"net/http"
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

func main() {
	// Redirection de la racine vers /home
	http.HandleFunc("/", rootRedirectHandler)

	// Page Home
	http.HandleFunc("/home", homeHandler)

	// Fichiers statiques
	setupFileServer("./website/styles", "/styles/")
	setupFileServer("./website/img", "/img/")
	setupFileServer("./website/js", "/js/")

	// Lancement du serveur
	if err := http.ListenAndServe(":2556", nil); err != nil {
		log.Fatalf("Erreur lors du démarrage du serveur: %v", err)
	}
}
