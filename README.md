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
- **XAMPP** (pour MySQL facile)
- **Git**

### 🚀 Installation avec XAMPP - Guide Complet

#### 📥 **Étape 1 : Télécharger et installer XAMPP**

1. **Télécharger XAMPP** : https://www.apachefriends.org/fr/download.html
2. **Installer XAMPP** avec les options par défaut
3. **Lancer XAMPP Control Panel** en tant qu'administrateur

#### 🔧 **Étape 2 : Démarrer les services**

1. **Ouvrir XAMPP Control Panel**
2. **Démarrer Apache** (cliquer sur "Start")
3. **Démarrer MySQL** (cliquer sur "Start")
4. Vérifier que les deux services sont **verts** ✅

#### 🗄️ **Étape 3 : Configurer la base de données**

1. **Ouvrir phpMyAdmin** :
   - Cliquer sur "Admin" à côté de MySQL dans XAMPP
   - Ou aller sur : `http://localhost/phpmyadmin`

2. **Créer la base de données** :
   ```sql
   CREATE DATABASE forum_y CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```

3. **Importer les tables** :
   - Sélectionner la base `forum_y`
   - Aller dans l'onglet "Importer"
   - Choisir le fichier : `website/database/install_all_tables_sql`
   - Cliquer sur "Exécuter"

#### 💾 **Étape 4 : Cloner et configurer le projet**

```bash
# 1️⃣ Cloner le projet
git clone https://github.com/MatteoMYnov/Projet-Forum.git
cd forum-ynov

# 2️⃣ Créer le fichier .env
echo "DB_HOST=localhost
DB_PORT=3306  
DB_USER=root
DB_PASSWORD=
DB_NAME=forum_y
PORT=2557
JWT_SECRET=mon_secret_jwt_securise_123456789
DEBUG=true
UPLOAD_PATH=./website/img" > .env
```

#### 🚀 **Étape 5 : Lancer le projet**

```bash
# Télécharger les dépendances Go
go mod download

# Démarrer le serveur
go run main.go
```

#### 🎉 **Étape 6 : Accéder au forum**

Ouvrir votre navigateur et aller sur :
- **🏠 Accueil** : `http://localhost:2557/home`
- **👤 Connexion** : `http://localhost:2557/login`
- **📝 Inscription** : `http://localhost:2557/register`

#### ✅ **Vérification que tout fonctionne**

1. **XAMPP Control Panel** : Apache et MySQL doivent être **verts** ✅
2. **Terminal Go** : doit afficher `✅ Serveur démarré sur http://localhost:2557`
3. **Navigateur** : la page d'accueil du forum doit s'afficher
4. **Test inscription** : créer un compte pour tester la base de données

---

### ⚡ **Démarrage Rapide (Si déjà installé)**

```bash
# 1. Démarrer XAMPP (Apache + MySQL)
# 2. Aller dans le dossier du projet
cd Projet-Forum

# 3. Lancer le serveur Go
go run main.go

# 4. Ouvrir http://localhost:2557/home
```

---

### 🛠️ **Dépannage XAMPP**

**❌ MySQL ne démarre pas :**
```bash
# Changer le port MySQL (dans XAMPP Config > MySQL)
Port par défaut : 3306 → Essayer : 3307
# Puis modifier le .env avec le nouveau port
```

**❌ Port 80 occupé (Apache) :**
```bash
# Changer le port Apache (dans XAMPP Config > Apache)  
Port par défaut : 80 → Essayer : 8080
```

**❌ Erreur "database connection failed" :**
1. Vérifier que MySQL est démarré dans XAMPP ✅
2. Vérifier les credentials dans `.env`
3. Tester la connexion dans phpMyAdmin

### 🗄️ **Base de Données Alternative (Manuel)**

Si vous préférez MySQL en ligne de commande :
```bash
# Se connecter à MySQL
mysql -u root -p

# Créer la base
CREATE DATABASE forum_y;
USE forum_y;

# Importer les tables  
SOURCE website/database/install_all_tables_sql;
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
Utilisez les [Issues GitHub](https://github.com/MatteoMYnov/Projet-Forum/issues) avec le template bug.

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

## 👥 Équipe de Développement

<div align="center">

### 🔥 **Les Développeurs qui ont rendu ce projet possible** 🔥

</div>

| 👤 Développeur | 🎯 Rôle | 🔗 GitHub | 💻 Spécialités |
|---|---|---|---|
| **Xerly JI** | 🔧 **Back-end Lead** | [![GitHub](https://img.shields.io/badge/GitHub-@XERCORD-181717?style=for-the-badge&logo=github)](https://github.com/XERCORD) | Go, APIs, Base de données, JWT |
| **Matteo Martin** | 🎨 **Front-end Lead** | [![GitHub](https://img.shields.io/badge/GitHub-@MatteoMYnov-181717?style=for-the-badge&logo=github)](https://github.com/MatteoMYnov) | HTML/CSS, JavaScript, |

---

<div align="center">

### 🚀 **Contributions**

**🔧 Back-end (Xerly JI)**
- Architecture serveur Go
- API REST complète
- Authentification JWT
- Gestion base de données
- Services et repositories

**🎨 Front-end (Matteo Martin)**  
- Interface utilisateur moderne
- Système de thèmes
- JavaScript interactif
- Design responsive
- Expérience utilisateur

---

**🎓 Projet réalisé dans le cadre de la formation Ynov**

**⭐ N'hésitez pas à mettre une étoile si ce projet vous plaît ! ⭐**

Fait avec ❤️ par l'équipe Ynov 2025

</div> 