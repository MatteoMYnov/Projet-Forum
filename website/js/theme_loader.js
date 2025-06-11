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
    })();