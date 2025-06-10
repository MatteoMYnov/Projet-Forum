(function () {
        const saved = localStorage.getItem('themeSettings');
        if (!saved) return;

        const values = JSON.parse(saved);
        document.documentElement.style.setProperty('--main-color', values.main);
        document.documentElement.style.setProperty('--main-color-hover', values.hover);
        document.documentElement.style.setProperty('--main-border-color', values.border);
        document.documentElement.style.setProperty('--main-text-color', values.text);
        document.documentElement.style.setProperty('--second-text-color', values.secondary);
    })();