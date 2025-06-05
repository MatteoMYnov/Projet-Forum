Projet : Site web dynamique avec navigation (Go + HTML/CSS/JS)

Structure du projet :
.
├── config/
├── website/
│   ├── database/
│   ├── img/
│   │   ├── avatars/
│   │   ├── banners/
│   │   ├── icon/
│   │   └── logo/
│   ├── js/
│   │   └── home.js
│   ├── styles/
│   │   ├── home.css
│   │   ├── login.css
│   │   ├── profile.css
│   │   └── register.css
│   └── template/
│       ├── home.html
│       ├── login.html
│       ├── profile.html
│       └── register.html
├── .env
├── .gitignore
├── go.mod
├── go.sum
└── main.go

Technos :
- Backend : Go (serveur HTTP)
- Frontend : HTML5, CSS3, JavaScript
- Routing Go vers les templates HTML
- Fichiers CSS séparés pour chaque page
- Utilisation de templates dans `/website/template/`
- Les ressources statiques sont dans `/website/img/`, `/website/js/`, `/website/styles/`

Fonctionnalités en place :
- Pages HTML statiques avec navigation latérale (barre contenant Home, Theme, Profile)
- Navigation entre `home.html` et `profile.html` grâce à des boutons cliquables (`<a href="">`)
- `main.go` probablement utilisé pour servir les fichiers statiques et templates HTML

En cours / À venir :
- Ajout d’interactions JS (ex: changement de thème, AJAX)
- Ajout d’une logique serveur pour gestion utilisateur (login, register)
- Intégration de la base de données (prévue dans `/database/`)
