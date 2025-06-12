// Fonction pour calculer la luminosité d'une couleur hexadécimale
function getLuminance(hexColor) {
    // Convertir la couleur hex en RGB
    const r = parseInt(hexColor.slice(1, 3), 16);
    const g = parseInt(hexColor.slice(3, 5), 16);
    const b = parseInt(hexColor.slice(5, 7), 16);
    
    // Calculer la luminosité relative selon la formule de luminance
    return (0.299 * r + 0.587 * g + 0.114 * b) / 255;
}

// Stocker les chemins originaux des icônes
const originalIconPaths = new Map();

// Fonction pour mettre à jour les icônes en fonction de la luminosité
function updateIcons() {
    const mainColor = getComputedStyle(document.documentElement).getPropertyValue('--main-color').trim();
    const luminance = getLuminance(mainColor);
    
    // Mettre à jour le logo
    const logoImg = document.querySelector('.icon-logo');
    const favicon = document.getElementById('favicon');
    const logoPath = luminance > 0.5 ? '/img/logo/classic.png' : '/img/logo/inverted.png';
    
    if (logoImg) {
        logoImg.src = logoPath;
    }
    if (favicon) {
        favicon.href = logoPath;
    }
    
    // Mettre à jour les icônes de navigation
    const navIcons = {
        'home-r.png': 'home-c.png',
        'brush-r.png': 'brush-c.png',
        'profile-r.png': 'profile-c.png'
    };
    
    // Sélectionner toutes les icônes de navigation
    const icons = document.querySelectorAll('.icon');
    icons.forEach(icon => {
        // Si c'est la première fois qu'on voit cette icône, stocker son chemin original
        if (!originalIconPaths.has(icon)) {
            originalIconPaths.set(icon, icon.src);
        }
        
        const originalSrc = originalIconPaths.get(icon);
        const fileName = originalSrc.split('/').pop();
        
        // Si c'est une icône que nous voulons changer
        if (navIcons[fileName]) {
            // Construire le nouveau chemin
            const newFileName = luminance > 0.5 ? navIcons[fileName] : fileName;
            const newSrc = originalSrc.replace(fileName, newFileName);
            icon.src = newSrc;
        }
    });
}

// Observer les changements de la variable CSS --main-color
function observeThemeChanges() {
    const observer = new MutationObserver((mutations) => {
        mutations.forEach((mutation) => {
            if (mutation.type === 'attributes' && mutation.attributeName === 'style') {
                updateIcons();
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
        updateIcons();
        observeThemeChanges();
    } else {
        // Attendre que le thème soit chargé
        const checkTheme = setInterval(() => {
            if (document.documentElement.style.getPropertyValue('--main-color')) {
                updateIcons();
                observeThemeChanges();
                clearInterval(checkTheme);
            }
        }, 100);
    }
}

// Initialiser au chargement de la page
document.addEventListener('DOMContentLoaded', initializeLogoSwitcher); 