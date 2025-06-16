document.addEventListener('DOMContentLoaded', function() {
    const searchForm = document.getElementById('search-form');
    const searchInput = document.getElementById('search-threads');
    const threadsContainer = document.querySelector('.threads-container');
    const threadCards = document.querySelectorAll('.thread-card');

    // Récupérer le terme de recherche depuis l'URL
    const urlParams = new URLSearchParams(window.location.search);
    const searchTerm = urlParams.get('q');

    // Si un terme de recherche est présent dans l'URL, l'afficher dans l'input
    if (searchTerm) {
        searchInput.value = searchTerm;
        searchThreads(searchTerm);
    }

    // Gérer la soumission du formulaire
    if (searchForm) {
        searchForm.addEventListener('submit', function(e) {
            e.preventDefault();
            const searchValue = searchInput.value.trim();
            
            // Mettre à jour l'URL avec le terme de recherche
            const newUrl = searchValue 
                ? `${window.location.pathname}?q=${encodeURIComponent(searchValue)}`
                : window.location.pathname;
            
            // Mettre à jour l'URL sans recharger la page
            window.history.pushState({}, '', newUrl);
            
            // Effectuer la recherche
            searchThreads(searchValue);
        });
    }

    function searchThreads(searchTerm) {
        if (!searchTerm) {
            // Si la recherche est vide, afficher tous les threads
            threadCards.forEach(card => {
                card.style.display = 'block';
            });
            return;
        }

        threadCards.forEach(card => {
            const title = card.querySelector('.thread-title').textContent.toLowerCase();
            const content = card.querySelector('.thread-preview').textContent.toLowerCase();
            const author = card.querySelector('.author-name').textContent.toLowerCase();
            const hashtags = Array.from(card.querySelectorAll('.hashtag'))
                .map(tag => tag.textContent.toLowerCase());

            // Rechercher dans le titre, le contenu, l'auteur et les hashtags
            const isMatch = title.includes(searchTerm.toLowerCase()) ||
                          content.includes(searchTerm.toLowerCase()) ||
                          author.includes(searchTerm.toLowerCase()) ||
                          hashtags.some(tag => tag.includes(searchTerm.toLowerCase()));

            card.style.display = isMatch ? 'block' : 'none';
        });

        // Afficher un message si aucun résultat n'est trouvé
        const visibleThreads = document.querySelectorAll('.thread-card[style="display: block"]');
        if (visibleThreads.length === 0) {
            const noResults = document.createElement('div');
            noResults.className = 'no-results';
            noResults.innerHTML = `
                <p>Aucun thread ne correspond à votre recherche "${searchTerm}"</p>
                <p>Essayez d'autres mots-clés ou consultez tous les threads</p>
            `;
            threadsContainer.appendChild(noResults);
        } else {
            const noResults = document.querySelector('.no-results');
            if (noResults) {
                noResults.remove();
            }
        }
    }
}); 