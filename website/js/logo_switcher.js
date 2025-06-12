// Fonction pour calculer la luminosité d'une couleur hexadécimale
function getLuminance(hexColor) {
    // Convertir la couleur hex en RGB
    const r = parseInt(hexColor.slice(1, 3), 16);
    const g = parseInt(hexColor.slice(3, 5), 16);
    const b = parseInt(hexColor.slice(5, 7), 16);
    
    // Calculer la luminosité relative selon la formule de luminance
    return (0.299 * r + 0.587 * g + 0.114 * b) / 255;
}

// Fonction pour mettre à jour le logo en fonction de la luminosité
function updateLogo() {
    const mainColor = getComputedStyle(document.documentElement).getPropertyValue('--main-color').trim();
    const luminance = getLuminance(mainColor);
    const logoImg = document.querySelector('.icon-logo');
    const favicon = document.getElementById('favicon');
    
    // Si la luminosité est supérieure à 0.5, utiliser le logo classique
    // Sinon, utiliser le logo inversé
    const logoPath = luminance > 0.5 ? '/img/logo/classic.png' : '/img/logo/inverted.png';
    
    if (logoImg) {
        logoImg.src = logoPath;
    }
    if (favicon) {
        favicon.href = logoPath;
    }
}

// Observer les changements de la variable CSS --main-color
function observeThemeChanges() {
    const observer = new MutationObserver((mutations) => {
        mutations.forEach((mutation) => {
            if (mutation.type === 'attributes' && mutation.attributeName === 'style') {
                updateLogo();
            }
        });
    });

    observer.observe(document.documentElement, {
        attributes: true,
        attributeFilter: ['style']
    });
}

// Attendre que le thème soit chargé avant d'initialiser
function initializeLogoSwitcher() {
    // Vérifier si le thème est déjà chargé
    if (document.documentElement.style.getPropertyValue('--main-color')) {
        updateLogo();
        observeThemeChanges();
    } else {
        // Attendre que le thème soit chargé
        const checkTheme = setInterval(() => {
            if (document.documentElement.style.getPropertyValue('--main-color')) {
                updateLogo();
                observeThemeChanges();
                clearInterval(checkTheme);
            }
        }, 100);
    }
}

// Initialiser au chargement de la page
document.addEventListener('DOMContentLoaded', initializeLogoSwitcher); 