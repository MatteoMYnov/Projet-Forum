package main

import (
	"forum/config"
	"forum/controllers"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	// 📝 Chargement du fichier .env
	log.Println("🚀 Démarrage du serveur Forum...")

	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Fichier .env non trouvé, utilisation des variables d'environnement système")
	}

	// 🔧 Chargement des variables d'environnement
	config.LoadEnv()

	// 🔌 Initialisation de la base de données
	config.InitDB()
	defer config.CloseDB()

	// 🧪 Test de connexion optionnel
	if config.Debug {
		if err := config.TestConnection(); err != nil {
			log.Printf("⚠️ Problème de test de connexion: %v", err)
		}
	}

	userController := controllers.NewUserControllers(config.DbContext) // Initialisation des contrôleurs

	// 📂 Serveur de fichiers statiques

	// 🌐 Routes principales (pages HTML)
	setupPageRoutes()

	/* // 🔗 Routes API (authentification et données)
	setupAPIRoutes()
	*/
	router := http.NewServeMux()
	setupStaticFiles(router)
	userController.UserRouter(router) // Enregistrement des routes du contrôleur

	// 🎯 Démarrage du serveur
	serverAddr := ":" + config.ServerPort
	log.Printf("✅ Serveur démarré sur http://localhost%s", serverAddr)
	log.Printf("📄 Pages disponibles:")
	log.Printf("   - Home: http://localhost%s/home", serverAddr)
	log.Printf("   - Theme: http://localhost%s/theme", serverAddr)
	log.Printf("   - Login: http://localhost%s/login", serverAddr)
	log.Printf("   - Register: http://localhost%s/register", serverAddr)
	log.Printf("   - Profile: http://localhost%s/profile", serverAddr)

	if err := http.ListenAndServe(serverAddr, router); err != nil {
		log.Fatalf("❌ Erreur serveur: %v", err)
	}
}

// setupStaticFiles configure le serveur de fichiers statiques
func setupStaticFiles(r *http.ServeMux) {
	// CSS, JS, Images
	r.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./website/styles/"))))
	r.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./website/js/"))))
	r.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./website/img/"))))
	log.Println("📂 Fichiers statiques configurés")
}

// setupPageRoutes configure les routes des pages HTML
func setupPageRoutes() {
	// Redirection racine
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/home", http.StatusFound)
			return
		}
		http.NotFound(w, r)
	})

	log.Println("🌐 Routes des pages configurées")
}

/* // setupAPIRoutes configure les routes API
func setupAPIRoutes() {
	// Routes d'authentification
	http.HandleFunc("/api/register", controllers.RegisterHandler)
	http.HandleFunc("/api/login", controllers.LoginHandler)
	http.HandleFunc("/api/logout", controllers.LogoutHandler)

	// Routes protégées
	http.HandleFunc("/api/profile", middleware.RequireAuth(controllers.ProfileAPI))

	log.Println("🔗 Routes API configurées")
} */
