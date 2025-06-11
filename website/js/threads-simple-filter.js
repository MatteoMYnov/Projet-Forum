// Système de filtrage simple pour les threads
class SimpleThreadsFilter {
    constructor() {
        this.allThreads = [];
        this.currentFilter = 'all';
        this.init();
    }

    init() {
        console.log('🔧 Initialisation SimpleThreadsFilter');
        this.cacheThreads();
        this.attachEventListeners();
        this.applyFilter(this.currentFilter);
    }

    cacheThreads() {
        const threadCards = document.querySelectorAll('.thread-card');
        this.allThreads = Array.from(threadCards).map(card => {
            return {
                element: card,
                id: card.dataset.threadId,
                views: this.extractNumber(card, '.stat-item:nth-child(1) .count'),
                comments: this.extractNumber(card, '.stat-item:nth-child(2) .count'),
                likes: this.extractNumber(card, '.stat-item:nth-child(3) .count'),
                title: card.querySelector('.thread-title a')?.textContent || '',
                time: card.querySelector('.thread-time')?.textContent || '',
                timeValue: this.parseTimeToMinutes(card.querySelector('.thread-time')?.textContent || '')
            };
        });
        console.log(`📊 ${this.allThreads.length} threads en cache`);
    }

    extractNumber(card, selector) {
        const element = card.querySelector(selector);
        if (!element) return 0;
        
        const text = element.textContent || '0';
        // Gérer les formats comme "1.2k", "500", etc.
        if (text.includes('k')) {
            return parseFloat(text.replace('k', '')) * 1000;
        }
        return parseInt(text) || 0;
    }

    parseTimeToMinutes(timeText) {
        // Convertir "il y a 2h", "il y a 3j", etc. en minutes pour le tri
        if (!timeText) return 0;
        
        if (timeText.includes('min') || timeText.includes('m')) {
            const match = timeText.match(/(\d+)/);
            return match ? parseInt(match[1]) : 0;
        }
        
        if (timeText.includes('h')) {
            const match = timeText.match(/(\d+)/);
            return match ? parseInt(match[1]) * 60 : 0;
        }
        
        if (timeText.includes('j')) {
            const match = timeText.match(/(\d+)/);
            return match ? parseInt(match[1]) * 24 * 60 : 0;
        }
        
        if (timeText.includes('semaine')) {
            const match = timeText.match(/(\d+)/);
            return match ? parseInt(match[1]) * 7 * 24 * 60 : 0;
        }
        
        return 0;
    }

    attachEventListeners() {
        document.querySelectorAll('.filter-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                this.handleFilterClick(e.target);
            });
        });

        // Rendre les cartes cliquables
        this.setupThreadCardListeners();
    }

    handleFilterClick(button) {
        // Mettre à jour l'état visuel
        document.querySelectorAll('.filter-btn').forEach(b => b.classList.remove('active'));
        button.classList.add('active');

        // Appliquer le filtre
        const filter = button.dataset.filter;
        this.currentFilter = filter;
        this.applyFilter(filter);
    }

    applyFilter(filter) {
        console.log(`🔄 Application du filtre: ${filter}`);
        
        let sortedThreads = [...this.allThreads];

        // Trier selon le filtre sélectionné
        switch (filter) {
            case 'all':
                // Afficher tous les threads dans l'ordre par défaut (plus récents d'abord)
                sortedThreads.sort((a, b) => a.timeValue - b.timeValue);
                break;
                
            case 'recent':
                // Tri par date (plus récent en premier = moins de minutes)
                sortedThreads.sort((a, b) => a.timeValue - b.timeValue);
                break;
                
            case 'views':
                // Tri par nombre de vues (décroissant)
                sortedThreads.sort((a, b) => b.views - a.views);
                break;
                
            case 'comments':
                // Tri par nombre de commentaires (décroissant)
                sortedThreads.sort((a, b) => b.comments - a.comments);
                break;
                
            default:
                // Garde l'ordre par défaut
                break;
        }

        this.displayThreads(sortedThreads);
    }

    displayThreads(threads) {
        const container = document.querySelector('.threads-container');
        if (!container) return;

        // Réorganiser les éléments dans l'ordre voulu
        threads.forEach((thread, index) => {
            container.appendChild(thread.element);
        });

        console.log(`📊 Threads réorganisés selon: ${this.currentFilter}`);
    }

    setupThreadCardListeners() {
        document.querySelectorAll('.thread-card').forEach(card => {
            card.addEventListener('click', (e) => {
                // Ne pas déclencher si on clique sur un lien
                if (!e.target.closest('a')) {
                    const threadId = card.dataset.threadId;
                    if (threadId) {
                        window.location.href = `/thread/${threadId}`;
                    }
                }
            });

            // Ajouter le style curseur
            card.style.cursor = 'pointer';
        });
    }
}

// Initialiser le filtrage quand la page est chargée
document.addEventListener('DOMContentLoaded', () => {
    console.log('🚀 Chargement SimpleThreadsFilter');
    window.simpleThreadsFilter = new SimpleThreadsFilter();
}); 