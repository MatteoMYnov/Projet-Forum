/**
 * Module de gestion de l'authentification pour tous les templates
 * Gère l'affichage des boutons de connexion/déconnexion dans la sidebar
 */

class AuthHandler {
    constructor() {
        this.init();
    }

    init() {
        console.log('🔧 Initialisation AuthHandler');
        this.setupAuthUI();
    }

    async setupAuthUI() {
        const authContainer = document.querySelector('.account');
        if (!authContainer) {
            console.log('⚠️ Container .account non trouvé');
            return;
        }

        try {
            const response = await fetch('/api/profile');
            
            if (response.ok) {
                const data = await response.json();
                const user = data.data;
                this.renderLoggedInUser(authContainer, user);
            } else {
                this.renderGuestUser(authContainer);
            }
        } catch (error) {
            console.log('ℹ️ Utilisateur non connecté');
            this.renderGuestUser(authContainer);
        }
    }

    renderLoggedInUser(container, user) {
        container.innerHTML = `
            <div class="auth-buttons">
                <div class="user-info">
                    <span class="username">@${user.username}</span>
                    <a href="/api/logout" class="logout-btn">Déconnexion</a>
                </div>
            </div>
        `;
        
        // Ajouter les styles inline si pas encore définis
        this.addAuthStyles();
    }

    renderGuestUser(container) {
        container.innerHTML = `
            <div class="auth-buttons">
                <div class="guest-links">
                    <a href="/login" class="auth-link login-link">Connexion</a>
                    <a href="/register" class="auth-link register-link">Inscription</a>
                </div>
            </div>
        `;
        
        // Ajouter les styles inline si pas encore définis
        this.addAuthStyles();
    }

    addAuthStyles() {
        // Vérifier si les styles sont déjà ajoutés
        if (document.getElementById('auth-handler-styles')) {
            return;
        }

        const style = document.createElement('style');
        style.id = 'auth-handler-styles';
        style.textContent = `
            .auth-buttons {
                padding: 10px;
                text-align: center;
            }
            
            .user-info {
                display: flex;
                flex-direction: column;
                gap: 8px;
                align-items: center;
            }
            
            .username {
                color: var(--main-text-color, #fff);
                font-size: 0.9rem;
                font-weight: 500;
            }
            
            .logout-btn {
                color: var(--second-text-color, #999);
                font-size: 0.8rem;
                text-decoration: none;
                transition: color 0.2s;
            }
            
            .logout-btn:hover {
                color: var(--main-text-color, #fff);
            }
            
            .guest-links {
                display: flex;
                flex-direction: column;
                gap: 8px;
                align-items: center;
            }
            
            .auth-link {
                text-decoration: none;
                font-size: 0.9rem;
                transition: color 0.2s;
            }
            
            .login-link {
                color: var(--main-text-color, #fff);
            }
            
            .register-link {
                color: var(--second-text-color, #999);
                font-size: 0.8rem;
            }
            
            .auth-link:hover {
                color: var(--accent-blue, #1d9bf0);
            }
        `;
        
        document.head.appendChild(style);
    }

    // Méthode pour rafraîchir l'état d'authentification
    async refresh() {
        await this.setupAuthUI();
    }
}

// Initialiser l'AuthHandler quand le DOM est chargé
document.addEventListener('DOMContentLoaded', () => {
    console.log('🚀 Chargement AuthHandler');
    window.authHandler = new AuthHandler();
});

// Exposer globalement pour utilisation dans d'autres scripts
window.AuthHandler = AuthHandler; 