package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"forum/middleware"
	"forum/models"
	"forum/services"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// UserRepository gère les opérations sur les utilisateurs
type UserControllers struct {
	authService    *services.AuthService
	threadService  *services.ThreadService
	uploadService  *services.UploadService
	reactionService *services.ReactionService
	messageService *services.MessageService
}

// NewUserControllers crée une nouvelle instance du controller
func NewUserControllers(db *sql.DB) *UserControllers {
	// Créer le service d'upload avec un dossier d'avatars et une taille max de 5MB
	uploadService := services.NewUploadService("./website/img/avatars", 5*1024*1024)
	
	return &UserControllers{
		authService:     services.NewAuthService(db),
		threadService:   services.NewThreadService(db),
		uploadService:   uploadService,
		reactionService: services.NewReactionService(db),
		messageService:  services.NewMessageService(db),
	}
}

func (c *UserControllers) UserRouter(r *http.ServeMux) {
	// Routes pour l'authentification
	r.HandleFunc("/register", c.RegisterPage)
	r.HandleFunc("/login", c.LoginPage)
	r.HandleFunc("/home", c.HomePage)
	r.HandleFunc("/theme", c.ThemePage)                  // Page de thème
	r.HandleFunc("/profile", middleware.RequireAuth(c.ProfilePage))
	
	// Routes pour les threads
	r.HandleFunc("/threads", c.ThreadsListPage)          // Liste des threads
	r.HandleFunc("/threads_demo", c.ThreadsDemoPage)     // Page de démo
	r.HandleFunc("/create-thread", middleware.RequireAuth(c.CreateThreadPage))
	r.HandleFunc("/thread/", c.ThreadPage)               // Pour afficher un thread spécifique

	// Handlers pour les actions
	r.HandleFunc("/api/register", c.RegisterHandler)
	r.HandleFunc("/api/login", c.LoginHandler)
	r.HandleFunc("/api/logout", c.LogoutHandler)
	r.HandleFunc("/api/profile", middleware.RequireAuth(c.ProfileAPI))
	
	// API pour les threads
	r.HandleFunc("/api/threads", middleware.RequireAuth(c.CreateThreadHandler))
	r.HandleFunc("/api/threads/", c.ThreadAPI) // Pour récupérer les données d'un thread
	
	// API pour les messages
	r.HandleFunc("/api/messages", middleware.RequireAuth(c.MessageHandler))
	r.HandleFunc("/api/messages/", c.MessageAPI) // Pour récupérer les messages d'un thread
	
	// API pour les réactions
	r.HandleFunc("/api/reactions", middleware.RequireAuth(c.ReactionHandler))
	r.HandleFunc("/api/reactions/", middleware.RequireAuth(c.ReactionAPI))
}

// RegisterPage affiche la page d'inscription
func (c *UserControllers) RegisterPage(w http.ResponseWriter, r *http.Request) {
	// Si déjà connecté, rediriger vers le profil
	if user := middleware.GetUserFromContext(r); user != nil {
		http.Redirect(w, r, "/profile", http.StatusFound)
		return
	}

	http.ServeFile(w, r, "./website/template/register.html")
}

// LoginPage affiche la page de connexion
func (c *UserControllers) LoginPage(w http.ResponseWriter, r *http.Request) {
	// Si déjà connecté, rediriger vers le profil
	if user := middleware.GetUserFromContext(r); user != nil {
		http.Redirect(w, r, "/profile", http.StatusFound)
		return
	}

	http.ServeFile(w, r, "./website/template/login.html")
}

// HomePage affiche la page d'accueil
func (c *UserControllers) HomePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./website/template/home.html")
}

// ThemePage affiche la page de thème
func (c *UserControllers) ThemePage(w http.ResponseWriter, r *http.Request) {
	log.Printf("🎨 ThemePage - Affichage de la page de thème")
	http.ServeFile(w, r, "./website/template/theme.html")
}

// ThreadsListPage affiche la liste de tous les threads
func (c *UserControllers) ThreadsListPage(w http.ResponseWriter, r *http.Request) {
	log.Printf("🧵 ThreadsListPage - Début de la fonction")

	// Récupérer les paramètres de pagination
	page := 1
	limit := 20
	
	if pageParam := r.URL.Query().Get("page"); pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}
	
	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Récupérer les threads avec pagination
	threads, meta, err := c.threadService.GetThreadsWithPagination(page, limit)
	if err != nil {
		log.Printf("❌ Erreur récupération threads: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// Récupérer les catégories
	categories, err := c.threadService.GetCategories()
	if err != nil {
		log.Printf("❌ Erreur récupération catégories: %v", err)
		categories = []models.Category{} // Valeur par défaut
	}

	log.Printf("✅ ThreadsListPage - %d threads trouvés, %d catégories, page %d/%d", 
		len(threads), len(categories), meta.Page, meta.TotalPages)

	// Lire le template
	templatePath := "./website/template/threads_list.html"
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		log.Printf("❌ Erreur lecture template: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// Traiter le template
	htmlContent := string(templateContent)
	processedHTML := processThreadsListTemplateWithPagination(htmlContent, threads, categories, meta)

	// Envoyer la réponse
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(processedHTML))
}

// ThreadsDemoPage affiche la page de démonstration
func (c *UserControllers) ThreadsDemoPage(w http.ResponseWriter, r *http.Request) {
	log.Printf("🎨 ThreadsDemoPage - Affichage de la page de démo")
	http.ServeFile(w, r, "./website/template/threads_demo.html")
}

// ProfilePage affiche la page de profil (nécessite authentification)
func (c *UserControllers) ProfilePage(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔍 ProfilePage - Début de la fonction")
	
	// Récupérer l'utilisateur depuis le contexte
	sessionInfo := middleware.GetUserFromContext(r)
	if sessionInfo == nil {
		log.Printf("❌ ProfilePage - Aucune session trouvée")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	log.Printf("✅ ProfilePage - Session trouvée: UserID=%d, Username=%s", sessionInfo.UserID, sessionInfo.Username)

	// Récupérer les données complètes de l'utilisateur
	user, err := c.authService.GetUserByID(sessionInfo.UserID)
	if err != nil {
		log.Printf("❌ Erreur récupération profil: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ ProfilePage - Données utilisateur récupérées: %s (%s)", user.Username, user.Email)

	// Lire le template HTML
	templatePath := "./website/template/profile.html"
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		log.Printf("❌ Erreur lecture template: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ ProfilePage - Template lu, taille: %d bytes", len(templateContent))

	// Remplacer les placeholders par les vraies données
	htmlContent := string(templateContent)
	
	// Traitement spécifique pour chaque placeholder
	htmlWithUserData := processProfileTemplate(htmlContent, user)

	log.Printf("✅ ProfilePage - Template traité, envoi de la réponse")

	// Envoyer la réponse
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(htmlWithUserData))
}

// RegisterHandler gère l'inscription des utilisateurs
func (c *UserControllers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer les données du formulaire
	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	log.Printf("📝 Tentative d'inscription: %s (%s)", username, email)

	// Validation basique
	if username == "" || email == "" || password == "" {
		log.Printf("❌ Champs manquants pour l'inscription")
		showErrorPage(w, r, "Tous les champs sont requis", "/register")
		return
	}

	// Gérer l'upload de l'image de profil (optionnel)
	var profilePicturePath *string
	file, header, err := r.FormFile("profile_picture")
	if err == nil {
		// Un fichier a été téléchargé
		defer file.Close()
		
		uploadedPath, uploadErr := c.uploadService.UploadProfilePicture(file, header)
		if uploadErr != nil {
			log.Printf("❌ Erreur upload image: %v", uploadErr)
			showErrorPage(w, r, "Erreur lors du téléchargement de l'image: "+uploadErr.Error(), "/register")
			return
		}
		
		profilePicturePath = &uploadedPath
		log.Printf("✅ Image de profil téléchargée: %s", uploadedPath)
	} else {
		// Aucun fichier téléchargé, utiliser l'image par défaut
		defaultPath := c.uploadService.GetDefaultAvatarPath()
		profilePicturePath = &defaultPath
		log.Printf("📷 Utilisation de l'image par défaut: %s", defaultPath)
	}

	// Créer la requête d'inscription
	registerReq := models.RegisterRequest{
		Username:       username,
		Email:          email,
		Password:       password,
		ProfilePicture: profilePicturePath,
	}

	// Appeler le service d'inscription
	user, err := c.authService.Register(registerReq)
	if err != nil {
		log.Printf("❌ Erreur inscription: %v", err)
		
		// Si une image a été téléchargée et que l'inscription échoue, la supprimer
		if profilePicturePath != nil && *profilePicturePath != c.uploadService.GetDefaultAvatarPath() {
			c.uploadService.DeleteProfilePicture(*profilePicturePath)
		}
		
		showErrorPage(w, r, err.Error(), "/register")
		return
	}

	log.Printf("✅ Inscription réussie: %s (ID: %d)", user.Username, user.ID)

	// Redirection avec message de succès
	http.Redirect(w, r, "/login?message=inscription_reussie", http.StatusSeeOther)
}

// LoginHandler gère la connexion des utilisateurs
func (c *UserControllers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer les données du formulaire
	identifier := strings.TrimSpace(r.FormValue("identifiant"))
	password := r.FormValue("password")

	log.Printf("🔑 Tentative de connexion: %s", identifier)

	// Validation basique
	if identifier == "" || password == "" {
		log.Printf("❌ Identifiants manquants")
		showErrorPage(w, r, "Tous les champs sont requis", "/login")
		return
	}

	// Créer la requête de connexion
	loginReq := models.LoginRequest{
		Identifier: identifier,
		Password:   password,
	}

	// Appeler le service de connexion
	user, token, err := c.authService.Login(loginReq)
	if err != nil {
		log.Printf("❌ Erreur connexion: %v", err)
		showErrorPage(w, r, "Identifiants incorrects", "/login")
		return
	}

	// Créer le cookie d'authentification
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true,
		Secure:   false, // Mettre à true en production avec HTTPS
		SameSite: http.SameSiteStrictMode,
		MaxAge:   24 * 60 * 60, // 24 heures
		Path:     "/",
	})

	log.Printf("✅ Connexion réussie: %s (ID: %d)", user.Username, user.ID)

	// Redirection vers le profil
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

// LogoutHandler gère la déconnexion
func (c *UserControllers) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Supprimer le cookie d'authentification
	http.SetCookie(w, &http.Cookie{
		Name:   "auth_token",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	})

	log.Printf("👋 Déconnexion effectuée")

	// Redirection vers la page d'accueil
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

// ProfileAPI retourne les données du profil en JSON
func (c *UserControllers) ProfileAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer l'utilisateur depuis le contexte
	sessionInfo := middleware.GetUserFromContext(r)
	if sessionInfo == nil {
		http.Error(w, "Non authentifié", http.StatusUnauthorized)
		return
	}

	// Récupérer les données complètes de l'utilisateur
	user, err := c.authService.GetUserByID(sessionInfo.UserID)
	if err != nil {
		log.Printf("❌ Erreur récupération profil: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// Retourner en JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.APIResponse{
		Success: true,
		Data:    user,
	})
}

// Helper functions

// showErrorPage affiche une page d'erreur simple ou redirige avec un message
func showErrorPage(w http.ResponseWriter, r *http.Request, message string, redirectTo string) {
	// Pour l'instant, on fait une redirection simple
	// En production, vous pourriez créer une vraie page d'erreur
	http.Redirect(w, r, redirectTo+"?error="+message, http.StatusSeeOther)
}

// GetCurrentUser helper pour récupérer l'utilisateur actuel
func (c *UserControllers) GetCurrentUser(r *http.Request) (*models.User, error) {
	sessionInfo := middleware.GetUserFromContext(r)
	if sessionInfo == nil {
		return nil, nil
	}

	return c.authService.GetUserByID(sessionInfo.UserID)
}

// RequireJSON middleware pour s'assurer que le content-type est JSON
func RequireJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Content-Type doit être application/json", http.StatusBadRequest)
			return
		}
		next(w, r)
	}
}

// ParseJSONBody helper pour parser le body JSON
func ParseJSONBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// WriteJSONResponse helper pour écrire une réponse JSON
func WriteJSONResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// WriteErrorResponse helper pour écrire une réponse d'erreur JSON
func WriteErrorResponse(w http.ResponseWriter, message string, status int) {
	WriteJSONResponse(w, models.APIResponse{
		Success: false,
		Error:   message,
	}, status)
}

// ParseIntParam helper pour parser un paramètre entier depuis l'URL
func ParseIntParam(r *http.Request, param string) (int, error) {
	value := r.URL.Query().Get(param)
	if value == "" {
		return 0, nil
	}
	return strconv.Atoi(value)
}

// processProfileTemplate remplace les placeholders dans le template de profil avec les vraies données
func processProfileTemplate(htmlContent string, user *models.User) string {
	log.Printf("🔄 Traitement template pour utilisateur: %s", user.Username)
	
	// Compter les placeholders avant traitement
	countBefore := strings.Count(htmlContent, "%s") + strings.Count(htmlContent, "%x")
	log.Printf("📊 Placeholders trouvés: %d", countBefore)
	
	// Déterminer l'image de profil à utiliser
	profilePicture := "/img/avatars/default-avatar.png"
	if user.ProfilePicture != nil && *user.ProfilePicture != "" {
		profilePicture = *user.ProfilePicture
		log.Printf("🖼️ Utilisation image personnalisée: %s", profilePicture)
	} else {
		log.Printf("🖼️ Utilisation image par défaut: %s (ProfilePicture=%v)", profilePicture, user.ProfilePicture)
	}
	
	// Remplacer l'image de profil dans la bannière
	htmlContent = strings.Replace(htmlContent, `src="../img/avatar/photo-profil.jpg"`, 
		fmt.Sprintf(`src="%s"`, profilePicture), 1)
	
	// Remplacer également l'avatar dans les posts (même image)
	htmlContent = strings.Replace(htmlContent, `src="../img/avatar/avatar-utilisateur.jpg"`, 
		fmt.Sprintf(`src="%s"`, profilePicture), 1)
	
	// Remplacer les placeholders spécifiques
	htmlContent = strings.Replace(htmlContent, `<h1 class="name">%s</h1>`, 
		fmt.Sprintf(`<h1 class="name">%s</h1>`, user.Username), 1)
	
	htmlContent = strings.Replace(htmlContent, `<span class="handle">@%s</span>`, 
		fmt.Sprintf(`<span class="handle">@%s</span>`, user.Username), 1)
	
	// Date d'inscription
	joinDate := user.CreatedAt.Format("January 2006")
	htmlContent = strings.Replace(htmlContent, `Joined September 2024`, 
		fmt.Sprintf(`Joined %s`, joinDate), 1)
	
	// Stats Following/Followers
	htmlContent = strings.Replace(htmlContent, `<span><strong>%x</strong> Following</span>`, 
		fmt.Sprintf(`<span><strong>%d</strong> Following</span>`, user.FollowingCount), 1)
	
	htmlContent = strings.Replace(htmlContent, `<span><strong>%x</strong> Followers</span>`, 
		fmt.Sprintf(`<span><strong>%d</strong> Followers</span>`, user.FollowerCount), 1)
	
	// Post utilisateur dans le mur
	htmlContent = strings.Replace(htmlContent, `<span class="post-user-name">%s</span>`, 
		fmt.Sprintf(`<span class="post-user-name">%s</span>`, user.Username), 1)
	
	htmlContent = strings.Replace(htmlContent, `<span class="post-user-handle">%s</span>`, 
		fmt.Sprintf(`<span class="post-user-handle">@%s</span>`, user.Username), 1)
	
	htmlContent = strings.Replace(htmlContent, `Félicitations %s pour ta nouvelle page ! 🎉`, 
		fmt.Sprintf(`Félicitations %s pour ta nouvelle page ! 🎉`, user.Username), 1)
	
	// Compter les placeholders après traitement
	countAfter := strings.Count(htmlContent, "%s") + strings.Count(htmlContent, "%x")
	log.Printf("✅ Template traité. Placeholders restants: %d", countAfter)
	
	return htmlContent
}

// processThreadDetailTemplate traite le template de détail d'un thread
func processThreadDetailTemplate(htmlContent string, thread models.Thread) string {
	log.Printf("🔄 Traitement template thread détail - Thread ID=%d", thread.ID)

	// Récupérer le nom de l'auteur et son avatar
	authorName := "Utilisateur inconnu"
	authorUsername := "unknown"
	authorAvatar := "/img/avatars/default-avatar.png"
	if thread.Author != nil {
		authorName = thread.Author.Username
		authorUsername = thread.Author.Username
		if thread.Author.ProfilePicture != nil && *thread.Author.ProfilePicture != "" {
			authorAvatar = *thread.Author.ProfilePicture
		}
	}

	// Récupérer le nom de la catégorie
	categoryName := "Général"
	if thread.Category != nil {
		categoryName = thread.Category.Name
	}

	// Formater la date
	timeAgo := formatTimeAgo(thread.CreatedAt)
	formattedDate := thread.CreatedAt.Format("15:04 · 2 Jan 2006")

	// Remplacer l'avatar de l'auteur dans le template
	htmlContent = strings.Replace(htmlContent, `src="../img/avatar/photo-profil.jpg"`, 
		fmt.Sprintf(`src="%s"`, authorAvatar), 1)

	// Remplacer toutes les informations du thread
	htmlContent = strings.ReplaceAll(htmlContent, "%THREAD_ID%", fmt.Sprintf("%d", thread.ID))
	htmlContent = strings.ReplaceAll(htmlContent, "%THREAD_TITLE%", thread.Title)
	htmlContent = strings.ReplaceAll(htmlContent, "%THREAD_CONTENT%", thread.Content)
	
	// Informations de l'auteur
	htmlContent = strings.ReplaceAll(htmlContent, "%AUTHOR_NAME%", authorName)
	htmlContent = strings.ReplaceAll(htmlContent, "%AUTHOR_USERNAME%", authorUsername)
	htmlContent = strings.ReplaceAll(htmlContent, "%AUTHOR_HANDLE%", "@"+authorUsername)
	
	// Dates et temps
	htmlContent = strings.ReplaceAll(htmlContent, "%CREATED_AT%", formattedDate)
	htmlContent = strings.ReplaceAll(htmlContent, "%THREAD_TIME%", timeAgo)
	htmlContent = strings.ReplaceAll(htmlContent, "%THREAD_DATE%", formattedDate)
	
	// Statistiques
	htmlContent = strings.ReplaceAll(htmlContent, "%VIEW_COUNT%", fmt.Sprintf("%d", thread.ViewCount))
	htmlContent = strings.ReplaceAll(htmlContent, "%LIKE_COUNT%", fmt.Sprintf("%d", thread.LikeCount))
	htmlContent = strings.ReplaceAll(htmlContent, "%DISLIKE_COUNT%", "0") // Pas encore implémenté
	htmlContent = strings.ReplaceAll(htmlContent, "%MESSAGE_COUNT%", fmt.Sprintf("%d", thread.MessageCount))
	
	// Catégorie
	htmlContent = strings.ReplaceAll(htmlContent, "%CATEGORY_NAME%", categoryName)
	
	// Tag de catégorie avec style
	categoryTagHTML := fmt.Sprintf(`<span class="category-tag">📁 %s</span>`, categoryName)
	htmlContent = strings.ReplaceAll(htmlContent, "%CATEGORY_TAG%", categoryTagHTML)
	
	// Hashtags
	hashtagsHTML := ""
	if len(thread.Hashtags) > 0 {
		for _, tag := range thread.Hashtags {
			hashtagsHTML += fmt.Sprintf(`<span class="hashtag">#%s</span>`, tag)
		}
	}
	htmlContent = strings.ReplaceAll(htmlContent, "%HASHTAGS%", hashtagsHTML)
	
	// Messages de réponse - les récupérer et les afficher
	messagesHTML := ""
	if len(thread.Messages) > 0 {
		for _, message := range thread.Messages {
			authorAvatar := "/img/avatars/default-avatar.png"
			if message.Author != nil && message.Author.ProfilePicture != nil && *message.Author.ProfilePicture != "" {
				authorAvatar = *message.Author.ProfilePicture
			}
			
			authorName := "Utilisateur inconnu"
			authorUsername := "unknown"
			if message.Author != nil {
				authorName = message.Author.Username
				authorUsername = message.Author.Username
			}
			
			messageTime := formatTimeAgo(message.CreatedAt)
			
			messagesHTML += fmt.Sprintf(`
			<div class="message">
				<div class="message-author">
					<img src="%s" alt="Avatar" class="message-avatar">
					<div class="message-author-info">
						<span class="message-author-name">%s</span>
						<span class="message-author-handle">@%s</span>
						<span class="message-time">%s</span>
					</div>
				</div>
				<div class="message-content">
					%s
				</div>
				<div class="message-actions">
					<button class="message-like">👍 %d</button>
				</div>
			</div>`,
				authorAvatar,
				authorName,
				authorUsername,
				messageTime,
				message.Content,
				message.LikeCount,
			)
		}
	} else {
		messagesHTML = `
		<div style="text-align: center; padding: 40px; color: var(--second-text-color);">
			<h3>Aucune réponse pour le moment</h3>
			<p>Soyez le premier à répondre à ce thread !</p>
		</div>`
	}
	htmlContent = strings.ReplaceAll(htmlContent, "%MESSAGES_LIST%", messagesHTML)
	
	log.Printf("✅ Template thread détail traité avec succès")
	return htmlContent
}

// =====================================
// CONTRÔLEURS POUR LES THREADS
// =====================================

// CreateThreadPage affiche la page de création de thread
func (c *UserControllers) CreateThreadPage(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔍 CreateThreadPage - Début de la fonction")
	
	// Récupérer les catégories pour le formulaire
	categories, err := c.threadService.GetCategories()
	if err != nil {
		log.Printf("❌ Erreur récupération catégories: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// Lire le template HTML
	templatePath := "./website/template/create_thread.html"
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		log.Printf("❌ Erreur lecture template: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// Remplacer les catégories dans le template
	htmlContent := string(templateContent)
	categoriesHTML := ""
	for _, category := range categories {
		categoriesHTML += fmt.Sprintf(`<option value="%d">%s</option>`, category.ID, category.Name)
	}
	
	htmlContent = strings.Replace(htmlContent, "%CATEGORIES%", categoriesHTML, 1)

	log.Printf("✅ CreateThreadPage - Template préparé avec %d catégories", len(categories))

	// Envoyer la réponse
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(htmlContent))
}

// CreateThreadHandler gère la création d'un nouveau thread
func (c *UserControllers) CreateThreadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	log.Printf("📝 CreateThreadHandler - Tentative de création de thread")

	// Récupérer l'utilisateur depuis le contexte
	sessionInfo := middleware.GetUserFromContext(r)
	if sessionInfo == nil {
		http.Error(w, "Non authentifié", http.StatusUnauthorized)
		return
	}

	// Récupérer les données du formulaire
	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))
	categoryIDStr := r.FormValue("category_id")
	hashtags := strings.TrimSpace(r.FormValue("hashtags"))

	log.Printf("📝 Données reçues - Titre: %s, Contenu: %d chars, Catégorie: %s", 
		title, len(content), categoryIDStr)

	// Traiter la catégorie
	var categoryID *int
	if categoryIDStr != "" {
		if catID, err := strconv.Atoi(categoryIDStr); err == nil && catID > 0 {
			categoryID = &catID
		}
	}

	// Créer la requête
	request := models.ThreadCreateRequest{
		Title:      title,
		Content:    content,
		CategoryID: categoryID,
	}

	// Traiter les hashtags si fournis
	if hashtags != "" {
		request.Hashtags = c.threadService.ProcessHashtagsFromRequest(hashtags)
	}

	// Créer le thread
	thread, err := c.threadService.CreateThread(request, sessionInfo.UserID)
	if err != nil {
		log.Printf("❌ Erreur création thread: %v", err)
		showErrorPage(w, r, err.Error(), "/create-thread")
		return
	}

	log.Printf("✅ Thread créé avec succès: ID=%d, Titre=%s", thread.ID, thread.Title)

	// Redirection vers le thread créé
	http.Redirect(w, r, fmt.Sprintf("/thread/%d", thread.ID), http.StatusSeeOther)
}

// ThreadPage affiche un thread spécifique
func (c *UserControllers) ThreadPage(w http.ResponseWriter, r *http.Request) {
	// Extraire l'ID du thread depuis l'URL
	path := strings.TrimPrefix(r.URL.Path, "/thread/")
	threadID, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "ID de thread invalide", http.StatusBadRequest)
		return
	}

	log.Printf("🔍 ThreadPage - Affichage du thread ID=%d", threadID)

	// Récupérer le thread
	thread, err := c.threadService.GetThread(threadID)
	if err != nil {
		log.Printf("❌ Erreur récupération thread %d: %v", threadID, err)
		http.Error(w, "Thread non trouvé", http.StatusNotFound)
		return
	}

	// Récupérer les messages du thread
	messages, err := c.messageService.GetMessagesByThread(threadID)
	if err != nil {
		log.Printf("⚠️ Erreur récupération messages pour thread %d: %v", threadID, err)
		messages = []models.Message{} // Valeur par défaut si erreur
	}

	// Attacher les messages au thread
	thread.Messages = messages

	// Lire le template HTML
	templatePath := "./website/template/thread_detail.html"
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		log.Printf("❌ Erreur lecture template: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// Traiter le template avec les données du thread
	htmlContent := string(templateContent)
	processedHTML := processThreadDetailTemplate(htmlContent, *thread)

	log.Printf("✅ ThreadPage - Thread %d affiché avec succès", threadID)

	// Envoyer la réponse
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(processedHTML))
}

// ThreadAPI gère les requêtes API pour les threads
func (c *UserControllers) ThreadAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Récupérer un thread ou tous les threads
		path := strings.TrimPrefix(r.URL.Path, "/api/threads/")
		
		if path == "" {
			// Récupérer tous les threads
			page, _ := strconv.Atoi(r.URL.Query().Get("page"))
			limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
			
			threads, err := c.threadService.GetAllThreads(page, limit)
			if err != nil {
				WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
				return
			}
			
			WriteJSONResponse(w, models.APIResponse{
				Success: true,
				Data:    threads,
			}, http.StatusOK)
			return
		}
		
		// Récupérer un thread spécifique
		threadID, err := strconv.Atoi(path)
		if err != nil {
			WriteErrorResponse(w, "ID de thread invalide", http.StatusBadRequest)
			return
		}
		
		thread, err := c.threadService.GetThread(threadID)
		if err != nil {
			WriteErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		}
		
		WriteJSONResponse(w, models.APIResponse{
			Success: true,
			Data:    thread,
		}, http.StatusOK)
		return
	}

	http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
}

// =====================================
// FONCTIONS HELPER POUR LES TEMPLATES
// =====================================

// processThreadsListTemplate traite le template de liste des threads
func processThreadsListTemplateWithPagination(htmlContent string, threads []models.Thread, categories []models.Category, meta *models.Meta) string {
	log.Printf("🔄 Traitement template threads - %d threads, %d catégories, page %d/%d", 
		len(threads), len(categories), meta.Page, meta.TotalPages)

	// Générer la liste des catégories
	categoriesList := ""
	for _, category := range categories {
		categoriesList += fmt.Sprintf(`<button class="category-pill" data-category="%d">%s</button>`, 
			category.ID, category.Name)
	}

	// Générer la liste des threads
	threadsList := ""
	if len(threads) == 0 {
		threadsList = `
		<div style="text-align: center; padding: 40px; color: #666;">
			<h3>Aucun thread trouvé</h3>
			<p>Soyez le premier à créer un thread !</p>
			<a href="/create-thread" style="background-color: #17bf63; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;">Créer un thread</a>
		</div>`
	} else {
		for _, thread := range threads {
			// Calculer le temps relatif
			timeAgo := formatTimeAgo(thread.CreatedAt)
			
			// Preview du contenu (max 150 chars)
			preview := thread.Content
			if len(preview) > 150 {
				preview = preview[:150] + "..."
			}

			// Générer les hashtags
			hashtags := ""
			for _, tag := range thread.Hashtags {
				hashtags += fmt.Sprintf(`<span class="hashtag">#%s</span>`, tag)
			}

			// Récupérer le nom de l'auteur et son avatar
			authorName := "Utilisateur inconnu"
			authorAvatar := "/img/avatars/default-avatar.png"
			if thread.Author != nil {
				authorName = thread.Author.Username
				if thread.Author.ProfilePicture != nil && *thread.Author.ProfilePicture != "" {
					authorAvatar = *thread.Author.ProfilePicture
				}
			}

			// Récupérer le nom de la catégorie
			categoryName := "Général"
			if thread.Category != nil {
				categoryName = thread.Category.Name
			}

			threadsList += fmt.Sprintf(`
			<div class="thread-card" data-thread-id="%d">
				<div class="thread-main">
					<div class="thread-author">
						<img src="%s" alt="Avatar" class="thread-avatar">
						<div class="author-info">
							<span class="author-name">%s</span>
							<span class="author-handle">@%s</span>
							<span class="thread-time">%s</span>
						</div>
					</div>
					
					<div class="thread-content">
						<h3 class="thread-title">
							<a href="/thread/%d">%s</a>
						</h3>
						<p class="thread-preview">%s</p>
						
						<div class="thread-tags">
							<span class="category-tag">%s</span>
							%s
						</div>
					</div>
				</div>
				
				<div class="thread-stats">
					<div class="stat-item">
						<span class="icon">👁️</span>
						<span class="count">%d</span>
					</div>
					<div class="stat-item">
						<span class="icon">💬</span>
						<span class="count">%d</span>
					</div>
					<div class="stat-item likes">
						<span class="icon">👍</span>
						<span class="count">%d</span>
					</div>
					<div class="stat-item dislikes">
						<span class="icon">👎</span>
						<span class="count">%d</span>
					</div>
					<div class="stat-item loves">
						<span class="icon">❤️</span>
						<span class="count">%d</span>
					</div>
				</div>
			</div>`,
				thread.ID,
				authorAvatar,
				authorName,
				authorName,
				timeAgo,
				thread.ID,
				thread.Title,
				preview,
				categoryName,
				hashtags,
				thread.ViewCount,
				thread.MessageCount,
				thread.LikeCount,
				thread.DislikeCount,
				thread.LoveCount,
			)
		}
	}

	// Statistiques avec métadonnées de pagination
	totalThreads := meta.TotalCount
	todayThreads := 0 // TODO: calculer les threads d'aujourd'hui
	weekThreads := 0  // TODO: calculer les threads de la semaine

	// Remplacer les placeholders
	htmlContent = strings.ReplaceAll(htmlContent, "%CATEGORIES_LIST%", categoriesList)
	htmlContent = strings.ReplaceAll(htmlContent, "%THREADS_LIST%", threadsList)
	htmlContent = strings.ReplaceAll(htmlContent, "%TOTAL_THREADS%", fmt.Sprintf("%d", totalThreads))
	htmlContent = strings.ReplaceAll(htmlContent, "%TODAY_THREADS%", fmt.Sprintf("%d", todayThreads))
	htmlContent = strings.ReplaceAll(htmlContent, "%WEEK_THREADS%", fmt.Sprintf("%d", weekThreads))
	
	// Trending threads (simplifié)
	trendingThreads := ""
	if len(threads) > 0 {
		// Prendre les 3 premiers threads comme "trending"
		for i, thread := range threads {
			if i >= 3 { break }
			trendingThreads += fmt.Sprintf(`
			<div class="trending-item">
				<span class="trending-title">%s</span>
				<span class="trending-stats">%d 👍 • %d 💬</span>
			</div>`,
				thread.Title,
				thread.LikeCount,
				thread.MessageCount,
			)
		}
	}
	htmlContent = strings.ReplaceAll(htmlContent, "%TRENDING_THREADS%", trendingThreads)

	// Catégories populaires
	popularCategories := ""
	categoryCount := make(map[string]int)
	for _, thread := range threads {
		categoryName := "Général"
		if thread.Category != nil {
			categoryName = thread.Category.Name
		}
		categoryCount[categoryName]++
	}
	
	for categoryName, count := range categoryCount {
		popularCategories += fmt.Sprintf(`
		<div class="category-item">
			<span class="category-icon">📂</span>
			<span class="category-name">%s</span>
			<span class="category-count">%d</span>
		</div>`,
			categoryName,
			count,
		)
	}
	htmlContent = strings.ReplaceAll(htmlContent, "%POPULAR_CATEGORIES%", popularCategories)

	// Génération de la pagination dynamique
	pagination := generatePagination(meta)
	htmlContent = strings.ReplaceAll(htmlContent, "%PAGINATION%", pagination)

	return htmlContent
}

// generatePagination génère le HTML de pagination
func generatePagination(meta *models.Meta) string {
	if meta.TotalPages <= 1 {
		return `<div class="pagination">
			<button class="page-btn" disabled>← Précédent</button>
			<div class="page-numbers">
				<button class="page-num active" data-page="1">1</button>
			</div>
			<button class="page-btn" disabled>Suivant →</button>
		</div>`
	}

	var pagination strings.Builder
	pagination.WriteString(`<div class="pagination">`)

	// Bouton Précédent
	if meta.Page > 1 {
		pagination.WriteString(fmt.Sprintf(`<button class="page-btn" data-page="%d">← Précédent</button>`, meta.Page-1))
	} else {
		pagination.WriteString(`<button class="page-btn" disabled>← Précédent</button>`)
	}

	pagination.WriteString(`<div class="page-numbers">`)

	// Logique d'affichage des numéros de page
	start := 1
	end := meta.TotalPages

	// Si plus de 7 pages, on affiche intelligemment
	if meta.TotalPages > 7 {
		if meta.Page <= 4 {
			// Début : 1 2 3 4 5 ... dernière
			end = 5
		} else if meta.Page >= meta.TotalPages-3 {
			// Fin : 1 ... (n-4) (n-3) (n-2) (n-1) n
			start = meta.TotalPages - 4
		} else {
			// Milieu : 1 ... (current-1) current (current+1) ... dernière
			start = meta.Page - 1
			end = meta.Page + 1
		}
	}

	// Première page si pas dans la plage
	if start > 1 {
		pagination.WriteString(`<button class="page-num" data-page="1">1</button>`)
		if start > 2 {
			pagination.WriteString(`<span class="page-dots">...</span>`)
		}
	}

	// Pages dans la plage
	for i := start; i <= end; i++ {
		if i == meta.Page {
			pagination.WriteString(fmt.Sprintf(`<button class="page-num active" data-page="%d">%d</button>`, i, i))
		} else {
			pagination.WriteString(fmt.Sprintf(`<button class="page-num" data-page="%d">%d</button>`, i, i))
		}
	}

	// Dernière page si pas dans la plage
	if end < meta.TotalPages {
		if end < meta.TotalPages-1 {
			pagination.WriteString(`<span class="page-dots">...</span>`)
		}
		pagination.WriteString(fmt.Sprintf(`<button class="page-num" data-page="%d">%d</button>`, meta.TotalPages, meta.TotalPages))
	}

	pagination.WriteString(`</div>`)

	// Bouton Suivant
	if meta.Page < meta.TotalPages {
		pagination.WriteString(fmt.Sprintf(`<button class="page-btn" data-page="%d">Suivant →</button>`, meta.Page+1))
	} else {
		pagination.WriteString(`<button class="page-btn" disabled>Suivant →</button>`)
	}

	pagination.WriteString(`</div>`)

	return pagination.String()
}

// formatTimeAgo formate une date en temps relatif
func formatTimeAgo(createdAt time.Time) string {
	now := time.Now()
	diff := now.Sub(createdAt)

	if diff < time.Minute {
		return "à l'instant"
	} else if diff < time.Hour {
		minutes := int(diff.Minutes())
		return fmt.Sprintf("il y a %d min", minutes)
	} else if diff < time.Hour*24 {
		hours := int(diff.Hours())
		return fmt.Sprintf("il y a %d h", hours)
	} else {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "il y a 1 jour"
		} else {
			return fmt.Sprintf("il y a %d jours", days)
		}
	}
}

// =====================================
// HANDLERS POUR LES RÉACTIONS
// =====================================

// ReactionHandler gère les requêtes de création/suppression de réactions
func (c *UserControllers) ReactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteErrorResponse(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	log.Printf("🔄 ReactionHandler - Nouvelle demande de réaction")

	// Récupérer l'utilisateur depuis le contexte
	sessionInfo := middleware.GetUserFromContext(r)
	if sessionInfo == nil {
		WriteErrorResponse(w, "Non authentifié", http.StatusUnauthorized)
		return
	}

	// Parser la requête JSON
	var request models.ReactionRequest
	err := ParseJSONBody(r, &request)
	if err != nil {
		log.Printf("❌ Erreur parsing JSON: %v", err)
		WriteErrorResponse(w, "Données JSON invalides", http.StatusBadRequest)
		return
	}

	log.Printf("📝 Données reçues - TargetType: %s, TargetID: %d, ReactionType: %s", 
		request.TargetType, request.TargetID, request.ReactionType)

	// Traiter la réaction
	reaction, err := c.reactionService.ProcessReaction(sessionInfo.UserID, request)
	if err != nil {
		log.Printf("❌ Erreur traitement réaction: %v", err)
		WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Réponse selon si la réaction a été ajoutée ou supprimée
	if reaction == nil {
		// Réaction supprimée
		WriteJSONResponse(w, models.APIResponse{
			Success: true,
			Message: "Réaction supprimée",
			Data: map[string]interface{}{
				"action": "removed",
				"target_type": request.TargetType,
				"target_id": request.TargetID,
				"reaction_type": request.ReactionType,
			},
		}, http.StatusOK)
	} else {
		// Réaction ajoutée/modifiée
		WriteJSONResponse(w, models.APIResponse{
			Success: true,
			Message: "Réaction ajoutée",
			Data: map[string]interface{}{
				"action": "added",
				"reaction": reaction,
			},
		}, http.StatusCreated)
	}

	log.Printf("✅ ReactionHandler - Réaction traitée avec succès")
}

// ReactionAPI gère les requêtes GET pour récupérer des informations sur les réactions
func (c *UserControllers) ReactionAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteErrorResponse(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer l'utilisateur (optionnel pour cette route)
	sessionInfo := middleware.GetUserFromContext(r)
	userID := 0
	if sessionInfo != nil {
		userID = sessionInfo.UserID
	}

	// Parser les paramètres de l'URL
	targetType := r.URL.Query().Get("target_type")
	targetIDStr := r.URL.Query().Get("target_id")

	if targetType == "" || targetIDStr == "" {
		WriteErrorResponse(w, "Paramètres target_type et target_id requis", http.StatusBadRequest)
		return
	}

	targetID, err := strconv.Atoi(targetIDStr)
	if err != nil {
		WriteErrorResponse(w, "target_id doit être un entier", http.StatusBadRequest)
		return
	}

	log.Printf("🔍 ReactionAPI - TargetType: %s, TargetID: %d, UserID: %d", targetType, targetID, userID)

	// Traiter selon le type de cible
	switch targetType {
	case "thread":
		summary, err := c.reactionService.GetReactionSummary(userID, targetID)
		if err != nil {
			log.Printf("❌ Erreur récupération résumé thread: %v", err)
			WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		WriteJSONResponse(w, models.APIResponse{
			Success: true,
			Data:    summary,
		}, http.StatusOK)

	case "message":
		counts, err := c.reactionService.GetMessageReactionCounts(targetID)
		if err != nil {
			log.Printf("❌ Erreur récupération comptes message: %v", err)
			WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Récupérer la réaction de l'utilisateur si connecté
		var userReaction *string
		if userID > 0 {
			reaction, err := c.reactionService.GetUserMessageReaction(userID, targetID)
			if err != nil {
				log.Printf("⚠️ Erreur récupération réaction utilisateur: %v", err)
			} else if reaction != nil {
				userReaction = &reaction.ReactionType
			}
		}

		response := map[string]interface{}{
			"counts":        counts,
			"user_reaction": userReaction,
		}

		WriteJSONResponse(w, models.APIResponse{
			Success: true,
			Data:    response,
		}, http.StatusOK)

	default:
		WriteErrorResponse(w, "Type de cible non supporté", http.StatusBadRequest)
	}
}

// =====================================
// HANDLERS POUR LES MESSAGES
// =====================================

// MessageHandler gère les requêtes de création de messages
func (c *UserControllers) MessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteErrorResponse(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	log.Printf("💬 MessageHandler - Nouvelle demande de création de message")

	// Récupérer l'utilisateur depuis le contexte
	sessionInfo := middleware.GetUserFromContext(r)
	if sessionInfo == nil {
		WriteErrorResponse(w, "Non authentifié", http.StatusUnauthorized)
		return
	}

	// Parser les données du formulaire
	threadIDStr := r.FormValue("thread_id")
	content := strings.TrimSpace(r.FormValue("content"))

	if threadIDStr == "" || content == "" {
		WriteErrorResponse(w, "thread_id et content sont requis", http.StatusBadRequest)
		return
	}

	threadID, err := strconv.Atoi(threadIDStr)
	if err != nil {
		WriteErrorResponse(w, "thread_id doit être un entier", http.StatusBadRequest)
		return
	}

	// Créer la requête de message
	request := models.MessageCreateRequest{
		ThreadID: threadID,
		Content:  content,
	}

	log.Printf("📝 Création message - ThreadID: %d, AuthorID: %d", threadID, sessionInfo.UserID)

	// Créer le message
	message, err := c.messageService.CreateMessage(request, sessionInfo.UserID)
	if err != nil {
		log.Printf("❌ Erreur création message: %v", err)
		WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("✅ Message créé avec succès - ID: %d", message.ID)

	// Rediriger vers le thread
	redirectURL := fmt.Sprintf("/thread/%d", threadID)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// MessageAPI gère les requêtes GET pour récupérer les messages d'un thread
func (c *UserControllers) MessageAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteErrorResponse(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Extraire l'ID du thread depuis l'URL
	path := strings.TrimPrefix(r.URL.Path, "/api/messages/")
	if path == "" {
		WriteErrorResponse(w, "ID de thread requis", http.StatusBadRequest)
		return
	}

	threadID, err := strconv.Atoi(path)
	if err != nil {
		WriteErrorResponse(w, "ID de thread invalide", http.StatusBadRequest)
		return
	}

	log.Printf("🔍 MessageAPI - Récupération messages pour thread %d", threadID)

	// Récupérer les messages
	messages, err := c.messageService.GetMessagesByThread(threadID)
	if err != nil {
		log.Printf("❌ Erreur récupération messages: %v", err)
		WriteErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("✅ %d messages trouvés pour thread %d", len(messages), threadID)

	WriteJSONResponse(w, models.APIResponse{
		Success: true,
		Data:    messages,
	}, http.StatusOK)
}
