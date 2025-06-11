/**
 * Module de gestion de la pagination
 * Gère les clics sur les boutons de pagination et la navigation
 */

class PaginationManager {
    constructor() {
        this.init();
    }

    init() {
        console.log('🔧 Initialisation PaginationManager');
        this.attachEventListeners();
    }

    attachEventListeners() {
        // Écouter les clics sur tous les boutons de pagination
        document.addEventListener('click', (e) => {
            // Boutons de numéro de page
            if (e.target.classList.contains('page-num') && !e.target.disabled) {
                e.preventDefault();
                const page = parseInt(e.target.dataset.page);
                if (page && page > 0) {
                    this.goToPage(page);
                }
            }
            
            // Boutons Précédent/Suivant
            if (e.target.classList.contains('page-btn') && !e.target.disabled) {
                e.preventDefault();
                const page = parseInt(e.target.dataset.page);
                if (page && page > 0) {
                    this.goToPage(page);
                }
            }
        });

        // Gérer les touches de clavier pour navigation rapide
        document.addEventListener('keydown', (e) => {
            // Seulement si on n'est pas dans un champ de saisie
            if (e.target.tagName !== 'INPUT' && e.target.tagName !== 'TEXTAREA') {
                if (e.key === 'ArrowLeft') {
                    // Flèche gauche = page précédente
                    const prevBtn = document.querySelector('.page-btn[data-page]:not([disabled])');
                    if (prevBtn && prevBtn.textContent.includes('Précédent')) {
                        e.preventDefault();
                        this.goToPage(parseInt(prevBtn.dataset.page));
                    }
                } else if (e.key === 'ArrowRight') {
                    // Flèche droite = page suivante
                    const nextBtn = document.querySelector('.page-btn[data-page]:not([disabled])');
                    if (nextBtn && nextBtn.textContent.includes('Suivant')) {
                        e.preventDefault();
                        this.goToPage(parseInt(nextBtn.dataset.page));
                    }
                }
            }
        });
    }

    goToPage(page) {
        console.log(`📄 Navigation vers la page ${page}`);
        
        // Construire l'URL avec le paramètre de page
        const url = new URL(window.location.href);
        url.searchParams.set('page', page);
        
        // Afficher un indicateur de chargement
        this.showLoadingState();
        
        // Naviguer vers la nouvelle page
        window.location.href = url.toString();
    }

    showLoadingState() {
        // Désactiver tous les boutons de pagination pendant le chargement
        const buttons = document.querySelectorAll('.page-btn, .page-num');
        buttons.forEach(btn => {
            btn.disabled = true;
            btn.style.opacity = '0.5';
        });

        // Afficher un indicateur de chargement
        const pagination = document.querySelector('.pagination');
        if (pagination) {
            pagination.style.position = 'relative';
            
            const loader = document.createElement('div');
            loader.className = 'pagination-loader';
            loader.innerHTML = '🔄 Chargement...';
            loader.style.cssText = `
                position: absolute;
                top: 50%;
                left: 50%;
                transform: translate(-50%, -50%);
                background: rgba(0, 0, 0, 0.8);
                color: white;
                padding: 10px 20px;
                border-radius: 5px;
                font-size: 14px;
                z-index: 1000;
            `;
            
            pagination.appendChild(loader);
        }
    }

    // Méthode pour mettre à jour l'état visuel après navigation
    updateVisualState(currentPage) {
        // Retirer la classe active de tous les boutons
        document.querySelectorAll('.page-num').forEach(btn => {
            btn.classList.remove('active');
        });

        // Ajouter la classe active au bouton de la page courante
        const currentBtn = document.querySelector(`.page-num[data-page="${currentPage}"]`);
        if (currentBtn) {
            currentBtn.classList.add('active');
        }
    }

    // Méthode pour obtenir la page actuelle depuis l'URL
    getCurrentPage() {
        const urlParams = new URLSearchParams(window.location.search);
        return parseInt(urlParams.get('page')) || 1;
    }

    // Méthode pour créer des liens de pagination pour SEO
    generateSEOLinks() {
        const pagination = document.querySelector('.pagination');
        if (!pagination) return;

        // Ajouter des liens cachés pour les moteurs de recherche
        const seoLinks = document.createElement('div');
        seoLinks.style.display = 'none';
        seoLinks.innerHTML = `
            <a href="?page=1" rel="first">Première page</a>
            <a href="?page=${this.getCurrentPage() - 1}" rel="prev">Page précédente</a>
            <a href="?page=${this.getCurrentPage() + 1}" rel="next">Page suivante</a>
        `;
        
        pagination.appendChild(seoLinks);
    }
}

// Initialiser la pagination quand le DOM est chargé
document.addEventListener('DOMContentLoaded', () => {
    console.log('🚀 Chargement PaginationManager');
    window.paginationManager = new PaginationManager();
    
    // Générer les liens SEO
    window.paginationManager.generateSEOLinks();
    
    // Mettre à jour l'état visuel
    const currentPage = window.paginationManager.getCurrentPage();
    window.paginationManager.updateVisualState(currentPage);
}); 