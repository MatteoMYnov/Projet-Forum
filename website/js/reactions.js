/**
 * Module de gestion des réactions (likes, dislikes, etc.)
 * Utilisé sur les pages de threads et de détails de threads
 */

class ReactionManager {
    constructor() {
        this.init();
    }

    init() {
        console.log('🔄 Initialisation du ReactionManager');
        this.attachEventListeners();
        this.loadReactionStates();
    }

    /**
     * Attache les event listeners aux boutons de réaction
     */
    attachEventListeners() {
        // Gérer les clics sur les boutons de réaction
        document.addEventListener('click', (e) => {
            // Réactions sur les threads
            if (e.target.closest('.reaction-btn')) {
                e.preventDefault();
                this.handleThreadReaction(e.target.closest('.reaction-btn'));
            }
            
            // Réactions sur les messages
            if (e.target.closest('.message-like')) {
                e.preventDefault();
                this.handleMessageReaction(e.target.closest('.message-like'));
            }
        });
    }

    /**
     * Charge les états actuels des réactions pour l'utilisateur
     */
    async loadReactionStates() {
        // Charger les réactions pour le thread actuel si on est sur une page de détail
        const threadId = this.getThreadIdFromPage();
        if (threadId) {
            await this.loadThreadReactions(threadId);
        }

        // Charger les réactions pour tous les messages visibles
        const messageElements = document.querySelectorAll('[data-message-id]');
        for (const element of messageElements) {
            const messageId = element.dataset.messageId;
            if (messageId) {
                await this.loadMessageReactions(messageId);
            }
        }
    }

    /**
     * Gère une réaction sur un thread
     */
    async handleThreadReaction(button) {
        const threadId = this.getThreadIdFromPage();
        const reactionType = button.dataset.action;

        if (!threadId || !reactionType) {
            console.error('❌ Thread ID ou type de réaction manquant');
            return;
        }

        console.log(`🔄 Réaction thread - ID: ${threadId}, Type: ${reactionType}`);

        // Désactiver le bouton temporairement
        const originalDisabled = button.disabled;
        button.disabled = true;

        try {
            const response = await this.sendReaction('thread', threadId, reactionType);
            
            if (response.success) {
                // Mettre à jour l'interface
                await this.updateThreadReactionUI(threadId, reactionType, response.data);
                
                // Animation de feedback
                this.animateButton(button, response.data.action === 'added');
                
                console.log(`✅ Réaction ${response.data.action} avec succès`);
            } else {
                console.error('❌ Erreur réaction:', response.error);
                this.showErrorMessage('Erreur lors de la réaction');
            }
        } catch (error) {
            console.error('❌ Erreur réaction:', error);
            this.showErrorMessage('Erreur de connexion');
        } finally {
            button.disabled = originalDisabled;
        }
    }

    /**
     * Gère une réaction sur un message
     */
    async handleMessageReaction(button) {
        const messageElement = button.closest('[data-message-id]');
        const messageId = messageElement?.dataset.messageId;
        const reactionType = 'like'; // Pour l'instant, seulement les likes sur les messages

        if (!messageId) {
            console.error('❌ Message ID manquant');
            return;
        }

        console.log(`🔄 Réaction message - ID: ${messageId}, Type: ${reactionType}`);

        // Désactiver le bouton temporairement
        const originalDisabled = button.disabled;
        button.disabled = true;

        try {
            const response = await this.sendReaction('message', messageId, reactionType);
            
            if (response.success) {
                // Mettre à jour l'interface
                await this.updateMessageReactionUI(messageId, reactionType, response.data);
                
                // Animation de feedback
                this.animateButton(button, response.data.action === 'added');
                
                console.log(`✅ Réaction message ${response.data.action} avec succès`);
            } else {
                console.error('❌ Erreur réaction message:', response.error);
                this.showErrorMessage('Erreur lors de la réaction');
            }
        } catch (error) {
            console.error('❌ Erreur réaction message:', error);
            this.showErrorMessage('Erreur de connexion');
        } finally {
            button.disabled = originalDisabled;
        }
    }

    /**
     * Envoie une réaction au serveur
     */
    async sendReaction(targetType, targetId, reactionType) {
        const response = await fetch('/api/reactions', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                target_type: targetType,
                target_id: parseInt(targetId),
                reaction_type: reactionType
            })
        });

        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }

        return await response.json();
    }

    /**
     * Charge les réactions pour un thread
     */
    async loadThreadReactions(threadId) {
        try {
            const response = await fetch(`/api/reactions/?target_type=thread&target_id=${threadId}`);
            if (response.ok) {
                const data = await response.json();
                if (data.success) {
                    this.updateThreadReactionCounts(threadId, data.data);
                }
            }
        } catch (error) {
            console.log('⚠️ Erreur chargement réactions thread:', error);
        }
    }

    /**
     * Charge les réactions pour un message
     */
    async loadMessageReactions(messageId) {
        try {
            const response = await fetch(`/api/reactions/?target_type=message&target_id=${messageId}`);
            if (response.ok) {
                const data = await response.json();
                if (data.success) {
                    this.updateMessageReactionCounts(messageId, data.data);
                }
            }
        } catch (error) {
            console.log('⚠️ Erreur chargement réactions message:', error);
        }
    }

    /**
     * Met à jour l'interface après une réaction sur un thread
     */
    async updateThreadReactionUI(threadId, reactionType, responseData) {
        // Recharger les comptes à jour
        await this.loadThreadReactions(threadId);
        
        // Mettre à jour l'état visuel des boutons
        this.updateButtonStates(reactionType, responseData.action === 'added');
    }

    /**
     * Met à jour l'interface après une réaction sur un message
     */
    async updateMessageReactionUI(messageId, reactionType, responseData) {
        // Recharger les comptes à jour
        await this.loadMessageReactions(messageId);
        
        // Trouver le bouton de like du message
        const messageElement = document.querySelector(`[data-message-id="${messageId}"]`);
        const likeButton = messageElement?.querySelector('.message-like');
        
        if (likeButton) {
            this.updateButtonStates(reactionType, responseData.action === 'added', likeButton.parentElement);
        }
    }

    /**
     * Met à jour les comptes de réactions pour un thread
     */
    updateThreadReactionCounts(threadId, reactionData) {
        const counts = reactionData.counts || {};
        const userReaction = reactionData.user_reaction;

        // Mettre à jour les comptes pour tous les types de réactions
        const reactionTypes = ['like', 'dislike', 'love', 'laugh', 'wow', 'sad', 'angry'];
        
        reactionTypes.forEach(type => {
            const countElement = document.querySelector(`.${type}-btn .count`);
            if (countElement) {
                countElement.textContent = counts[type] || 0;
            }
        });

        // Mettre à jour l'état des boutons selon la réaction de l'utilisateur
        document.querySelectorAll('.reaction-btn').forEach(btn => {
            btn.classList.remove('active', 'user-reacted');
        });

        if (userReaction) {
            const activeBtn = document.querySelector(`[data-action="${userReaction}"]`);
            if (activeBtn) {
                activeBtn.classList.add('active', 'user-reacted');
            }
        }
    }

    /**
     * Met à jour les comptes de réactions pour un message
     */
    updateMessageReactionCounts(messageId, reactionData) {
        const counts = reactionData.counts || {};
        const userReaction = reactionData.user_reaction;

        // Trouver l'élément du message
        const messageElement = document.querySelector(`[data-message-id="${messageId}"]`);
        if (!messageElement) return;

        // Mettre à jour le bouton de like
        const likeButton = messageElement.querySelector('.message-like');
        if (likeButton) {
            const likeCount = counts.like || 0;
            likeButton.textContent = `👍 ${likeCount}`;
            
            // Indiquer si l'utilisateur a liké
            if (userReaction === 'like') {
                likeButton.classList.add('active', 'user-reacted');
            } else {
                likeButton.classList.remove('active', 'user-reacted');
            }
        }
    }

    /**
     * Met à jour l'état visuel des boutons de réaction
     */
    updateButtonStates(reactionType, isActive, container = document) {
        const buttons = container.querySelectorAll('.reaction-btn, .message-like');
        
        buttons.forEach(btn => {
            const btnType = btn.dataset.action || 'like';
            
            if (btnType === reactionType) {
                if (isActive) {
                    btn.classList.add('active', 'user-reacted');
                } else {
                    btn.classList.remove('active', 'user-reacted');
                }
            } else {
                // Désactiver les autres types de réactions (on ne peut avoir qu'une seule réaction)
                btn.classList.remove('active', 'user-reacted');
            }
        });
    }

    /**
     * Animation de feedback pour les boutons
     */
    animateButton(button, isAdded) {
        // Animation de "pulse" pour le feedback
        button.style.transform = 'scale(1.1)';
        button.style.transition = 'transform 0.15s ease';
        
        setTimeout(() => {
            button.style.transform = 'scale(1)';
        }, 150);

        // Changer temporairement la couleur
        const originalColor = button.style.color;
        if (isAdded) {
            button.style.color = '#17bf63'; // Vert pour ajouté
        } else {
            button.style.color = '#666'; // Gris pour supprimé
        }

        setTimeout(() => {
            button.style.color = originalColor;
        }, 300);
    }

    /**
     * Affiche un message d'erreur temporaire
     */
    showErrorMessage(message) {
        // Créer un toast d'erreur temporaire
        const toast = document.createElement('div');
        toast.className = 'error-toast';
        toast.textContent = message;
        toast.style.cssText = `
            position: fixed;
            top: 20px;
            right: 20px;
            background: #f4212e;
            color: white;
            padding: 12px 20px;
            border-radius: 5px;
            z-index: 10000;
            font-size: 14px;
        `;

        document.body.appendChild(toast);

        // Supprimer après 3 secondes
        setTimeout(() => {
            if (toast.parentNode) {
                toast.parentNode.removeChild(toast);
            }
        }, 3000);
    }

    /**
     * Récupère l'ID du thread depuis la page
     */
    getThreadIdFromPage() {
        // Essayer de récupérer depuis l'URL
        const urlMatch = window.location.pathname.match(/\/thread\/(\d+)/);
        if (urlMatch) {
            return parseInt(urlMatch[1]);
        }

        // Essayer de récupérer depuis un élément data
        const threadElement = document.querySelector('[data-thread-id]');
        if (threadElement) {
            return parseInt(threadElement.dataset.threadId);
        }

        // Essayer de récupérer depuis les placeholders du template
        const threadIdElements = document.querySelectorAll('script, [data-thread-id]');
        for (const element of threadIdElements) {
            const content = element.textContent || element.innerHTML;
            const match = content.match(/threadId['":\s]*(\d+)/);
            if (match) {
                return parseInt(match[1]);
            }
        }

        return null;
    }
}

// Initialiser quand le DOM est prêt
document.addEventListener('DOMContentLoaded', () => {
    console.log('🚀 Initialisation des réactions');
    window.reactionManager = new ReactionManager();
});

// Styles CSS pour les réactions actives
const styles = `
    .reaction-btn.active,
    .message-like.active {
        background-color: var(--accent-green, #17bf63) !important;
        color: white !important;
        border-color: var(--accent-green, #17bf63) !important;
    }

    .reaction-btn.user-reacted,
    .message-like.user-reacted {
        font-weight: bold;
    }

    .dislike-btn.active {
        background-color: var(--accent-red, #f4212e) !important;
        border-color: var(--accent-red, #f4212e) !important;
    }

    .love-btn.active {
        background-color: #e91e63 !important;
        border-color: #e91e63 !important;
    }

    .laugh-btn.active {
        background-color: #ff9800 !important;
        border-color: #ff9800 !important;
    }

    .wow-btn.active {
        background-color: #2196f3 !important;
        border-color: #2196f3 !important;
    }

    .sad-btn.active {
        background-color: #607d8b !important;
        border-color: #607d8b !important;
    }

    .angry-btn.active {
        background-color: #f44336 !important;
        border-color: #f44336 !important;
    }

    .reaction-btn:disabled {
        opacity: 0.6;
        cursor: not-allowed;
    }

    .error-toast {
        animation: slideInRight 0.3s ease;
    }

    @keyframes slideInRight {
        from {
            transform: translateX(100%);
            opacity: 0;
        }
        to {
            transform: translateX(0);
            opacity: 1;
        }
    }
`;

// Injecter les styles
const styleSheet = document.createElement('style');
styleSheet.textContent = styles;
document.head.appendChild(styleSheet); 