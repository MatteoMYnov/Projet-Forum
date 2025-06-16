(function () {
        const saved = localStorage.getItem('themeSettings');
        if (!saved) return;

        const values = JSON.parse(saved);
        
        // Variables communes
        document.documentElement.style.setProperty('--main-color', values.main);
        document.documentElement.style.setProperty('--main-color-hover', values.hover);
        document.documentElement.style.setProperty('--main-text-color', values.text);
        document.documentElement.style.setProperty('--second-text-color', values.secondary);
        
        // Variables de bordure (compatibilité)
        document.documentElement.style.setProperty('--main-border-color', values.border);
        document.documentElement.style.setProperty('--border-color', values.border);
        
        // Variables spécifiques à thread.css
        document.documentElement.style.setProperty('--primary-bg', values.main);
        document.documentElement.style.setProperty('--secondary-bg', values.hover);
        document.documentElement.style.setProperty('--card-bg', values.hover);
        document.documentElement.style.setProperty('--hover-bg', values.hover);
        document.documentElement.style.setProperty('--main-text', values.text);
        document.documentElement.style.setProperty('--secondary-text', values.secondary);
        document.documentElement.style.setProperty('--muted-text', values.secondary);
        
        // Variables pour les formulaires et éléments interactifs
        document.documentElement.style.setProperty('--input-background', values.hover);
        document.documentElement.style.setProperty('--input-focus-background', values.hover);
        document.documentElement.style.setProperty('--accent-color', values.accent || '#1d9bf0');
        document.documentElement.style.setProperty('--accent-color-alpha', values.accentAlpha || 'rgba(29, 155, 240, 0.2)');
        document.documentElement.style.setProperty('--accent-gradient', values.accentGradient || 'linear-gradient(135deg, #1d9bf0 0%, #7877c6 50%, #ff649e 100%)');
        
        // Variables pour les erreurs
        document.documentElement.style.setProperty('--error-color', '#ff6b6b');
        document.documentElement.style.setProperty('--error-background', 'rgba(244, 33, 46, 0.1)');
        document.documentElement.style.setProperty('--error-border', 'rgba(244, 33, 46, 0.3)');
        
        // Déclencher la mise à jour des icônes après le chargement du thème
        document.addEventListener('DOMContentLoaded', function() {
            // Attendre que logo_switcher.js soit chargé et que updateIcons soit disponible
            setTimeout(function() {
                if (typeof updateIcons === 'function') {
                    updateIcons();
                } else if (typeof window.updateIcons === 'function') {
                    window.updateIcons();
                }
            }, 100);
        });
    })();