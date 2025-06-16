// Fonction pour calculer la luminosité d'une couleur hexadécimale
function getLuminance(hexColor) {
    // Enlever le # si présent
    hexColor = hexColor.replace('#', '');
    
    // Si la couleur n'est pas au bon format, retourner 0 (sombre)
    if (hexColor.length !== 6) {
        console.log('Format de couleur incorrect:', hexColor);
        return 0;
    }
    
    // Convertir la couleur hex en RGB
    const r = parseInt(hexColor.slice(0, 2), 16);
    const g = parseInt(hexColor.slice(2, 4), 16);
    const b = parseInt(hexColor.slice(4, 6), 16);
    
    console.log('RGB:', r, g, b);
    
    // Calculer la luminosité relative selon la formule de luminance
    const luminance = (0.299 * r + 0.587 * g + 0.114 * b) / 255;
    console.log('Luminance calculée:', luminance);
    return luminance;
}

// Stocker les chemins originaux des icônes
const originalIconPaths = new Map();

// Fonction pour déterminer si on doit utiliser des icônes sombres ou claires
function shouldUseDarkIcons() {
    const mainColor = getComputedStyle(document.documentElement).getPropertyValue('--main-color').trim();
    console.log('Couleur principale détectée:', mainColor);
    
    // Si pas de couleur ou noir, utiliser icônes claires
    if (!mainColor || mainColor === '' || mainColor === '#000000' || mainColor === 'rgb(0, 0, 0)') {
        console.log('Thème sombre détecté (noir ou vide)');
        return false; // icônes claires sur fond sombre
    }
    
    // Convertir rgb() en hex si nécessaire
    let hexColor = mainColor;
    if (mainColor.startsWith('rgb')) {
        const rgbMatch = mainColor.match(/rgb\((\d+),\s*(\d+),\s*(\d+)\)/);
        if (rgbMatch) {
            const r = parseInt(rgbMatch[1]);
            const g = parseInt(rgbMatch[2]);
            const b = parseInt(rgbMatch[3]);
            hexColor = '#' + [r, g, b].map(x => x.toString(16).padStart(2, '0')).join('');
            console.log('Conversion RGB vers HEX:', mainColor, '->', hexColor);
        }
    }
    
    const luminance = getLuminance(hexColor);
    console.log('Décision icônes:', luminance > 0.3 ? 'sombres (fond clair)' : 'claires (fond sombre)');
    
    // Si luminance > 0.3, c'est un fond clair, donc utiliser icônes sombres
    return luminance > 0.3;
}

// Fonction pour mettre à jour les icônes
function updateIcons() {
    console.log('=== Mise à jour des icônes ===');
    
    const useDarkIcons = shouldUseDarkIcons();
    
    // Mettre à jour le logo
    const logoImg = document.querySelector('.icon-logo');
    const favicon = document.getElementById('favicon');
    
    const logoPath = useDarkIcons ? '/img/logo/classic.png' : '/img/logo/inverted.png';
    console.log('Logo choisi:', logoPath, '(sombre:', useDarkIcons, ')');
    
    if (logoImg) {
        logoImg.src = logoPath;
        console.log('Logo mis à jour');
    }
    if (favicon) {
        favicon.href = logoPath;
        console.log('Favicon mis à jour');
    }
    
    // Mettre à jour les icônes de navigation
    const icons = document.querySelectorAll('.icon');
    console.log('Nombre d\'icônes trouvées:', icons.length);
    
    icons.forEach((icon, index) => {
        // Stocker le chemin original si c'est la première fois
        if (!originalIconPaths.has(icon)) {
            originalIconPaths.set(icon, icon.src);
            console.log('Stockage original icône', index, ':', icon.src);
        }
        
        const originalSrc = originalIconPaths.get(icon);
        const fileName = originalSrc.split('/').pop();
        
        // Déterminer la nouvelle icône
        let newFileName = fileName;
        
        if (useDarkIcons) {
            // Fond clair -> utiliser icônes sombres (-c.png)
            if (fileName.includes('-r.png')) {
                newFileName = fileName.replace('-r.png', '-c.png');
            }
        } else {
            // Fond sombre -> utiliser icônes claires (-r.png)  
            if (fileName.includes('-c.png')) {
                newFileName = fileName.replace('-c.png', '-r.png');
            }
        }
        
        if (newFileName !== fileName) {
            const newSrc = originalSrc.replace(fileName, newFileName);
            console.log('Changement icône', index, ':', fileName, '->', newFileName);
            icon.src = newSrc;
        } else {
            console.log('Icône', index, 'inchangée:', fileName);
        }
    });
    
    console.log('=== Fin mise à jour des icônes ===');
}

// Fonction pour forcer le thème sombre
function updateIconsForDarkTheme() {
    console.log('Application forcée du thème sombre');
    
    const logoImg = document.querySelector('.icon-logo');
    const favicon = document.getElementById('favicon');
    
    if (logoImg) logoImg.src = '/img/logo/inverted.png';
    if (favicon) favicon.href = '/img/logo/inverted.png';
    
    const icons = document.querySelectorAll('.icon');
    icons.forEach(icon => {
        if (!originalIconPaths.has(icon)) {
            originalIconPaths.set(icon, icon.src);
        }
        
        const originalSrc = originalIconPaths.get(icon);
        const fileName = originalSrc.split('/').pop();
        
        // Pour le thème sombre, utiliser les icônes claires (-r.png)
        if (fileName.includes('-c.png')) {
            const newSrc = originalSrc.replace('-c.png', '-r.png');
            icon.src = newSrc;
        }
    });
}

// Observer les changements de la variable CSS --main-color
function observeThemeChanges() {
    const observer = new MutationObserver((mutations) => {
        mutations.forEach((mutation) => {
            if (mutation.type === 'attributes' && mutation.attributeName === 'style') {
                console.log('Changement de style détecté');
                setTimeout(updateIcons, 50); // Petit délai pour laisser le temps au style de s'appliquer
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
    console.log('=== Initialisation du changeur de logo ===');
    
    // Forcer une première mise à jour après un délai
    setTimeout(() => {
        console.log('Mise à jour forcée après délai');
        updateIcons();
    }, 200);
    
    // Observer les changements
    observeThemeChanges();
    
    // Vérification périodique pour les changements non détectés
    setInterval(() => {
        updateIcons();
    }, 2000);
}

// Exposer updateIcons globalement
window.updateIcons = updateIcons;

// Initialiser au chargement de la page
document.addEventListener('DOMContentLoaded', initializeLogoSwitcher);

// Initialiser immédiatement si le DOM est déjà chargé
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initializeLogoSwitcher);
} else {
    initializeLogoSwitcher();
} 