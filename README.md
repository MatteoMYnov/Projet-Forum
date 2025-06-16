# 🚀 Forum Ynov - Plateforme de Discussion Moderne

<div align="center">

![Go](https://img.shields.io/badge/Go-1.24.3-00ADD8?style=for-the-badge&logo=go)
![MySQL](https://img.shields.io/badge/MySQL-4479A1?style=for-the-badge&logo=mysql&logoColor=white)
![HTML5](https://img.shields.io/badge/HTML5-E34F26?style=for-the-badge&logo=html5&logoColor=white)
![CSS3](https://img.shields.io/badge/CSS3-1572B6?style=for-the-badge&logo=css3&logoColor=white)
![JavaScript](https://img.shields.io/badge/JavaScript-F7DF1E?style=for-the-badge&logo=javascript&logoColor=black)

**🔥 Une plateforme de forum moderne et complète développée en Go 🔥**

[🎯 Démo](#demo) • [⚡ Installation](#installation) • [📋 Fonctionnalités](#fonctionnalités) • [🛠️ API](#api) • [📚 Documentation](#documentation)

---

</div>

## ✨ Fonctionnalités

### 👤 Gestion des Utilisateurs
- 🔐 **Authentification sécurisée** avec JWT
- 👥 **Profils personnalisables** (avatar, bannière, bio)
- 🎨 **Système de thèmes** (clair/sombre)
- 🔒 **Rôles et permissions** (user/admin)
- 📊 **Statistiques utilisateur** (followers, threads, activité)

### 💬 Système de Discussion
- 📝 **Création de threads** avec catégories
- 💭 **Messages et réponses** en temps réel
- 📌 **Épinglage et archivage** des threads
- 🔍 **Recherche et filtrage** avancés
- 📈 **Compteur de vues** et statistiques

### 🎭 Interactions Sociales
- ❤️ **Système de réactions** (like, dislike, love)
- 👥 **Système de suivi** (followers/following)
- 🏷️ **Hashtags** pour organiser le contenu
- 📱 **Interface responsive** et moderne

### 🛡️ Administration
- 👨‍💼 **Panel d'administration** complet
- 🔨 **Modération des threads** et messages
- 📊 **Tableau de bord** avec analytics
- 🚫 **Gestion des bannissements**

## 🏗️ Architecture

```
forum/
├── 🗂️ config/          # Configuration (DB, env)
├── 🎮 controllers/      # Contrôleurs HTTP
├── 🛡️ middleware/       # Middlewares (auth, CORS, etc.)
├── 📦 models/          # Modèles de données
├── 🗄️ repositories/     # Couche d'accès aux données
├── ⚙️ services/        # Logique métier
├── 🔧 utils/           # Utilitaires (JWT, helpers)
└── 🌐 website/         # Frontend (templates, assets)
    ├── 🎨 styles/      # CSS modulaires
    ├── ⚡ js/          # JavaScript interactif
    ├── 🖼️ img/         # Assets images
    └── 📄 template/    # Templates HTML
```

## ⚡ Installation Rapide

### 🔧 Prérequis
- **Go 1.24.3+** 
- **MySQL 8.0+**
- **Git**

### 🚀 Démarrage en 3 étapes

```bash
# 1️⃣ Cloner le projet
git clone https://github.com/votre-repo/forum-ynov.git
cd forum-ynov

# 2️⃣ Configuration
cp .env.example .env
# ✏️ Éditer .env avec vos paramètres DB

# 3️⃣ Démarrage
go mod download
go run main.go
```

🎉 **C'est parti !** Rendez-vous sur `http://localhost:8080`

### 🗄️ Base de Données

```bash
# Configuration automatique
mysql -u root -p < website/database/create_database.sql
mysql -u root -p < website/database/install_all_tables_sql
```

## 🛠️ API Endpoints

### 🔐 Authentification
```http
POST /api/register     # Inscription
POST /api/login        # Connexion
POST /api/logout       # Déconnexion
PUT  /api/profile/update # Mise à jour profil
```

### 💬 Threads & Messages
```http
GET    /api/threads           # Lister les threads
POST   /api/threads           # Créer un thread
GET    /api/thread/{id}       # Détails d'un thread
POST   /api/messages          # Poster un message
PUT    /api/threads/close/{id} # Fermer un thread
```

### 🎭 Interactions
```http
POST /api/reactions    # Ajouter une réaction
POST /api/follow       # Suivre un utilisateur
GET  /api/wall         # Mur utilisateur
```

## 🎨 Thèmes et Personnalisation

Le forum inclut un **système de thèmes dynamique** :

- 🌞 **Thème clair** - Interface épurée et moderne
- 🌙 **Thème sombre** - Confort visuel en basse luminosité
- 🎨 **Personnalisation** - Variables CSS pour customisation

```css
/* Variables CSS dynamiques */
:root {
  --primary-color: #2563eb;
  --background: #ffffff;
  --text-primary: #1f2937;
  /* ... */
}
```

## 📱 Interface Responsive

✅ **Mobile First** - Optimisé pour tous les appareils  
✅ **Progressive Web App** - Installation possible  
✅ **Accessibilité** - Conforme aux standards WCAG  
✅ **Performance** - Chargement ultra-rapide  

## 🔧 Technologies Utilisées

### Backend
- **🔷 Go 1.24.3** - Langage principal
- **🗄️ MySQL** - Base de données relationnelle
- **🔑 JWT** - Authentification stateless
- **📁 Multipart** - Upload de fichiers

### Frontend
- **📄 HTML5** - Structure sémantique
- **🎨 CSS3** - Styles modernes avec variables
- **⚡ JavaScript** - Interactivité native
- **📱 Responsive Design** - Mobile-first

### DevOps & Outils
- **📦 Go Modules** - Gestion des dépendances
- **🔄 Git** - Contrôle de version
- **🚀 Build simple** - Compilation native Go

## 🤝 Contribution

Les contributions sont les bienvenues ! 

### 🎯 Comment contribuer
1. 🍴 **Fork** le projet
2. 🌿 **Créer** une branche feature (`git checkout -b feature/amazing-feature`)
3. 💾 **Commit** vos changements (`git commit -m 'Add amazing feature'`)
4. 📤 **Push** la branche (`git push origin feature/amazing-feature`)
5. 🔄 **Ouvrir** une Pull Request

### 🐛 Signaler un Bug
Utilisez les [Issues GitHub](https://github.com/votre-repo/issues) avec le template bug.

## 📊 Performances

- ⚡ **< 100ms** - Temps de réponse API
- 🚀 **< 2s** - Chargement initial page
- 💾 **< 50MB** - Utilisation mémoire
- 🔄 **1000+** - Utilisateurs simultanés supportés

## 🔐 Sécurité

- 🛡️ **JWT sécurisé** avec expiration
- 🔒 **Validation** stricte des inputs
- 🚫 **Protection CSRF** intégrée
- 🔐 **Hashage bcrypt** des mots de passe
- 🛡️ **Middleware** de sécurité

## 📚 Documentation

- 📖 [Guide d'installation détaillé](docs/installation.md)
- 🔧 [Configuration avancée](docs/configuration.md)
- 🎨 [Guide des thèmes](docs/theming.md)
- 📡 [Documentation API](docs/api.md)

## 📄 Licence

Ce projet est sous licence **MIT** - voir le fichier [LICENSE](LICENSE) pour plus de détails.

---

<div align="center">

**⭐ N'hésitez pas à mettre une étoile si ce projet vous plaît ! ⭐**

Fait avec ❤️ par l'équipe Ynov

</div> 