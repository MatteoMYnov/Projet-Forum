// Script de test pour vérifier le changement d'icônes
console.log('=== SCRIPT DE TEST ICONS CHARGÉ ===');

// Test thème blanc
window.testWhiteTheme = function() {
    console.log('🌞 TEST THEME BLANC');
    const whiteTheme = {
        main: '#ffffff',
        hover: '#f0f0f0', 
        border: '#cccccc',
        text: '#1c1c1c',
        secondary: '#555555'
    };
    
    // Appliquer le thème blanc
    document.documentElement.style.setProperty('--main-color', whiteTheme.main);
    document.documentElement.style.setProperty('--main-color-hover', whiteTheme.hover);
    document.documentElement.style.setProperty('--main-text-color', whiteTheme.text);
    document.documentElement.style.setProperty('--second-text-color', whiteTheme.secondary);
    document.documentElement.style.setProperty('--main-border-color', whiteTheme.border);
    document.documentElement.style.setProperty('--border-color', whiteTheme.border);
    
    // Forcer la mise à jour des icônes
    setTimeout(() => {
        if (typeof window.updateIcons === 'function') {
            console.log('Forçage updateIcons pour thème blanc');
            window.updateIcons();
        }
    }, 100);
};

// Test thème noir
window.testBlackTheme = function() {
    console.log('🌙 TEST THEME NOIR');
    const blackTheme = {
        main: '#000000',
        hover: '#181818',
        border: '#2F3336',
        text: '#DBDDDE',
        secondary: '#71767B'
    };
    
    // Appliquer le thème noir
    document.documentElement.style.setProperty('--main-color', blackTheme.main);
    document.documentElement.style.setProperty('--main-color-hover', blackTheme.hover);
    document.documentElement.style.setProperty('--main-text-color', blackTheme.text);
    document.documentElement.style.setProperty('--second-text-color', blackTheme.secondary);
    document.documentElement.style.setProperty('--main-border-color', blackTheme.border);
    document.documentElement.style.setProperty('--border-color', blackTheme.border);
    
    setTimeout(() => {
        if (typeof window.updateIcons === 'function') {
            console.log('Forçage updateIcons pour thème noir');
            window.updateIcons();
        }
    }, 100);
};

// Auto-test au chargement
setTimeout(() => {
    console.log('🧪 AUTO-TEST: détection de la couleur actuelle');
    const currentColor = getComputedStyle(document.documentElement).getPropertyValue('--main-color').trim();
    console.log('Couleur actuelle:', currentColor);
    
    if (typeof window.updateIcons === 'function') {
        console.log('Lancement updateIcons');
        window.updateIcons();
    } else {
        console.log('❌ updateIcons non disponible');
    }
}, 500);

console.log('💡 Utilisez testWhiteTheme() ou testBlackTheme() dans la console pour tester'); 